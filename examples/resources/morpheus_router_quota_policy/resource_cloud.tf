resource "morpheus_router_quota_policy" "tf_example_router_quota_policy_cloud" {
  name        = "tf_example_router_quota_policy_cloud"
  description = "terraform example cloud router quota policy"
  enabled     = true
  max_routers = 20
  scope       = "cloud"
  cloud_id    = 1
}