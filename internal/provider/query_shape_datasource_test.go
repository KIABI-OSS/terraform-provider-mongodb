package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQueryShapeDatasourceFind(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "mongodb_query_shape" "find_test" {
  database   = "test"
  collection = "test"
  command    = "find"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mongodb_query_shape.find_test", "database", "test"),
					resource.TestCheckResourceAttr("data.mongodb_query_shape.find_test", "collection", "test"),
					resource.TestCheckResourceAttr("data.mongodb_query_shape.find_test", "command", "find"),
					resource.TestCheckResourceAttrSet("data.mongodb_query_shape.find_test", "hash"),
				),
			},
		},
	})
}

func TestAccQueryShapeDatasourceFindWithFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "mongodb_query_shape" "find_filter" {
  database   = "test"
  collection = "test"
  command    = "find"
  filter     = jsonencode({"status" : "active", "category" : "electronics"})
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mongodb_query_shape.find_filter", "database", "test"),
					resource.TestCheckResourceAttr("data.mongodb_query_shape.find_filter", "collection", "test"),
					resource.TestCheckResourceAttr("data.mongodb_query_shape.find_filter", "command", "find"),
					resource.TestCheckResourceAttrSet("data.mongodb_query_shape.find_filter", "hash"),
				),
			},
		},
	})
}

func TestAccQueryShapeDatasourceFindWithSortAndProjection(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "mongodb_query_shape" "find_sort_proj" {
  database   = "test"
  collection = "test"
  command    = "find"
  sort       = jsonencode({"created_at" : -1, "name" : 1})
  projection = jsonencode({"_id" : 0, "name" : 1, "price" : 1})
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mongodb_query_shape.find_sort_proj", "database", "test"),
					resource.TestCheckResourceAttr("data.mongodb_query_shape.find_sort_proj", "command", "find"),
					resource.TestCheckResourceAttrSet("data.mongodb_query_shape.find_sort_proj", "hash"),
				),
			},
		},
	})
}

func TestAccQueryShapeDatasourceFindWithSkipAndLimit(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "mongodb_query_shape" "find_pagination" {
  database   = "test"
  collection = "test"
  command    = "find"
  skip       = 10
  limit      = 20
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mongodb_query_shape.find_pagination", "database", "test"),
					resource.TestCheckResourceAttr("data.mongodb_query_shape.find_pagination", "command", "find"),
					resource.TestCheckResourceAttr("data.mongodb_query_shape.find_pagination", "skip", "10"),
					resource.TestCheckResourceAttr("data.mongodb_query_shape.find_pagination", "limit", "20"),
					resource.TestCheckResourceAttrSet("data.mongodb_query_shape.find_pagination", "hash"),
				),
			},
		},
	})
}

func TestAccQueryShapeDatasourceAggregate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "mongodb_query_shape" "agg_test" {
  database   = "test"
  collection = "test"
  command    = "aggregate"
  pipeline   = jsonencode([{"$match" : {"status" : "active"}}, {"$group" : {"_id" : "$category", "count" : {"$sum" : 1}}}])
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mongodb_query_shape.agg_test", "database", "test"),
					resource.TestCheckResourceAttr("data.mongodb_query_shape.agg_test", "collection", "test"),
					resource.TestCheckResourceAttr("data.mongodb_query_shape.agg_test", "command", "aggregate"),
					resource.TestCheckResourceAttrSet("data.mongodb_query_shape.agg_test", "hash"),
				),
			},
		},
	})
}

func TestAccQueryShapeDatasourceAggregateWithAllowDiskUse(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "mongodb_query_shape" "agg_disk_use" {
  database      = "test"
  collection    = "test"
  command       = "aggregate"
  pipeline      = jsonencode([{"$sort" : {"created_at" : -1}}, {"$limit" : 100}])
  allow_disk_use = true
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mongodb_query_shape.agg_disk_use", "database", "test"),
					resource.TestCheckResourceAttr("data.mongodb_query_shape.agg_disk_use", "command", "aggregate"),
					resource.TestCheckResourceAttr("data.mongodb_query_shape.agg_disk_use", "allow_disk_use", "true"),
					resource.TestCheckResourceAttrSet("data.mongodb_query_shape.agg_disk_use", "hash"),
				),
			},
		},
	})
}

func TestAccQueryShapeDatasourceDistinct(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "mongodb_query_shape" "distinct_test" {
  database   = "test"
  collection = "test"
  command    = "distinct"
  key        = "status"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mongodb_query_shape.distinct_test", "database", "test"),
					resource.TestCheckResourceAttr("data.mongodb_query_shape.distinct_test", "collection", "test"),
					resource.TestCheckResourceAttr("data.mongodb_query_shape.distinct_test", "command", "distinct"),
					resource.TestCheckResourceAttr("data.mongodb_query_shape.distinct_test", "key", "status"),
					resource.TestCheckResourceAttrSet("data.mongodb_query_shape.distinct_test", "hash"),
				),
			},
		},
	})
}

func TestAccQueryShapeDatasourceDistinctWithFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "mongodb_query_shape" "distinct_filter" {
  database   = "test"
  collection = "test"
  command    = "distinct"
  key        = "category"
  filter     = jsonencode({"active" : true})
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mongodb_query_shape.distinct_filter", "database", "test"),
					resource.TestCheckResourceAttr("data.mongodb_query_shape.distinct_filter", "command", "distinct"),
					resource.TestCheckResourceAttr("data.mongodb_query_shape.distinct_filter", "key", "category"),
					resource.TestCheckResourceAttrSet("data.mongodb_query_shape.distinct_filter", "hash"),
				),
			},
		},
	})
}

func TestAccQueryShapeDatasourceFindWithCollation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "mongodb_query_shape" "find_collation" {
  database   = "test"
  collection = "test"
  command    = "find"
  collation  = jsonencode({"locale" : "en", "strength" : 2})
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mongodb_query_shape.find_collation", "database", "test"),
					resource.TestCheckResourceAttr("data.mongodb_query_shape.find_collation", "command", "find"),
					resource.TestCheckResourceAttrSet("data.mongodb_query_shape.find_collation", "hash"),
				),
			},
		},
	})
}

func TestAccQueryShapeDatasourceFindWithBatchSize(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "mongodb_query_shape" "find_batch" {
  database   = "test"
  collection = "test"
  command    = "find"
  batch_size = 100
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mongodb_query_shape.find_batch", "database", "test"),
					resource.TestCheckResourceAttr("data.mongodb_query_shape.find_batch", "command", "find"),
					resource.TestCheckResourceAttr("data.mongodb_query_shape.find_batch", "batch_size", "100"),
					resource.TestCheckResourceAttrSet("data.mongodb_query_shape.find_batch", "hash"),
				),
			},
		},
	})
}

func TestAccQueryShapeDatasourceFindWithHint(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "mongodb_query_shape" "find_hint" {
  database   = "test"
  collection = "test"
  command    = "find"
  filter     = jsonencode({"status" : "active"})
  hint       = jsonencode("_id_")
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mongodb_query_shape.find_hint", "database", "test"),
					resource.TestCheckResourceAttr("data.mongodb_query_shape.find_hint", "command", "find"),
					resource.TestCheckResourceAttrSet("data.mongodb_query_shape.find_hint", "hash"),
				),
			},
		},
	})
}

func TestAccQueryShapeDatasourceAggregateComplexPipeline(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "mongodb_query_shape" "agg_complex" {
  database   = "test"
  collection = "test"
  command    = "aggregate"
  pipeline   = jsonencode([
    {"$match" : {"status" : "active", "deleted" : {"$ne" : true}}},
    {"$group" : {"_id" : "$category", "total" : {"$sum" : "$amount"}, "count" : {"$sum" : 1}}},
    {"$sort" : {"total" : -1}},
    {"$limit" : 10},
    {"$project" : {"_id" : 1, "total" : 1, "count" : 1}}
  ])
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mongodb_query_shape.agg_complex", "database", "test"),
					resource.TestCheckResourceAttr("data.mongodb_query_shape.agg_complex", "command", "aggregate"),
					resource.TestCheckResourceAttrSet("data.mongodb_query_shape.agg_complex", "hash"),
				),
			},
		},
	})
}

func TestAccQueryShapeDatasourceHashNotEmpty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "mongodb_query_shape" "hash_check" {
  database   = "test"
  collection = "test"
  command    = "find"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.mongodb_query_shape.hash_check", "hash"),
				),
			},
		},
	})
}

func TestAccQueryShapeDatasourceFindWithPipelineShouldError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "mongodb_query_shape" "find_pipeline" {
  database   = "test"
  collection = "test"
  command    = "find"
  pipeline   = jsonencode([{"$match" : {"x" : 1}}])
}
`,
				ExpectError: regexp.MustCompile(`.*pipeline cannot be set when command is "find".*`),
			},
		},
	})
}

func TestAccQueryShapeDatasourceFindWithKeyShouldError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "mongodb_query_shape" "find_key" {
  database   = "test"
  collection = "test"
  command    = "find"
  key        = "status"
}
`,
				ExpectError: regexp.MustCompile(`.*key cannot be set when command is "find".*`),
			},
		},
	})
}

func TestAccQueryShapeDatasourceAggregateWithFilterShouldError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "mongodb_query_shape" "agg_filter" {
  database   = "test"
  collection = "test"
  command    = "aggregate"
  filter     = jsonencode({"x" : 1})
  pipeline   = jsonencode([{"$match" : {"x" : 1}}])
}
`,
				ExpectError: regexp.MustCompile(`.*filter cannot be set when command is "aggregate".*`),
			},
		},
	})
}

func TestAccQueryShapeDatasourceAggregateWithSortShouldError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "mongodb_query_shape" "agg_sort" {
  database   = "test"
  collection = "test"
  command    = "aggregate"
  sort       = jsonencode({"x" : 1})
  pipeline   = jsonencode([{"$match" : {"x" : 1}}])
}
`,
				ExpectError: regexp.MustCompile(`.*sort cannot be set when command is "aggregate".*`),
			},
		},
	})
}

func TestAccQueryShapeDatasourceAggregateWithSkipShouldError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "mongodb_query_shape" "agg_skip" {
  database   = "test"
  collection = "test"
  command    = "aggregate"
  skip       = 10
  pipeline   = jsonencode([{"$match" : {"x" : 1}}])
}
`,
				ExpectError: regexp.MustCompile(`.*skip cannot be set when command is "aggregate".*`),
			},
		},
	})
}

func TestAccQueryShapeDatasourceAggregateWithKeyShouldError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "mongodb_query_shape" "agg_key" {
  database   = "test"
  collection = "test"
  command    = "aggregate"
  key        = "status"
  pipeline   = jsonencode([{"$match" : {"x" : 1}}])
}
`,
				ExpectError: regexp.MustCompile(`.*key cannot be set when command is "aggregate".*`),
			},
		},
	})
}

func TestAccQueryShapeDatasourceDistinctWithPipelineShouldError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "mongodb_query_shape" "distinct_pipeline" {
  database   = "test"
  collection = "test"
  command    = "distinct"
  key        = "status"
  pipeline   = jsonencode([{"$match" : {"x" : 1}}])
}
`,
				ExpectError: regexp.MustCompile(`.*pipeline cannot be set when command is "distinct".*`),
			},
		},
	})
}

func TestAccQueryShapeDatasourceDistinctWithAllowDiskUseShouldError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "mongodb_query_shape" "distinct_disk" {
  database       = "test"
  collection     = "test"
  command        = "distinct"
  key            = "status"
  allow_disk_use = true
}
`,
				ExpectError: regexp.MustCompile(`.*allow_disk_use cannot be set when command is "distinct".*`),
			},
		},
	})
}

func TestAccQueryShapeDatasourceDistinctWithLimitShouldError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "mongodb_query_shape" "distinct_limit" {
  database   = "test"
  collection = "test"
  command    = "distinct"
  key        = "status"
  limit      = 10
}
`,
				ExpectError: regexp.MustCompile(`.*limit cannot be set when command is "distinct".*`),
			},
		},
	})
}

func TestAccQueryShapeDatasourceCountWithPipelineShouldError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "mongodb_query_shape" "count_pipeline" {
  database   = "test"
  collection = "test"
  command    = "count"
  pipeline   = jsonencode([{"$match" : {"x" : 1}}])
}
`,
				ExpectError: regexp.MustCompile(`.*pipeline cannot be set when command is "count".*`),
			},
		},
	})
}

func TestAccQueryShapeDatasourceCountWithKeyShouldError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "mongodb_query_shape" "count_key" {
  database   = "test"
  collection = "test"
  command    = "count"
  key        = "status"
}
`,
				ExpectError: regexp.MustCompile(`.*key cannot be set when command is "count".*`),
			},
		},
	})
}

func TestAccQueryShapeDatasourceCountWithSortShouldError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "mongodb_query_shape" "count_sort" {
  database   = "test"
  collection = "test"
  command    = "count"
  sort       = jsonencode({"x" : 1})
}
`,
				ExpectError: regexp.MustCompile(`.*sort cannot be set when command is "count".*`),
			},
		},
	})
}

func TestAccQueryShapeDatasourceHashConsistencySameConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "mongodb_query_shape" "a" {
  database   = "test"
  collection = "test"
  command    = "find"
  filter     = jsonencode({"status" : "active", "category" : "electronics", "nested" : {"a" : 1, "b" : 2}})
  sort       = jsonencode({"created_at" : -1, "name" : 1})
  projection = jsonencode({"_id" : 0, "name" : 1, "price" : 1})
}

data "mongodb_query_shape" "b" {
  database   = "test"
  collection = "test"
  command    = "find"
  filter     = jsonencode({"status" : "active", "category" : "electronics", "nested" : {"a" : 1, "b" : 2}})
  sort       = jsonencode({"created_at" : -1, "name" : 1})
  projection = jsonencode({"_id" : 0, "name" : 1, "price" : 1})
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.mongodb_query_shape.a", "hash",
						"data.mongodb_query_shape.b", "hash",
					),
				),
			},
		},
	})
}

func TestAccQueryShapeDatasourceHashConsistencyCrossStep(t *testing.T) {
	var hash string

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "mongodb_query_shape" "consistency" {
  database   = "test"
  collection = "test"
  command    = "find"
  filter     = jsonencode({"status" : "active", "category" : "electronics", "nested" : {"a" : 1, "b" : 2}})
  sort       = jsonencode({"created_at" : -1, "name" : 1})
  projection = jsonencode({"_id" : 0, "name" : 1, "price" : 1})
  hint       = jsonencode("_id_")
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.mongodb_query_shape.consistency", "hash"),
					resource.TestCheckResourceAttrWith("data.mongodb_query_shape.consistency", "hash", func(v string) error {
						hash = v
						return nil
					}),
				),
			},
			{
				Config: providerConfig + `
data "mongodb_query_shape" "consistency" {
  database   = "test"
  collection = "test"
  command    = "find"
  filter     = jsonencode({"status" : "active", "category" : "electronics", "nested" : {"a" : 1, "b" : 2}})
  sort       = jsonencode({"created_at" : -1, "name" : 1})
  projection = jsonencode({"_id" : 0, "name" : 1, "price" : 1})
  hint       = jsonencode("_id_")
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith("data.mongodb_query_shape.consistency", "hash", func(v string) error {
						if v != hash {
							return fmt.Errorf("query hash is not deterministic across applies: first produced %q, second produced %q", hash, v)
						}
						return nil
					}),
				),
			},
		},
	})
}

func TestAccQueryShapeDatasourceFindWithAllowDiskUseShouldError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "mongodb_query_shape" "find_disk" {
  database       = "test"
  collection     = "test"
  command        = "find"
  allow_disk_use = true
}
`,
				ExpectError: regexp.MustCompile(`.*allow_disk_use cannot be set when command is "find".*`),
			},
		},
	})
}