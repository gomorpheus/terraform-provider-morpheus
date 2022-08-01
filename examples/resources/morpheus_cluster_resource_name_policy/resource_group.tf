resource "morpheus_cluster_resource_name_policy" "tf_example_cluster_resource_name_policy_group" {
  name                   = "tf_example_cluster_resource_name_policy_group"
  description            = "terraform example group cluster resource name policy"
  enabled                = true
  enforcement_type       = "fixed"
  naming_pattern         = "$${userInitials.toUpperCase()}DMCLSTR$${sequence+1000}"
  auto_resolve_conflicts = true
  scope                  = "group"
  group_id               = 1
}