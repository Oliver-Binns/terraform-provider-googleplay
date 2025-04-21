resource "googleplay_app_iam" "test_app" {
  app_id  = "0000000000000000000"
  user_id = googleplay_user.test.email
  permissions = [
    "CAN_VIEW_NON_FINANCIAL_DATA", "CAN_REPLY_TO_REVIEWS"
  ]
}