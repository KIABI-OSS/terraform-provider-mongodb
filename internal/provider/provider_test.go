package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	providerConfig = `
provider "mongodb" {
  host = "localhost"
  port = "27017"
  username = "test"
  password = "test"
}
`
	providerConfigWithURL = `
provider "mongodb" {
  url = "mongodb://localhost:27017"
}
`
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"mongodb": providerserver.NewProtocol6WithError(New("test")()),
	}
)

func TestMongodbProvider_Configure(t *testing.T) {
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig,
			},
		},
	})
}

func TestMongodbProvider_Configure_WithURL(t *testing.T) {
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfigWithURL,
			},
		},
	})
}
