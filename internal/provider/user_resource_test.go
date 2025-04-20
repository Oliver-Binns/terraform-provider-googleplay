package provider

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccUserResource(t *testing.T) {
	accountEmail := fmt.Sprintf(
		"%s@oliverbinns.co.uk",
		uuid.New().String(),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccUserResourceConfig(
					accountEmail,
					`"CAN_VIEW_NON_FINANCIAL_DATA_GLOBAL", "CAN_REPLY_TO_REVIEWS_GLOBAL"`,
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"googleplay_user.oliver",
						tfjsonpath.New("email"),
						knownvalue.StringExact(accountEmail),
					),
					statecheck.ExpectKnownValue(
						"googleplay_user.oliver",
						tfjsonpath.New("name"),
						knownvalue.StringExact(
							fmt.Sprintf(
								"developers/5166846112789481453/users/%s",
								accountEmail,
							),
						),
					),
					statecheck.ExpectKnownValue(
						"googleplay_user.oliver",
						tfjsonpath.New("global_permissions"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("CAN_REPLY_TO_REVIEWS_GLOBAL"),
							knownvalue.StringExact("CAN_VIEW_NON_FINANCIAL_DATA_GLOBAL"),
						}),
					),
					statecheck.ExpectKnownValue(
						"googleplay_user.oliver",
						tfjsonpath.New("app_permissions"),
						knownvalue.ListExact([]knownvalue.Check{}),
					),
				},
			},
			// Update and Read testing
			{
				Config: testAccUserResourceConfig(
					accountEmail,
					`"CAN_MANAGE_TRACK_USERS_GLOBAL"`,
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"googleplay_user.oliver",
						tfjsonpath.New("email"),
						knownvalue.StringExact(accountEmail),
					),
					statecheck.ExpectKnownValue(
						"googleplay_user.oliver",
						tfjsonpath.New("name"),
						knownvalue.StringExact(
							fmt.Sprintf(
								"developers/5166846112789481453/users/%s",
								accountEmail,
							),
						),
					),
					statecheck.ExpectKnownValue(
						"googleplay_user.oliver",
						tfjsonpath.New("global_permissions"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("CAN_MANAGE_TRACK_USERS_GLOBAL"),
						}),
					),
					statecheck.ExpectKnownValue(
						"googleplay_user.oliver",
						tfjsonpath.New("app_permissions"),
						knownvalue.ListExact([]knownvalue.Check{}),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccUserResourceConfig(accountEmail string, permissions string) string {
	return fmt.Sprintf(`
resource "googleplay_user" "oliver" {
  email = "%s"
  global_permissions = [
    %s
  ]
  app_permissions = []
}

provider "googleplay" {
  service_account_json_base64 = filebase64("~/service-account.json")
  developer_id = "5166846112789481453"
}`, accountEmail, permissions)
}
