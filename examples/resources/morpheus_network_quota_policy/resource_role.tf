resource "morpheus_network_quota_policy" "tf_example_network_quota_policy_role" {
  name            = "tf_example_network_quota_policy_role"
  description     = "terraform example role network quota policy"
  enabled         = true
  max_networks    = 10
  scope           = "role"
  role_id         = 1
  apply_each_user = true
}