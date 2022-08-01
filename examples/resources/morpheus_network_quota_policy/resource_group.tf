resource "morpheus_network_quota_policy" "tf_example_network_quota_policy_group" {
  name         = "tf_example_network_quota_policy_group"
  description  = "terraform example group network quota policy"
  enabled      = true
  max_networks = 10
  scope        = "group"
  group_id     = 1
}