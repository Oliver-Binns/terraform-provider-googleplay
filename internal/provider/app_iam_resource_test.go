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

func TestAccAppIAMResource(t *testing.T) {
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
				Config: testAccAppIAMResourceConfig(
					accountEmail,
					`"CAN_VIEW_APP_QUALITY"`,
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"googleplay_app_iam.test_app",
						tfjsonpath.New("app_id"),
						knownvalue.StringExact("4973279986054171407"),
					),
					statecheck.ExpectKnownValue(
						"googleplay_app_iam.test_app",
						tfjsonpath.New("user_id"),
						knownvalue.StringExact(accountEmail),
					),
					statecheck.ExpectKnownValue(
						"googleplay_app_iam.test_app",
						tfjsonpath.New("permissions"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("CAN_VIEW_APP_QUALITY"),
						}),
					),
					statecheck.ExpectKnownValue(
						"googleplay_app_iam.test_app",
						tfjsonpath.New("expanded_permissions"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("CAN_VIEW_APP_QUALITY"),
						}),
					),
				},
			},
			// Test update permissions
			{
				Config: testAccAppIAMResourceConfig(
					accountEmail,
					`"CAN_VIEW_APP_QUALITY", "CAN_VIEW_NON_FINANCIAL_DATA"`,
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"googleplay_app_iam.test_app",
						tfjsonpath.New("app_id"),
						knownvalue.StringExact("4973279986054171407"),
					),
					statecheck.ExpectKnownValue(
						"googleplay_app_iam.test_app",
						tfjsonpath.New("user_id"),
						knownvalue.StringExact(accountEmail),
					),
					statecheck.ExpectKnownValue(
						"googleplay_app_iam.test_app",
						tfjsonpath.New("permissions"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("CAN_VIEW_NON_FINANCIAL_DATA"),
							knownvalue.StringExact("CAN_VIEW_APP_QUALITY"),
						}),
					),
					statecheck.ExpectKnownValue(
						"googleplay_app_iam.test_app",
						tfjsonpath.New("expanded_permissions"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("CAN_VIEW_NON_FINANCIAL_DATA"),
							knownvalue.StringExact("CAN_VIEW_APP_QUALITY"),
						}),
					),
				},
			},
			// Test update permissions with implicit grant
			{
				Config: testAccAppIAMResourceConfig(
					accountEmail,
					`"CAN_REPLY_TO_REVIEWS"`,
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"googleplay_app_iam.test_app",
						tfjsonpath.New("app_id"),
						knownvalue.StringExact("4973279986054171407"),
					),
					statecheck.ExpectKnownValue(
						"googleplay_app_iam.test_app",
						tfjsonpath.New("user_id"),
						knownvalue.StringExact(accountEmail),
					),
					statecheck.ExpectKnownValue(
						"googleplay_app_iam.test_app",
						tfjsonpath.New("permissions"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("CAN_REPLY_TO_REVIEWS"),
						}),
					),
					statecheck.ExpectKnownValue(
						"googleplay_app_iam.test_app",
						tfjsonpath.New("expanded_permissions"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("CAN_REPLY_TO_REVIEWS"),
							knownvalue.StringExact("CAN_VIEW_NON_FINANCIAL_DATA"),
							knownvalue.StringExact("CAN_VIEW_APP_QUALITY"),
						}),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccAppIAMResourceConfig(email string, permissions string) string {
	return fmt.Sprintf(`
resource "googleplay_user" "test" {
  email = "%s"
  global_permissions = [
  	"CAN_EDIT_GAMES_GLOBAL"
  ]
}

resource "googleplay_app_iam" "test_app" {
  app_id = "4973279986054171407"
  user_id = googleplay_user.test.email
  permissions = [
    %s
  ]
}

provider "googleplay" {
  developer_id = "5166846112789481453"
}`, email, permissions)
}
