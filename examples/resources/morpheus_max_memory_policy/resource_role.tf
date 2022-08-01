resource "morpheus_max_memory_policy" "tf_example_max_memory_policy_role" {
  name            = "tf_example_max_memory_policy_role"
  description     = "terraform example role max memory policy"
  enabled         = true
  max_memory      = 256
  scope           = "role"
  role_id         = 1
  apply_each_user = true
}