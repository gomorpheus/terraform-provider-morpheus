resource "morpheus_user_group_creation_policy" "tf_example_user_group_creation_policy_global" {
  name             = "tf_example_user_group_creation_policy_global"
  description      = "terraform example global user group creation policy"
  enabled          = true
  user_group_id = 1
  scope            = "global"
}