resource "morpheus_user_group_creation_policy" "tf_example_user_group_creation_policy_role" {
  name            = "tf_example_user_group_creation_policy_role"
  description     = "terraform example role user group creation policy"
  enabled         = true
  user_group_id   = 1
  scope           = "role"
  role_id         = 1
  apply_each_user = true
}