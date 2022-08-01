resource "morpheus_user_creation_policy" "tf_example_user_creation_policy_user" {
  name             = "tf_example_user_creation_policy_user"
  description      = "terraform example user user creation policy"
  enabled          = true
  enforcement_type = "fixed"
  create_user      = true
  scope            = "user"
  user_id          = 1
}