package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccUserResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccUserResourceConfig(
					`"CAN_SEE_ALL_APPS", "CAN_REPLY_TO_REVIEWS_GLOBAL"`,
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"googleplay_user.oliver",
						tfjsonpath.New("email"),
						knownvalue.StringExact("example@oliverbinns.co.uk"),
					),
					statecheck.ExpectKnownValue(
						"googleplay_user.oliver",
						tfjsonpath.New("name"),
						knownvalue.StringExact("developers/5166846112789481453/users/example@oliverbinns.co.uk"),
					),
					statecheck.ExpectKnownValue(
						"googleplay_user.oliver",
						tfjsonpath.New("permissions"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("CAN_SEE_ALL_APPS"),
							knownvalue.StringExact("CAN_REPLY_TO_REVIEWS_GLOBAL"),
						}),
					),
				},
			},
			// Update and Read testing
			{
				Config: testAccUserResourceConfig(
					`"CAN_VIEW_NON_FINANCIAL_DATA_GLOBAL"`,
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"googleplay_user.oliver",
						tfjsonpath.New("email"),
						knownvalue.StringExact("example@oliverbinns.co.uk"),
					),
					statecheck.ExpectKnownValue(
						"googleplay_user.oliver",
						tfjsonpath.New("name"),
						knownvalue.StringExact("developers/5166846112789481453/users/example@oliverbinns.co.uk"),
					),
					statecheck.ExpectKnownValue(
						"googleplay_user.oliver",
						tfjsonpath.New("permissions"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("CAN_VIEW_NON_FINANCIAL_DATA_GLOBAL"),
						}),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccUserResourceConfig(permissions string) string {
	return fmt.Sprintf(`
resource "googleplay_user" "oliver" {
  email = "example@oliverbinns.co.uk"
  permissions = [
    %s
  ]
}

provider "googleplay" {
  service_account_json_base64 = filebase64("~/service-account.json")
  developer_id = "5166846112789481453"
}`, permissions)
}
