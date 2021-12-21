resource "morpheus_max_hosts_policy" "tf_example_max_hosts_policy_cloud" {
  name        = "tf_example_max_hosts_policy_cloud"
  description = "Terraform example Morpheus max hosts policy"
  enabled     = true
  max_hosts   = 35
  scope       = "cloud"
  cloud_id    = 1
}