resource "morpheus_max_hosts_policy" "tf_example_max_hosts_policy_group" {
  name        = "tf_example_max_hosts_policy_group"
  description = "Terraform example Morpheus max hosts policy"
  enabled     = true
  max_hosts   = 35
  scope       = "group"
  group_id    = 1
}