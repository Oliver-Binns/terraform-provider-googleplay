resource "googleplay_user" "test" {
  email = "test@oliverbinns.co.uk"
  global_permissions = [
    "CAN_MANAGE_PERMISSIONS_GLOBAL"
  ]
}
