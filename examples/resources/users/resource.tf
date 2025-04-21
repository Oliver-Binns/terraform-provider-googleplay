resource "googleplay_user" "test" {
  email = "test@oliverbinns.co.uk"
  global_permissions = [
    "CAN_MANAGE_PERMISSIONS_GLOBAL"
  ]
  app_permissions = [
    {
      app_id = "4973279986054171407"
      permissions = [
        "CAN_VIEW_APP_QUALITY"
      ]
    }
  ]
}