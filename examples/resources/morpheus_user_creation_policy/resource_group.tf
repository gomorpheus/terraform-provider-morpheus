resource "morpheus_user_creation_policy" "tf_example_user_creation_policy_group" {
  name             = "tf_example_user_creation_policy_group"
  description      = "terraform example group user creation policy"
  enabled          = true
  enforcement_type = "fixed"
  create_user      = true
  scope            = "group"
  group_id         = 1
}