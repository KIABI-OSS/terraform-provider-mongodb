package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"strconv"

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
	Host               types.String `tfsdk:"host"`
	Port               types.String `tfsdk:"port"`
	Certificate        types.String `tfsdk:"certificate"`
	Username           types.String `tfsdk:"username"`
	Password           types.String `tfsdk:"password"`
	AuthDatabase       types.String `tfsdk:"auth_database"`
	ReplicaSet         types.String `tfsdk:"replica_set"`
	InsecureSkipVerify types.Bool   `tfsdk:"insecure_skip_verify"`
	SSL                types.Bool   `tfsdk:"ssl"`
	Direct             types.Bool   `tfsdk:"direct"`
	RetryWrites        types.Bool   `tfsdk:"retrywrites"`
	Proxy              types.String `tfsdk:"proxy"`
	Url                types.String `tfsdk:"url"`
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
			"host": schema.StringAttribute{
				Optional:    true,
				Description: "The mongodb server address.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("host"),
						path.MatchRoot("url"),
					),
				},
			},
			"port": schema.StringAttribute{
				Optional:    true,
				Description: "The mongodb server port",
			},
			"certificate": schema.StringAttribute{
				Optional:    true,
				Description: "PEM-encoded content of Mongodb host CA certificate",
			},
			"username": schema.StringAttribute{
				Optional:    true,
				Description: "The mongodb user",
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Description: "The mongodb password",
			},
			"auth_database": schema.StringAttribute{
				Optional:    true,
				Description: "The mongodb auth database",
			},
			"replica_set": schema.StringAttribute{
				Optional:    true,
				Description: "The mongodb replica set",
			},
			"insecure_skip_verify": schema.BoolAttribute{
				Optional:    true,
				Description: "ignore hostname verification",
			},
			"ssl": schema.BoolAttribute{
				Optional:    true,
				Description: "ssl activation",
			},
			"direct": schema.BoolAttribute{
				Optional:    true,
				Description: "enforces a direct connection instead of discovery",
			},
			"retrywrites": schema.BoolAttribute{
				Optional:    true,
				Description: "Retryable Writes",
			},
			"proxy": schema.StringAttribute{
				Optional:    true,
				Description: "Proxy through which to connect to MongoDB. Supported protocols are http, https, and socks5. ",
			},
			"url": schema.StringAttribute{
				Optional:    true,
				Description: "The url of the mongodb server.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("host"),
						path.MatchRoot("url"),
					),
				},
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
	if config.Url.ValueString() == "" && config.Host.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing host or url",
			"The provider cannot create the MongoDB client as there is an unknown configuration value for the host. Please specify either host or url.",
		)
		resp.Diagnostics.AddError("fatal", "Missing host or url")
		return
	}

	if config.Url.ValueString() != "" && config.Host.ValueString() != "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Conflicting host and url",
			"The provider cannot create the MongoDB client as there are conflicting configuration values for host and url. Please specify either host or url, but not both.",
		)
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	var uri string
	if config.Url.ValueString() != "" {
		uri = config.Url.ValueString()
	} else {
		var arguments = ""

		arguments = addArgs(arguments, "retrywrites="+strconv.FormatBool(config.RetryWrites.ValueBool()))

		if config.SSL.ValueBool() {
			arguments = addArgs(arguments, "ssl=true")
		}

		if config.ReplicaSet.ValueString() != "" && !config.Direct.ValueBool() {
			arguments = addArgs(arguments, "replicaSet="+config.ReplicaSet.ValueString())
		}

		if config.Direct.ValueBool() {
			arguments = addArgs(arguments, "connect="+"direct")
		}

		uri = "mongodb://" + config.Host.ValueString() + ":" + config.Port.ValueString() + arguments
	}

	// Create a new client using the configuration values
	tflog.Info(ctx, "Creating MongoDB client")

	dialer, dialerErr := proxyDialer(config.Proxy.ValueString())

	if dialerErr != nil {
		resp.Diagnostics.AddError(
			"Unable to create proxy dialer",
			"An unexpected error occurred when creating the proxy dialer. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Error: "+dialerErr.Error(),
		)
		return
	}

	var opts *options.ClientOptions
	var verify = false
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	if config.InsecureSkipVerify.ValueBool() {
		verify = true
	}

	if config.Certificate.ValueString() != "" {
		tlsConfig, err := getTLSConfigWithAllServerCertificates([]byte(config.Certificate.ValueString()), verify)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to read certificate",
				"An unexpected error occurred when reading the certificate. "+
					"If the error is not clear, please contact the provider developers.\n\n"+
					"Error: "+err.Error(),
			)
			return
		}

		opts = options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI).SetAuth(options.Credential{
			AuthSource: config.AuthDatabase.ValueString(), Username: config.Username.ValueString(), Password: config.Password.ValueString(),
		}).SetTLSConfig(tlsConfig).SetDialer(dialer)

	} else {
		opts = options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI).SetAuth(options.Credential{
			AuthSource: config.AuthDatabase.ValueString(), Username: config.Username.ValueString(), Password: config.Password.ValueString(),
		}).SetDialer(dialer)
	}

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
		NewDatabaseResource,
		NewCollectionResource,
	}
}
