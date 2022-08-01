resource "morpheus_router_quota_policy" "tf_example_router_quota_policy_global" {
  name        = "tf_example_router_quota_policy_global"
  description = "terraform example global router quota policy"
  enabled     = true
  max_routers = 20
  scope       = "global"
}