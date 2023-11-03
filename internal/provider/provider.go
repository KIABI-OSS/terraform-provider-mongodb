package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &mongodbProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &mongodbProvider{
			version: version,
		}
	}
}

// mongodbProvider is the provider implementation.
type mongodbProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type mongodbProviderModel struct {
	Url types.String `tfsdk:"url"`
}

// Metadata returns the provider type name.
func (p *mongodbProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "mongodb"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *mongodbProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create resources in MongoDB.",
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Optional:    true,
				Description: "URL of the MongoDB instance to connect to.",
			},
		},
	}
}

func (p *mongodbProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring MongoDB provider")

	// Retrieve provider data from configuration
	var config mongodbProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.
	if config.Url.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("url"),
			"Unknown MongoDB Url",
			"The provider cannot create the MongoDB client as there is an unknown configuration value for the url. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the MONGODB_URL environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	url := os.Getenv("MONGODB_URL")

	if !config.Url.IsNull() {
		url = config.Url.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if url == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("url"),
			"Missing Url",
			"The provider cannot create the MongoDB client as there is a missing or empty value for the url. "+
				"Set the host value in the configuration or use the MONGODB_URL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new client using the configuration values
	tflog.Info(ctx, "Creating MongoDB client")

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(url).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create MongoDB Client",
			"An unexpected error occurred when creating the MongoDB client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Error: "+err.Error(),
		)
		return
	}

	// Make the client available during DataSource and Resource type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured MongoDB provider")
}

// DataSources defines the data sources implemented in the provider.
func (p *mongodbProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// Resources defines the resources implemented in the provider.
func (p *mongodbProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewIndexResource,
	}
}
