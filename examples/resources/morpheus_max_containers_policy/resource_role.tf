resource "morpheus_max_containers_policy" "tf_example_max_containers_policy_role" {
  name            = "tf_example_max_containers_policy_role"
  description     = "terraform example role max containers policy"
  enabled         = true
  max_containers  = 50
  scope           = "role"
  role_id         = 1
  apply_each_user = true
}