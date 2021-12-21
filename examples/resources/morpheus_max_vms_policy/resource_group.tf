resource "morpheus_max_vms_policy" "tf_example_max_vms_policy_group" {
  name        = "tf_example_max_vms_policy_group"
  description = "Terraform example Morpheus max vms policy"
  enabled     = true
  max_vms     = 35
  scope       = "group"
  group_id    = 1
}