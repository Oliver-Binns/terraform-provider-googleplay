# Copyright (c) HashiCorp, Inc.

provider "googleplay" {
  service_account_json = file("service_account.json")
}
