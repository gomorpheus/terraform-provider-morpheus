resource "morpheus_network_quota_policy" "tf_example_network_quota_policy_cloud" {
  name         = "tf_example_network_quota_policy_cloud"
  description  = "terraform example cloud network quota policy"
  enabled      = true
  max_networks = 10
  scope        = "cloud"
  cloud_id     = 1
}