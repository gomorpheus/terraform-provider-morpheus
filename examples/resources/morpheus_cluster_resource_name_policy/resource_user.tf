resource "morpheus_cluster_resource_name_policy" "tf_example_cluster_resource_name_policy_user" {
  name                   = "tf_example_cluster_resource_name_policy_user"
  description            = "terraform example user cluster resource name policy"
  enabled                = true
  enforcement_type       = "fixed"
  naming_pattern         = "$${userInitials.toUpperCase()}DMCLSTR$${sequence+1000}"
  auto_resolve_conflicts = true
  scope                  = "user"
  user_id                = 1
}