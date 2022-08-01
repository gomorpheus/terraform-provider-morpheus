resource "morpheus_network_quota_policy" "tf_example_network_quota_policy_global" {
  name         = "tf_example_network_quota_policy_global"
  description  = "terraform example global network quota policy"
  enabled      = true
  max_networks = 10
  scope        = "global"
}