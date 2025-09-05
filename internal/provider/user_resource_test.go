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
			// Create and read testing
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
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("CAN_REPLY_TO_REVIEWS_GLOBAL"),
							knownvalue.StringExact("CAN_VIEW_NON_FINANCIAL_DATA_GLOBAL"),
						}),
					),
				},
			},
			// ImportState testing
			{
				ResourceName:      "googleplay_user.oliver",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"global_permissions",
				},
			},
			// Test update global permissions
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
						tfjsonpath.New("expanded_permissions"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("CAN_MANAGE_TRACK_USERS_GLOBAL"),
						}),
					),
				},
			},
			// Test update global permissions with implicit grant
			{
				Config: testAccUserResourceConfig(
					accountEmail,
					`"CAN_MANAGE_PERMISSIONS_GLOBAL"`,
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
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("CAN_MANAGE_PERMISSIONS_GLOBAL"),
						}),
					),
					statecheck.ExpectKnownValue(
						"googleplay_user.oliver",
						tfjsonpath.New("expanded_permissions"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("CAN_VIEW_FINANCIAL_DATA_GLOBAL"),
							knownvalue.StringExact("CAN_MANAGE_PERMISSIONS_GLOBAL"),
							knownvalue.StringExact("CAN_EDIT_GAMES_GLOBAL"),
							knownvalue.StringExact("CAN_PUBLISH_GAMES_GLOBAL"),
							knownvalue.StringExact("CAN_REPLY_TO_REVIEWS_GLOBAL"),
							knownvalue.StringExact("CAN_MANAGE_PUBLIC_APKS_GLOBAL"),
							knownvalue.StringExact("CAN_MANAGE_TRACK_APKS_GLOBAL"),
							knownvalue.StringExact("CAN_MANAGE_TRACK_USERS_GLOBAL"),
							knownvalue.StringExact("CAN_MANAGE_PUBLIC_LISTING_GLOBAL"),
							knownvalue.StringExact("CAN_MANAGE_DRAFT_APPS_GLOBAL"),
							knownvalue.StringExact("CAN_CREATE_MANAGED_PLAY_APPS_GLOBAL"),
							knownvalue.StringExact("CAN_MANAGE_ORDERS_GLOBAL"),
							knownvalue.StringExact("CAN_MANAGE_APP_CONTENT_GLOBAL"),
							knownvalue.StringExact("CAN_VIEW_NON_FINANCIAL_DATA_GLOBAL"),
							knownvalue.StringExact("CAN_VIEW_APP_QUALITY_GLOBAL"),
							knownvalue.StringExact("CAN_MANAGE_DEEPLINKS_GLOBAL"),
						}),
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
}

provider "googleplay" {
  service_account_json_base64 = filebase64("~/service-account.json")
  developer_id = "5166846112789481453"
}`, accountEmail, permissions)
}
