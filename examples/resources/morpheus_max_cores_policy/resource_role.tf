resource "morpheus_max_cores_policy" "tf_example_max_cores_policy_role" {
  name            = "tf_example_max_cores_policy_role"
  description     = "Terraform example Morpheus max cores policy"
  enabled         = true
  max_cores       = 35
  scope           = "role"
  group_id        = 1
  apply_each_user = true
}