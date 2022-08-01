resource "morpheus_router_quota_policy" "tf_example_router_quota_policy_user" {
  name        = "tf_example_router_quota_policy_user"
  description = "terraform example user router quota policy"
  enabled     = true
  max_routers = 20
  scope       = "user"
  user_id     = 1
}