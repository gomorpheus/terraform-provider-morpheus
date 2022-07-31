resource "morpheus_max_vms_policy" "tf_example_max_vms_policy_cloud" {
  name        = "tf_example_max_vms_policy_cloud"
  description = "Terraform example Morpheus max vms policy"
  enabled     = true
  max_vms     = 35
  scope       = "cloud"
  cloud_id    = 1
}