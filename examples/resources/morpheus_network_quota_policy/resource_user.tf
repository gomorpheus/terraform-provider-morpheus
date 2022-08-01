resource "morpheus_network_quota_policy" "tf_example_network_quota_policy_user" {
  name         = "tf_example_network_quota_policy_user"
  description  = "terraform example user network quota policy"
  enabled      = true
  max_networks = 10
  scope        = "user"
  user_id      = 1
}