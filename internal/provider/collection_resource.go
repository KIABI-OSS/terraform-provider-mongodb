package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &collectionResource{}
	_ resource.ResourceWithConfigure   = &collectionResource{}
	_ resource.ResourceWithImportState = &collectionResource{}
)

// collectionResource is the resource implementation.
type collectionResource struct {
	client *mongo.Client
}

// collectionResourceModel maps the resource schema data.
type collectionResourceModel struct {
	Database   string       `tfsdk:"database"`
	Name       string       `tfsdk:"name"`
	Validation *validation  `tfsdk:"validation"`
	Id         types.String `tfsdk:"id"`
}

type validation struct {
	Validator string `tfsdk:"validator"`
}

// NewCollectionResource is a helper function to simplify the provider implementation.
func NewCollectionResource() resource.Resource {
	return &collectionResource{}
}

// Configure adds the provider configured client to the resource.
func (r *collectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Info(ctx, "Configuring MongoDB collection resource")
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
	tflog.Info(ctx, "Configured MongoDB collection resource")
}

// Metadata returns the resource type name.
func (r *collectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_collection"
}

// Schema defines the schema for the resource.
func (r *collectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create collections in MongoDB.",
		Attributes: map[string]schema.Attribute{
			"database": schema.StringAttribute{
				Description: "Name of the database where to create the collection.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the collection to create.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"validation": schema.SingleNestedAttribute{
				Description: "Collection validation rules.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"validator": schema.StringAttribute{
						Description: "JSON schema validation rules for the collection.",
						Required:    true,
					},
				},
			},
			"id": schema.StringAttribute{
				Computed:           true,
				DeprecationMessage: "Just there for compatibility reasons",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *collectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan collectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	databaseName := plan.Database
	collectionName := plan.Name

	tflog.Debug(ctx, fmt.Sprintf("Creating collection %s.%s", databaseName, collectionName))

	db := r.client.Database(databaseName)

	opts := options.CreateCollection()
	if plan.Validation != nil {
		opts.SetValidator(plan.Validation.Validator)
	}

	err := db.CreateCollection(ctx, collectionName, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create collection",
			"An unexpected error occurred when creating collection. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Error: "+err.Error(),
		)
		return
	}

	plan.Id = types.StringValue(fmt.Sprintf("%s.%s", databaseName, collectionName))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Collection %s.%s created", databaseName, collectionName))
}

// Read refreshes the Terraform state with the latest data.
func (r *collectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state collectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	databaseName := state.Database
	collectionName := state.Name

	tflog.Debug(ctx, fmt.Sprintf("Reading collection %s.%s", databaseName, collectionName))

	db := r.client.Database(databaseName)
	collections, err := db.ListCollectionNames(ctx, map[string]interface{}{
		"name": collectionName,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to list collections",
			"An unexpected error occurred when listing collections. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Error: "+err.Error(),
		)
		return
	}

	if len(collections) == 0 {
		resp.Diagnostics.AddError(
			"Collection not found",
			fmt.Sprintf("Collection %s.%s does not exist", databaseName, collectionName),
		)
		return
	}

	// Set the state
	state.Id = types.StringValue(fmt.Sprintf("%s.%s", databaseName, collectionName))

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Read collection %s.%s", databaseName, collectionName))
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *collectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Updates not supported",
		"Collection updates are not supported. Changes to collection configuration require recreation.",
	)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *collectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state collectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	databaseName := state.Database
	collectionName := state.Name

	tflog.Debug(ctx, fmt.Sprintf("Dropping collection %s.%s", databaseName, collectionName))

	db := r.client.Database(databaseName)
	err := db.Collection(collectionName).Drop(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to drop collection",
			"An unexpected error occurred when dropping collection. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Error: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Dropped collection %s.%s", databaseName, collectionName))
}

// ImportState imports an existing resource into Terraform state.
func (r *collectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := parseCollectionId(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid id format. Should be <database>.<collection>.",
			"An unexpected error occurred when importing collection. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("database"), id.database)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), id.collection)...)
}

type collectionId struct {
	database   string
	collection string
}

func parseCollectionId(id string) (*collectionId, error) {
	parts := strings.Split(id, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format: %s", id)
	}
	return &collectionId{
		database:   parts[0],
		collection: parts[1],
	}, nil
}
