resource "morpheus_user_group_creation_policy" "tf_example_user_group_creation_policy_group" {
  name          = "tf_example_user_group_creation_policy_group"
  description   = "terraform example group user group creation policy"
  enabled       = true
  user_group_id = 1
  scope         = "group"
  group_id      = 1
}