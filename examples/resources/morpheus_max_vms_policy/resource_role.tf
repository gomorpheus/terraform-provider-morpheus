resource "morpheus_max_vms_policy" "tf_example_max_vms_policy_role" {
  name            = "tf_example_max_vms_policy_role"
  description     = "Terraform example Morpheus max vms policy"
  enabled         = true
  max_vms         = 35
  scope           = "role"
  group_id        = 1
  apply_each_user = true
}