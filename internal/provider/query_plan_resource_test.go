package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQueryPlanResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "mongodb_query_plan" "test" {
  query_hash = "0000000000000000000000000000000000000000000000000000000000000000"
  database   = "test"
  collection = "test"
  allowed_indexes = [
    "_id_"
  ]
  comment = "terraform acceptance test"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mongodb_query_plan.test", "query_hash", "0000000000000000000000000000000000000000000000000000000000000000"),
					resource.TestCheckResourceAttr("mongodb_query_plan.test", "database", "test"),
					resource.TestCheckResourceAttr("mongodb_query_plan.test", "collection", "test"),
					resource.TestCheckResourceAttr("mongodb_query_plan.test", "allowed_indexes.0", "_id_"),
					resource.TestCheckResourceAttr("mongodb_query_plan.test", "comment", "terraform acceptance test"),
				),
			},
			{
				Config: providerConfig + `
resource "mongodb_query_plan" "test" {
  query_hash = "0000000000000000000000000000000000000000000000000000000000000000"
  database   = "test"
  collection = "test"
  allowed_indexes = [
    "_id_"
  ]
  comment = "updated comment"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mongodb_query_plan.test", "query_hash", "0000000000000000000000000000000000000000000000000000000000000000"),
					resource.TestCheckResourceAttr("mongodb_query_plan.test", "database", "test"),
					resource.TestCheckResourceAttr("mongodb_query_plan.test", "collection", "test"),
					resource.TestCheckResourceAttr("mongodb_query_plan.test", "allowed_indexes.0", "_id_"),
					resource.TestCheckResourceAttr("mongodb_query_plan.test", "comment", "updated comment"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccQueryPlanResourceWildcard(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "mongodb_query_plan" "wildcard" {
  query_hash = "1111111111111111111111111111111111111111111111111111111111111111"
  database   = "test"
  collection = "test"
  allowed_indexes = ["*"]
  comment = "allow all indexes"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mongodb_query_plan.wildcard", "query_hash", "1111111111111111111111111111111111111111111111111111111111111111"),
					resource.TestCheckResourceAttr("mongodb_query_plan.wildcard", "database", "test"),
					resource.TestCheckResourceAttr("mongodb_query_plan.wildcard", "collection", "test"),
					resource.TestCheckResourceAttr("mongodb_query_plan.wildcard", "allowed_indexes.0", "*"),
				),
			},
		},
	})
}
