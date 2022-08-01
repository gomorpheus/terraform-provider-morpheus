resource "morpheus_router_quota_policy" "tf_example_router_quota_policy_group" {
  name        = "tf_example_router_quota_policy_group"
  description = "terraform example group router quota policy"
  enabled     = true
  max_routers = 20
  scope       = "group"
  group_id    = 1
}