resource "morpheus_max_vms_policy" "tf_example_max_vms_policy_global" {
  name        = "tf_example_max_vms_policy_global"
  description = "Terraform example Morpheus max vms policy"
  enabled     = true
  max_vms     = 35
  scope       = "global"
}