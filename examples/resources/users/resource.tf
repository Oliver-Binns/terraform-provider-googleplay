resource "googleplay_user" "test" {
  email = "test@oliverbinns.co.uk"
  permissions = [
    "CAN_MANAGE_PERMISSIONS_GLOBAL"
  ]
}