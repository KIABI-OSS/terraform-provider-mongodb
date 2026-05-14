package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	_ datasource.DataSource                    = &queryShapeDataSource{}
	_ datasource.DataSourceWithConfigure       = &queryShapeDataSource{}
	_ datasource.DataSourceWithConfigValidators = &queryShapeDataSource{}
)

type queryShapeDataSource struct {
	client *mongo.Client
}

type queryShapeModel struct {
	Database     types.String `tfsdk:"database"`
	Collection   types.String `tfsdk:"collection"`
	Command      types.String `tfsdk:"command"`
	Filter       types.String `tfsdk:"filter"`
	Sort         types.String `tfsdk:"sort"`
	Projection   types.String `tfsdk:"projection"`
	Pipeline     types.String `tfsdk:"pipeline"`
	Hint         types.String `tfsdk:"hint"`
	Collation    types.String `tfsdk:"collation"`
	Skip         types.Int64  `tfsdk:"skip"`
	Limit        types.Int64  `tfsdk:"limit"`
	BatchSize    types.Int64  `tfsdk:"batch_size"`
	AllowDiskUse types.Bool   `tfsdk:"allow_disk_use"`
	Key          types.String `tfsdk:"key"`
	Hash         types.String `tfsdk:"hash"`
}

type queryShapeValidator struct{}

func (v queryShapeValidator) Description(_ context.Context) string {
	return "Validates that command-specific attributes are only used with compatible commands"
}

func (v queryShapeValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v queryShapeValidator) ValidateDataSource(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var config queryShapeModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Command.IsNull() || config.Command.IsUnknown() {
		return
	}

	switch config.Command.ValueString() {
	case "find":
		if !config.Pipeline.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("pipeline"),
				"Invalid Attribute Combination",
				"pipeline cannot be set when command is \"find\"",
			)
		}
		if !config.AllowDiskUse.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("allow_disk_use"),
				"Invalid Attribute Combination",
				"allow_disk_use cannot be set when command is \"find\"",
			)
		}
		if !config.Key.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("key"),
				"Invalid Attribute Combination",
				"key cannot be set when command is \"find\"",
			)
		}
	case "aggregate":
		if !config.Filter.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("filter"),
				"Invalid Attribute Combination",
				"filter cannot be set when command is \"aggregate\"",
			)
		}
		if !config.Sort.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("sort"),
				"Invalid Attribute Combination",
				"sort cannot be set when command is \"aggregate\"",
			)
		}
		if !config.Projection.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("projection"),
				"Invalid Attribute Combination",
				"projection cannot be set when command is \"aggregate\"",
			)
		}
		if !config.Skip.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("skip"),
				"Invalid Attribute Combination",
				"skip cannot be set when command is \"aggregate\"",
			)
		}
		if !config.Limit.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("limit"),
				"Invalid Attribute Combination",
				"limit cannot be set when command is \"aggregate\"",
			)
		}
		if !config.Key.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("key"),
				"Invalid Attribute Combination",
				"key cannot be set when command is \"aggregate\"",
			)
		}
	case "distinct":
		if !config.Sort.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("sort"),
				"Invalid Attribute Combination",
				"sort cannot be set when command is \"distinct\"",
			)
		}
		if !config.Projection.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("projection"),
				"Invalid Attribute Combination",
				"projection cannot be set when command is \"distinct\"",
			)
		}
		if !config.Pipeline.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("pipeline"),
				"Invalid Attribute Combination",
				"pipeline cannot be set when command is \"distinct\"",
			)
		}
		if !config.Skip.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("skip"),
				"Invalid Attribute Combination",
				"skip cannot be set when command is \"distinct\"",
			)
		}
		if !config.Limit.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("limit"),
				"Invalid Attribute Combination",
				"limit cannot be set when command is \"distinct\"",
			)
		}
		if !config.BatchSize.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("batch_size"),
				"Invalid Attribute Combination",
				"batch_size cannot be set when command is \"distinct\"",
			)
		}
		if !config.AllowDiskUse.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("allow_disk_use"),
				"Invalid Attribute Combination",
				"allow_disk_use cannot be set when command is \"distinct\"",
			)
		}
	case "count":
		if !config.Sort.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("sort"),
				"Invalid Attribute Combination",
				"sort cannot be set when command is \"count\"",
			)
		}
		if !config.Projection.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("projection"),
				"Invalid Attribute Combination",
				"projection cannot be set when command is \"count\"",
			)
		}
		if !config.Pipeline.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("pipeline"),
				"Invalid Attribute Combination",
				"pipeline cannot be set when command is \"count\"",
			)
		}
		if !config.Skip.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("skip"),
				"Invalid Attribute Combination",
				"skip cannot be set when command is \"count\"",
			)
		}
		if !config.BatchSize.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("batch_size"),
				"Invalid Attribute Combination",
				"batch_size cannot be set when command is \"count\"",
			)
		}
		if !config.AllowDiskUse.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("allow_disk_use"),
				"Invalid Attribute Combination",
				"allow_disk_use cannot be set when command is \"count\"",
			)
		}
		if !config.Key.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("key"),
				"Invalid Attribute Combination",
				"key cannot be set when command is \"count\"",
			)
		}
	}
}

func NewQueryShapeDataSource() datasource.DataSource {
	return &queryShapeDataSource{}
}

func (d *queryShapeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	tflog.Info(ctx, "Configuring MongoDB query shape datasource")
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*mongo.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected *mongo.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.client = client
	tflog.Info(ctx, "Configured MongoDB query shape datasource")
}

func (d *queryShapeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_query_shape"
}

func (d *queryShapeDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		queryShapeValidator{},
	}
}

func (d *queryShapeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve the query shape hash from MongoDB.",
		Attributes: map[string]schema.Attribute{
			"database": schema.StringAttribute{
				Description: "Name of the database.",
				Required:    true,
			},
			"collection": schema.StringAttribute{
				Description: "Name of the collection.",
				Required:    true,
			},
			"command": schema.StringAttribute{
				Description: "Type of command. One of: find, aggregate, distinct, count.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("find", "aggregate", "distinct", "count"),
				},
			},
			"filter": schema.StringAttribute{
				Description: "JSON string representing the query filter.",
				Optional:    true,
			},
			"sort": schema.StringAttribute{
				Description: "JSON string representing the sort specification.",
				Optional:    true,
			},
			"projection": schema.StringAttribute{
				Description: "JSON string representing the projection specification.",
				Optional:    true,
			},
			"pipeline": schema.StringAttribute{
				Description: "JSON string representing the aggregation pipeline array.",
				Optional:    true,
			},
			"hint": schema.StringAttribute{
				Description: "JSON string representing the index hint.",
				Optional:    true,
			},
			"collation": schema.StringAttribute{
				Description: "JSON string representing the collation specification.",
				Optional:    true,
			},
			"skip": schema.Int64Attribute{
				Description: "Number of documents to skip.",
				Optional:    true,
			},
			"limit": schema.Int64Attribute{
				Description: "Maximum number of documents to return.",
				Optional:    true,
			},
			"batch_size": schema.Int64Attribute{
				Description: "Number of documents to return per batch.",
				Optional:    true,
			},
			"allow_disk_use": schema.BoolAttribute{
				Description: "Allow disk use for aggregation operations.",
				Optional:    true,
			},
			"key": schema.StringAttribute{
				Description: "The field to use for distinct command.",
				Optional:    true,
			},
			"hash": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *queryShapeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan queryShapeModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	db := d.client.Database(plan.Database.ValueString())
	coll := plan.Collection.ValueString()
	cmd := plan.Command.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Getting query shape hash for %s.%s (%s)", plan.Database.ValueString(), coll, cmd))

	var explainCmd bson.D
	switch cmd {
	case "find":
		explainCmd = d.buildFindExplain(plan, coll)
	case "aggregate":
		explainCmd = d.buildAggregateExplain(plan, coll)
	case "distinct":
		explainCmd = d.buildDistinctExplain(plan, coll)
	case "count":
		explainCmd = d.buildCountExplain(plan, coll)
	}

	if explainBytes, err := bson.MarshalExtJSON(explainCmd, false, false); err == nil {
		tflog.Debug(ctx, fmt.Sprintf("Sending explain command: %s", string(explainBytes)))
	}

	var result bson.M
	err := db.RunCommand(ctx, explainCmd).Decode(&result)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to explain query",
			fmt.Sprintf("Failed to run explain for %s: %s", cmd, err.Error()),
		)
		return
	}

	if resultBytes, err := json.Marshal(result); err == nil {
		tflog.Debug(ctx, fmt.Sprintf("Explain response: %s", string(resultBytes)))
	}

	hash := extractQueryShapeHash(result)
	if hash == "" {
		resp.Diagnostics.AddError(
			"queryShapeHash not found",
			"queryShapeHash was not found in the explain output. MongoDB 8.0+ is required.",
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Extracted queryShapeHash: %s", hash))

	plan.Hash = types.StringValue(hash)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (d *queryShapeDataSource) buildFindExplain(plan queryShapeModel, collection string) bson.D {
	args := bson.D{
		{Key: "find", Value: collection},
		{Key: "filter", Value: parseFilter(plan.Filter.ValueString())},
	}

	setIf(&args, "sort", parseFilter(plan.Sort.ValueString()))
	setIf(&args, "projection", parseFilter(plan.Projection.ValueString()))
	if !plan.Skip.IsNull() {
		args = append(args, bson.E{Key: "skip", Value: plan.Skip.ValueInt64()})
	}
	if !plan.Limit.IsNull() {
		args = append(args, bson.E{Key: "limit", Value: plan.Limit.ValueInt64()})
	}
	if !plan.BatchSize.IsNull() {
		args = append(args, bson.E{Key: "batchSize", Value: plan.BatchSize.ValueInt64()})
	}
	setIf(&args, "hint", parseFilter(plan.Hint.ValueString()))
	setCollation(&args, plan.Collation.ValueString())

	return bson.D{{Key: "explain", Value: args}}
}

func (d *queryShapeDataSource) buildAggregateExplain(plan queryShapeModel, collection string) bson.D {
	cursor := bson.M{}
	if plan.BatchSize.ValueInt64() > 0 {
		cursor["batchSize"] = plan.BatchSize.ValueInt64()
	}

	aggArgs := bson.D{
		{Key: "aggregate", Value: collection},
		{Key: "pipeline", Value: parsePipeline(plan.Pipeline.ValueString())},
		{Key: "cursor", Value: cursor},
	}

	if plan.AllowDiskUse.ValueBool() {
		aggArgs = append(aggArgs, bson.E{Key: "allowDiskUse", Value: true})
	}
	setIf(&aggArgs, "hint", parseFilter(plan.Hint.ValueString()))
	setCollation(&aggArgs, plan.Collation.ValueString())

	return bson.D{{Key: "explain", Value: aggArgs}}
}

func (d *queryShapeDataSource) buildDistinctExplain(plan queryShapeModel, collection string) bson.D {
	args := bson.D{
		{Key: "distinct", Value: collection},
		{Key: "key", Value: plan.Key.ValueString()},
	}

	setIf(&args, "query", parseFilter(plan.Filter.ValueString()))
	setCollation(&args, plan.Collation.ValueString())

	return bson.D{{Key: "explain", Value: args}}
}

func (d *queryShapeDataSource) buildCountExplain(plan queryShapeModel, collection string) bson.D {
	args := bson.D{
		{Key: "count", Value: collection},
		{Key: "query", Value: parseFilter(plan.Filter.ValueString())},
	}

	if !plan.Limit.IsNull() {
		args = append(args, bson.E{Key: "limit", Value: plan.Limit.ValueInt64()})
	}
	setIf(&args, "hint", parseFilter(plan.Hint.ValueString()))
	setCollation(&args, plan.Collation.ValueString())

	return bson.D{{Key: "explain", Value: args}}
}

func parseFilter(s string) bson.M {
	if s == "" {
		return bson.M{}
	}
	var v bson.M
	if json.Unmarshal([]byte(s), &v) == nil {
		return v
	}
	return bson.M{}
}

func parsePipeline(s string) bson.A {
	if s == "" {
		return bson.A{}
	}
	var v bson.A
	if json.Unmarshal([]byte(s), &v) == nil {
		return v
	}
	return bson.A{}
}

func setIf(args *bson.D, key string, v bson.M) {
	if v != nil && len(v) > 0 {
		*args = append(*args, bson.E{Key: key, Value: v})
	}
}

func setIfNotZero(args *bson.D, key string, v int64) {
	if v > 0 {
		*args = append(*args, bson.E{Key: key, Value: v})
	}
}

func setCollation(args *bson.D, s string) {
	if s == "" {
		return
	}
	var col bson.M
	if json.Unmarshal([]byte(s), &col) != nil {
		return
	}
	o := options.Collation{Locale: col["locale"].(string)}
	if v, ok := col["caseLevel"].(bool); ok {
		o.CaseLevel = v
	}
	if v, ok := col["caseFirst"].(string); ok {
		o.CaseFirst = v
	}
	if v, ok := col["strength"].(float64); ok {
		o.Strength = int(v)
	}
	if v, ok := col["numericOrdering"].(bool); ok {
		o.NumericOrdering = v
	}
	if v, ok := col["alternate"].(string); ok {
		o.Alternate = v
	}
	if v, ok := col["maxVariable"].(string); ok {
		o.MaxVariable = v
	}
	if v, ok := col["normalization"].(bool); ok {
		o.Normalization = v
	}
	if v, ok := col["backwards"].(bool); ok {
		o.Backwards = v
	}
	*args = append(*args, bson.E{Key: "collation", Value: o})
}

func extractQueryShapeHash(result bson.M) string {
	if h, ok := result["queryShapeHash"].(string); ok && h != "" {
		return h
	}
	if p, ok := result["queryPlanner"].(bson.M); ok {
		if h, ok := p["queryShapeHash"].(string); ok && h != "" {
			return h
		}
	}
	return ""
}