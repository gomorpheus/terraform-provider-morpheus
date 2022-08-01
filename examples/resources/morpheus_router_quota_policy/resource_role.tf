resource "morpheus_router_quota_policy" "tf_example_router_quota_policy_role" {
  name            = "tf_example_router_quota_policy_role"
  description     = "terraform example role router quota policy"
  enabled         = true
  max_routers     = 20
  scope           = "role"
  role_id         = 1
  apply_each_user = true
}