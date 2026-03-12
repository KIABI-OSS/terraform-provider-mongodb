package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCollectionResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "mongodb_database" "test" {
	name = "test_db"
}

resource "mongodb_collection" "test" {
	database = mongodb_database.test.name
	name = "test"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mongodb_database.test", "name", "test_db"),
					resource.TestCheckResourceAttr("mongodb_collection.test", "name", "test"),
					resource.TestCheckResourceAttr("mongodb_collection.test", "database", "test_db"),
				),
			},
		},
	})
}
