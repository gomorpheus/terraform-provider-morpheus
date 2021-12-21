resource "morpheus_max_hosts_policy" "tf_example_max_hosts_policy_role" {
  name            = "tf_example_max_hosts_policy_role"
  description     = "Terraform example Morpheus max hosts policy"
  enabled         = true
  max_hosts       = 35
  scope           = "role"
  group_id        = 1
  apply_each_user = true
}