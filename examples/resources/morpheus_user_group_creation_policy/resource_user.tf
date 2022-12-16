resource "morpheus_user_group_creation_policy" "tf_example_user_group_creation_policy_user" {
  name          = "tf_example_user_group_creation_policy_user"
  description   = "terraform example user user group creation policy"
  enabled       = true
  user_group_id = 1
  scope         = "user"
  user_id       = 1
}