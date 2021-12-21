resource "morpheus_max_hosts_policy" "tf_example_max_hosts_policy_global" {
  name        = "tf_example_max_hosts_policy_global"
  description = "Terraform example Morpheus max hosts policy"
  enabled     = true
  max_hosts   = 35
  scope       = "global"
}