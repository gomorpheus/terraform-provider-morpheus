resource "morpheus_cluster_resource_name_policy" "tf_example_cluster_resource_name_policy_global" {
  name                   = "tf_example_cluster_resource_name_policy_global"
  description            = "terraform example global cluster resource name policy"
  enabled                = true
  enforcement_type       = "fixed"
  naming_pattern         = "$${userInitials.toUpperCase()}DMCLSTR$${sequence+1000}"
  auto_resolve_conflicts = true
  scope                  = "global"
}