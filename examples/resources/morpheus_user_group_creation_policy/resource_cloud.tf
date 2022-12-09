resource "morpheus_user_group_creation_policy" "tf_example_user_group_creation_policy_cloud" {
  name          = "tf_example_user_group_creation_policy_cloud"
  description   = "terraform example cloud user group creation policy"
  enabled       = true
  user_group_id = 1
  scope         = "cloud"
  cloud_id      = 1
}