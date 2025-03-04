// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

const name = "Oliver Binns"
const email = "mail@oliverbinns.co.uk"

func TestAccUserResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccUserResourceConfig(name, email),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"googleplay_user.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(name),
					),
					statecheck.ExpectKnownValue(
						"googleplay_user.test",
						tfjsonpath.New("email"),
						knownvalue.StringExact(email),
					),
				},
			},
			// ImportState testing
			{
				ResourceName:                         "googleplay_user.test",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "email",
			},
			// Update and Read testing
			{
				Config: testAccUserResourceConfig(name, "test@oliverbinns.co.uk"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"googleplay_user.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(name),
					),
					statecheck.ExpectKnownValue(
						"googleplay_user.test",
						tfjsonpath.New("email"),
						knownvalue.StringExact("test@oliverbinns.co.uk"),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccUserResourceConfig(name string, email string) string {
	return fmt.Sprintf(`
resource "googleplay_user" "test" {
  name  = %[1]q
  email = %[2]q
}

provider "googleplay" {
  service_account_json = "{}"
}
`, name, email)
}
