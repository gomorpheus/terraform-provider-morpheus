resource "morpheus_max_vms_policy" "tf_example_max_vms_policy_user" {
  name        = "tf_example_max_vms_policy_user"
  description = "Terraform example Morpheus max vms policy"
  enabled     = true
  max_vms     = 35
  scope       = "user"
  user_id     = 1
}