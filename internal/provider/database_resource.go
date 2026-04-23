package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"go.mongodb.org/mongo-driver/mongo"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &databaseResource{}
	_ resource.ResourceWithConfigure   = &databaseResource{}
	_ resource.ResourceWithImportState = &databaseResource{}
)

// databaseResource is the resource implementation.
type databaseResource struct {
	client *mongo.Client
}

// databaseResourceModel maps the resource schema data.
type databaseResourceModel struct {
	Name string       `tfsdk:"name"`
	Id   types.String `tfsdk:"id"`
}

// NewDatabaseResource is a helper function to simplify the provider implementation.
func NewDatabaseResource() resource.Resource {
	return &databaseResource{}
}

// Configure adds the provider configured client to the resource.
func (r *databaseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Info(ctx, "Configuring MongoDB database resource")
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
	tflog.Info(ctx, "Configured MongoDB database resource")
}

// Metadata returns the resource type name.
func (r *databaseResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database"
}

// Schema defines the schema for the resource.
func (r *databaseResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create databases in MongoDB.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name of the database to create.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
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
func (r *databaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan databaseResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	databaseName := plan.Name

	tflog.Debug(ctx, fmt.Sprintf("Creating database %s", databaseName))

	// In MongoDB, databases are created implicitly when you first store data in them.
	// We'll create a dummy collection to ensure the database exists.
	db := r.client.Database(databaseName)
	err := db.CreateCollection(ctx, "_terraform_created")
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create database",
			"An unexpected error occurred when creating database. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Error: "+err.Error(),
		)
		return
	}

	plan.Id = types.StringValue(databaseName)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Database %s created", databaseName))
}

// Read refreshes the Terraform state with the latest data.
func (r *databaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state databaseResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	databaseName := state.Name

	tflog.Debug(ctx, fmt.Sprintf("Reading database %s", databaseName))

	// List all databases to check if our database exists
	databases, err := r.client.ListDatabaseNames(ctx, map[string]interface{}{
		"name": databaseName,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to list databases",
			"An unexpected error occurred when listing databases. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Error: "+err.Error(),
		)
		return
	}

	if len(databases) == 0 {
		resp.Diagnostics.AddError(
			"Database not found",
			fmt.Sprintf("Database %s does not exist", databaseName),
		)
		return
	}

	// Set the state
	state.Id = types.StringValue(databaseName)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Read database %s", databaseName))
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *databaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Updates not supported",
		"Database updates are not supported. Changes to database configuration require recreation.",
	)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *databaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state databaseResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	databaseName := state.Name

	tflog.Debug(ctx, fmt.Sprintf("Dropping database %s", databaseName))

	err := r.client.Database(databaseName).Drop(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to drop database",
			"An unexpected error occurred when dropping database. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Error: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Dropped database %s", databaseName))
}

// ImportState imports an existing resource into Terraform state.
func (r *databaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), req.ID)...)
}
