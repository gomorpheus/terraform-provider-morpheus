resource "morpheus_user_creation_policy" "tf_example_user_creation_policy_role" {
  name             = "tf_example_user_creation_policy_role"
  description      = "terraform example role user creation policy"
  enabled          = true
  enforcement_type = "fixed"
  create_user      = true
  scope            = "role"
  role_id          = 1
  apply_each_user  = true
}