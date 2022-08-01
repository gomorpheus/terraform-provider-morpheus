resource "morpheus_user_creation_policy" "tf_example_user_creation_policy_global" {
  name             = "tf_example_user_creation_policy_global"
  description      = "terraform example global user creation policy"
  enabled          = true
  enforcement_type = "fixed"
  create_user      = true
  scope            = "global"
}