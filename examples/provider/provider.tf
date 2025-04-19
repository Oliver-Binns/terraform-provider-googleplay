provider "googleplay" {
  service_account_json_base64 = filebase64("service_account.json")
  developer_id                = "5166846112789481453"
}
