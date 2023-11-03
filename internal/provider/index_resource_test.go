package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIndexResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "mongodb_index" "acc_test" {
  database   = "test"
  collection = "test"
  name       = "tf_acc_test"
  keys = [
    {
      "field" : "field1"
      "type" : "asc"
    }
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of items
					resource.TestCheckResourceAttr("mongodb_index.acc_test", "database", "test"),
					resource.TestCheckResourceAttr("mongodb_index.acc_test", "collection", "test"),
					resource.TestCheckResourceAttr("mongodb_index.acc_test", "name", "tf_acc_test"),
					resource.TestCheckResourceAttr("mongodb_index.acc_test", "keys.0.field", "field1"),
					resource.TestCheckResourceAttr("mongodb_index.acc_test", "keys.0.type", "asc"),
					resource.TestCheckNoResourceAttr("mongodb_index.acc_test", "sparse"),
					resource.TestCheckNoResourceAttr("mongodb_index.acc_test", "expire_after_seconds"),
					resource.TestCheckNoResourceAttr("mongodb_index.acc_test", "unique"),
					resource.TestCheckNoResourceAttr("mongodb_index.acc_test", "wildcard_projection"),
				),
			},
			// Replace and Read testing
			{
				Config: providerConfig + `
resource "mongodb_index" "acc_test" {
  database   = "test"
  collection = "test"
  name       = "tf_acc_test"
  keys = [
    {
      "field" : "field1"
      "type" : "asc"
    },
	{
	  "field" : "field2"
	  "type" : "desc"
	}
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mongodb_index.acc_test", "database", "test"),
					resource.TestCheckResourceAttr("mongodb_index.acc_test", "collection", "test"),
					resource.TestCheckResourceAttr("mongodb_index.acc_test", "name", "tf_acc_test"),
					resource.TestCheckResourceAttr("mongodb_index.acc_test", "keys.0.field", "field1"),
					resource.TestCheckResourceAttr("mongodb_index.acc_test", "keys.0.type", "asc"),
					resource.TestCheckResourceAttr("mongodb_index.acc_test", "keys.1.field", "field2"),
					resource.TestCheckResourceAttr("mongodb_index.acc_test", "keys.1.type", "desc"),
					resource.TestCheckNoResourceAttr("mongodb_index.acc_test", "sparse"),
					resource.TestCheckNoResourceAttr("mongodb_index.acc_test", "expire_after_seconds"),
					resource.TestCheckNoResourceAttr("mongodb_index.acc_test", "unique"),
					resource.TestCheckNoResourceAttr("mongodb_index.acc_test", "wildcard_projection"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
