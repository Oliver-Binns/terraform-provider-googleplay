package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccUserDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccExampleDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.googleplay_user.oliver",
						tfjsonpath.New("name"),
						knownvalue.StringExact("Oliver Binns"),
					),
					statecheck.ExpectKnownValue(
						"data.googleplay_user.oliver",
						tfjsonpath.New("email"),
						knownvalue.StringExact("mail@oliverbinns.co.uk"),
					),
				},
			},
		},
	})
}

const testAccExampleDataSourceConfig = `
data "googleplay_user" "oliver" {
  email = "mail@oliverbinns.co.uk"
  name = "Oliver Binns"
}

provider "googleplay" {
  service_account_json_base64 = filebase64("~/service-account.json")
  developer_id = "5166846112789481453"
}
`
