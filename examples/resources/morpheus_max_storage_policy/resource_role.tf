resource "morpheus_max_storage_policy" "tf_example_max_storage_policy_role" {
  name            = "tf_example_max_storage_policy_role"
  description     = "terraform example role max storage policy"
  enabled         = true
  max_storage     = 100
  scope           = "role"
  role_id         = 1
  apply_each_user = true
}