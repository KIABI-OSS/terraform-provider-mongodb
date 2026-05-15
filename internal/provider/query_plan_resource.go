package provider

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	_ resource.Resource                = &queryPlanResource{}
	_ resource.ResourceWithConfigure   = &queryPlanResource{}
	_ resource.ResourceWithImportState = &queryPlanResource{}
)

type queryPlanResource struct {
	client *mongo.Client
}

type queryPlanResourceModel struct {
	QueryHash      types.String `tfsdk:"query_hash"`
	Database       types.String `tfsdk:"database"`
	Collection     types.String `tfsdk:"collection"`
	AllowedIndexes types.List   `tfsdk:"allowed_indexes"`
	Comment        types.String `tfsdk:"comment"`
	Id             types.String `tfsdk:"id"`
}

func NewQueryPlanResource() resource.Resource {
	return &queryPlanResource{}
}

func (r *queryPlanResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Info(ctx, "Configuring MongoDB query plan resource")
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*mongo.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *mongo.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
	tflog.Info(ctx, "Configured MongoDB query plan resource")
}

func (r *queryPlanResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_query_plan"
}

func (r *queryPlanResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage query plan settings in MongoDB. Forces query settings for a given query shape hash.",
		Attributes: map[string]schema.Attribute{
			"query_hash": schema.StringAttribute{
				Description: "The SHA256 query shape hash to apply settings to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		"database": schema.StringAttribute{
			Description: "Name of the database.",
			Required:    true,
		},
		"collection": schema.StringAttribute{
			Description: "Name of the collection.",
			Required:    true,
		},
			"allowed_indexes": schema.ListAttribute{
				Description: "List of index names allowed for this query shape. Set to [\"*\"] to allow all indexes.",
				Required:    true,
				ElementType: types.StringType,
			},
			"comment": schema.StringAttribute{
				Description: "Comment describing the reason for the query settings.",
				Optional:    true,
			},
			"id": schema.StringAttribute{
				Computed:           true,
				DeprecationMessage: "Just there for compatibility reasons",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *queryPlanResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan queryPlanResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Creating query plan for hash %s", plan.QueryHash.ValueString()))

	allowedIndexes := make([]string, 0)
	diags = plan.AllowedIndexes.ElementsAs(ctx, &allowedIndexes, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.setQuerySettings(ctx, plan.QueryHash.ValueString(), plan.Database.ValueString(), plan.Collection.ValueString(), allowedIndexes, plan.Comment.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create query plan",
			"An unexpected error occurred when creating query plan. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Error: "+err.Error(),
		)
		return
	}

	plan.Id = types.StringValue("to_be_ignored")

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Created query plan for hash %s", plan.QueryHash.ValueString()))
}

func (r *queryPlanResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state queryPlanResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	queryHash := state.QueryHash.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Reading query plan for hash %s", queryHash))

	settings, err := r.getQuerySettings(ctx, queryHash)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read query plan",
			"An unexpected error occurred when reading query plan. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Error: "+err.Error(),
		)
		return
	}

	if settings == nil {
		resp.Diagnostics.AddError(
			"Query plan not found",
			fmt.Sprintf("No query settings found for hash %s. The resource may have been deleted outside of Terraform.", queryHash),
		)
		return
	}

	if settings.Comment != "" {
		state.Comment = types.StringValue(settings.Comment)
	}
	state.Id = types.StringValue("to_be_ignored")

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Read query plan for hash %s", queryHash))
}

func (r *queryPlanResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan queryPlanResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Updating query plan for hash %s", plan.QueryHash.ValueString()))

	allowedIndexes := make([]string, 0)
	diags = plan.AllowedIndexes.ElementsAs(ctx, &allowedIndexes, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.setQuerySettings(ctx, plan.QueryHash.ValueString(), plan.Database.ValueString(), plan.Collection.ValueString(), allowedIndexes, plan.Comment.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update query plan",
			"An unexpected error occurred when updating query plan. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Error: "+err.Error(),
		)
		return
	}

	plan.Id = types.StringValue("to_be_ignored")

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Updated query plan for hash %s", plan.QueryHash.ValueString()))
}

func (r *queryPlanResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state queryPlanResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	queryHash := state.QueryHash.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Deleting query plan for hash %s", queryHash))

	adminDb := r.client.Database("admin")
	err := adminDb.RunCommand(ctx, bson.D{
		{Key: "removeQuerySettings", Value: queryHash},
	}).Err()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete query plan",
			"An unexpected error occurred when deleting query plan. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Error: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Deleted query plan for hash %s", queryHash))
}

func (r *queryPlanResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.AddError(
		"Import not supported for mongodb_query_plan",
		"The mongodb_query_plan resource does not support import. "+
			"Query shape hashes cannot be deterministically mapped to a Terraform resource state without additional metadata.",
	)
}

type queryPlanSettings struct {
	Database       string
	Collection     string
	AllowedIndexes []string
	Comment        string
}

func (r *queryPlanResource) setQuerySettings(ctx context.Context, queryHash, database, collection string, allowedIndexes []string, comment string) error {
	adminDb := r.client.Database("admin")

	settings := bson.M{
		"indexHints": bson.M{
			"ns": bson.M{
				"db":   database,
				"coll": collection,
			},
			"allowedIndexes": allowedIndexes,
		},
	}

	if comment != "" {
		settings["comment"] = comment
	}

	cmd := bson.D{
		{Key: "setQuerySettings", Value: queryHash},
		{Key: "settings", Value: settings},
	}

	var result bson.M
	err := adminDb.RunCommand(ctx, cmd).Decode(&result)
	if err != nil {
		return fmt.Errorf("failed to set query settings: %w", err)
	}

	if ok, _ := result["ok"].(float64); ok != 1 {
		errMsg := "unknown"
		if errMsgRaw, exists := result["errmsg"]; exists {
			errMsg = fmt.Sprintf("%v", errMsgRaw)
		}
		return fmt.Errorf("setQuerySettings command failed: %s", errMsg)
	}

	return nil
}

func (r *queryPlanResource) getQuerySettings(ctx context.Context, queryHash string) (*queryPlanSettings, error) {
	adminDb := r.client.Database("admin")

	hashBytes, err := hex.DecodeString(queryHash)
	if err != nil {
		return nil, fmt.Errorf("invalid query hash hex string: %w", err)
	}

	cursor, err := adminDb.Aggregate(ctx, mongo.Pipeline{
		bson.D{{Key: "$querySettings", Value: bson.M{}}},
	}, options.Aggregate().SetAllowDiskUse(true))
	if err != nil {
		return nil, fmt.Errorf("failed to list query settings: %w", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var raw bson.M
		if err := cursor.Decode(&raw); err != nil {
			return nil, fmt.Errorf("failed to decode query settings: %w", err)
		}

		tflog.Debug(ctx, fmt.Sprintf("query settings entry: %v", raw))

		rawHash, ok := raw["queryShapeHash"]
		if !ok {
			tflog.Debug(ctx, "entry missing queryShapeHash field")
			continue
		}

		if hashMatches(rawHash, hashBytes, queryHash) {
			return parseSettingsFromDoc(raw)
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return nil, nil
}

func hashMatches(rawHash interface{}, hashBytes []byte, hashStr string) bool {
	switch v := rawHash.(type) {
	case string:
		// Compare case-insensitively since MongoDB may normalize the hex case
		if len(v) != len(hashStr) {
			return false
		}
		for i := 0; i < len(v); i++ {
			if toUpper(v[i]) != toUpper(hashStr[i]) {
				return false
			}
		}
		return true
	case primitive.Binary:
		return bytes.Equal(v.Data, hashBytes)
	case []byte:
		return bytes.Equal(v, hashBytes)
	}
	return false
}

func toUpper(b byte) byte {
	if b >= 'a' && b <= 'z' {
		return b - 32
	}
	return b
}

func parseSettingsFromDoc(doc bson.M) (*queryPlanSettings, error) {
	settings := &queryPlanSettings{}

	settingsRaw, ok := doc["settings"].(bson.M)
	if !ok {
		return nil, nil
	}

	if comment, ok := settingsRaw["comment"].(string); ok {
		settings.Comment = comment
	}

	indexHints, ok := settingsRaw["indexHints"].(bson.M)
	if !ok {
		return settings, nil
	}

	if ns, ok := indexHints["ns"].(bson.M); ok {
		if db, ok := ns["db"].(string); ok {
			settings.Database = db
		}
		if coll, ok := ns["coll"].(string); ok {
			settings.Collection = coll
		}
	}

	if allowedIndexes, ok := indexHints["allowedIndexes"].(bson.A); ok {
		for _, idx := range allowedIndexes {
			if s, ok := idx.(string); ok {
				settings.AllowedIndexes = append(settings.AllowedIndexes, s)
			}
		}
	}

	return settings, nil
}


