package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &indexResource{}
	_ resource.ResourceWithConfigure   = &indexResource{}
	_ resource.ResourceWithImportState = &indexResource{}
)

// indexResource is the resource implementation.
type indexResource struct {
	client *mongo.Client
}

// indexResourceModel maps the resource schema data.
type indexResourceModel struct {
	Database           string            `tfsdk:"database"`
	Collection         string            `tfsdk:"collection"`
	Name               string            `tfsdk:"name"`
	Keys               []indexKey        `tfsdk:"keys"`
	Sparse             *bool             `tfsdk:"sparse"`
	ExpireAfterSeconds *int32            `tfsdk:"expire_after_seconds"`
	Unique             *bool             `tfsdk:"unique"`
	WildcardProjection *map[string]int32 `tfsdk:"wildcard_projection"`
	Collation          *collation        `tfsdk:"collation"`
	background		   *bool             `tfsdk:"background"`

	// see https://developer.hashicorp.com/terraform/plugin/framework/acctests#implement-id-attribute
	Id types.String `tfsdk:"id"`
}

type indexKey struct {
	Field string `tfsdk:"field"`
	Type  string `tfsdk:"type"`
}

type collation struct {
	Locale          string  `tfsdk:"locale"`
	CaseLevel       *bool   `tfsdk:"case_level"`
	CaseFirst       *string `tfsdk:"case_first"`
	Strength        *int    `tfsdk:"strength"`
	NumericOrdering *bool   `tfsdk:"numeric_ordering"`
	Alternate       *string `tfsdk:"alternate"`
	MaxVariable     *string `tfsdk:"max_variable"`
	Normalization   *bool   `tfsdk:"normalization"`
	Backwards       *bool   `tfsdk:"backwards"`
}

// NewIndexResource is a helper function to simplify the provider implementation.
func NewIndexResource() resource.Resource {
	return &indexResource{}
}

// Configure adds the provider configured client to the resource.
func (d *indexResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Info(ctx, "Configuring MongoDB index resource")
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

	d.client = client
	tflog.Info(ctx, "Configured MongoDB index resource")
}

// Metadata returns the resource type name.
func (r *indexResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_index"
}

// Schema defines the schema for the resource.
func (r *indexResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create indexes in MongoDB.",
		Attributes: map[string]schema.Attribute{
			"database": schema.StringAttribute{
				Description: "Name of the database where to create the index.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"collection": schema.StringAttribute{
				Description: "Name of the collection where to create the index.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the index to create.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"keys": schema.ListNestedAttribute{
				Description: "The list of fields composing the index.",
				Required:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"field": schema.StringAttribute{
							Description: "The name of the field to index.",
							Required:    true,
						},
						"type": schema.StringAttribute{
							Description: "The type of index for this field.",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.NoneOf("text"),
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
			"sparse": schema.BoolAttribute{
				Description: "Is it a sparse index.",
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"expire_after_seconds": schema.Int64Attribute{
				Description: "Documents ttl in seconds for ttl indexes.",
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"unique": schema.BoolAttribute{
				Description: "Is it a unique index.",
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"wildcard_projection": schema.MapAttribute{
				Description: "Projection for wirldcard indexes.",
				ElementType: types.Int64Type,
				Optional:    true,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
			"background": schema.BoolAttribute{
				Description: "Create the index in the background.",
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			"collation": schema.SingleNestedAttribute{
				Description: "Index collation.",
				Optional:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: map[string]schema.Attribute{
					"locale": schema.StringAttribute{
						Description: "The locale.",
						Required:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"case_level": schema.BoolAttribute{
						Description: "The case level.",
						Optional:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.RequiresReplace(),
						},
					},
					"case_first": schema.StringAttribute{
						Description: "The case ordering.",
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"strength": schema.Int64Attribute{
						Description: "The number of comparison levels to use.",
						Optional:    true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.RequiresReplace(),
						},
					},
					"numeric_ordering": schema.BoolAttribute{
						Description: "Whether to order numbers based on numerical order and not collation order.",
						Optional:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.RequiresReplace(),
						},
					},
					"alternate": schema.StringAttribute{
						Description: "Whether spaces and punctuation are considered base characters.",
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"max_variable": schema.StringAttribute{
						Description: "Which characters are affected by alternate: 'shifted'.",
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"normalization": schema.BoolAttribute{
						Description: "Causes text to be normalized into Unicode NFD.",
						Optional:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.RequiresReplace(),
						},
					},
					"backwards": schema.BoolAttribute{
						Description: "Causes secondary differences to be considered in reverse order, as it is done in the French language.",
						Optional:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.RequiresReplace(),
						},
					},
				},
			},
			// see https://developer.hashicorp.com/terraform/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed:           true,
				DeprecationMessage: "Just there for compatibility reasons",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *indexResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan indexResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	databaseName := plan.Database
	collectionName := plan.Collection
	indexName := plan.Name

	tflog.Debug(ctx, fmt.Sprintf("Creating index %s.%s.%s", databaseName, collectionName, indexName))

	keys := bson.D{}
	for _, key := range plan.Keys {
		keys = append(keys, bson.E{Key: key.Field, Value: convertToMongoIndexType(key.Type)})
	}

	db := r.client.Database(databaseName)
	collection := db.Collection(collectionName)

	options := &options.IndexOptions{
		Name:               &indexName,
		Sparse:             plan.Sparse,
		ExpireAfterSeconds: plan.ExpireAfterSeconds,
		Unique:             plan.Unique,
		Collation:          plan.Collation.toMongoCollation(),
		Background: 	    plan.background,
	}
	if plan.WildcardProjection != nil {
		options.WildcardProjection = plan.WildcardProjection
	}

	name, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{Keys: keys, Options: options})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create index",
			"An unexpected error occurred when creating index. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Error: "+err.Error(),
		)
		return
	}

	plan.Name = name
	plan.Id = types.StringValue("to_be_ignored")

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Index %s.%s.%s created", databaseName, collectionName, indexName))
}

// Read refreshes the Terraform state with the latest data.
func (r *indexResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state indexResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	databaseName := state.Database
	collectionName := state.Collection
	indexName := state.Name

	tflog.Debug(ctx, fmt.Sprintf("Getting index %s.%s.%s", databaseName, collectionName, indexName))

	db := r.client.Database(databaseName)
	collection := db.Collection(collectionName)
	indexes, err := collection.Indexes().ListSpecifications(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to list indexes",
			"An unexpected error occurred when listing indexes. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Error: "+err.Error(),
		)
		return
	}

	var foundIndex *mongo.IndexSpecification
	for _, index := range indexes {
		if index.Name == indexName {
			foundIndex = index
			break
		}
	}

	if foundIndex == nil {
		resp.Diagnostics.AddError(
			"Unable to find index with name "+indexName,
			"The requested index does not exist. ",
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found index %s.%s.%s", databaseName, collectionName, indexName))

	var foundKeys bson.D
	err = bson.Unmarshal(foundIndex.KeysDocument, &foundKeys)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to parse keys from fetched index",
			"An unexpected error occurred when parsing index keys. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Error: "+err.Error(),
		)
		return
	}

	state.Keys = make([]indexKey, 0)
	for _, v := range foundKeys {
		typ, err := convertToTfIndexType(v.Value)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to convert key type from fetched index",
				"An unexpected error occurred when parsing index keys. "+
					"If the error is not clear, please contact the provider developers.\n\n"+
					"Error: "+err.Error(),
			)
			return
		}
		state.Keys = append(state.Keys, indexKey{Field: v.Key, Type: typ})
	}

	state.Sparse = foundIndex.Sparse
	state.ExpireAfterSeconds = foundIndex.ExpireAfterSeconds
	state.Unique = foundIndex.Unique
	state.Id = types.StringValue("to_be_ignored")

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Refreshed index %s.%s.%s", databaseName, collectionName, indexName))
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *indexResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"An update has been triggered when none should have been.",
		" Changes in index should always result in resource recreation. ",
	)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *indexResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state indexResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete index
	databaseName := state.Database
	collectionName := state.Collection
	indexName := state.Name

	tflog.Debug(ctx, fmt.Sprintf("Dropping index %s.%s.%s", databaseName, collectionName, indexName))

	db := r.client.Database(databaseName)
	collection := db.Collection(collectionName)

	_, err := collection.Indexes().DropOne(ctx, indexName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update (drop) index",
			"An unexpected error occurred when creating index. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Error: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Dropped index %s.%s.%s", databaseName, collectionName, indexName))
}

func (r *indexResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute, parse it and set it has the state of the resource to import
	id, err := parseIndexId(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid id format. Should be <database>.<collection>.<index_name>.",
			"An unexpected error occurred when creating index. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("database"), id.database)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("collection"), id.collection)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), id.indexName)...)
}
