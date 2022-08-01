resource "morpheus_cluster_resource_name_policy" "tf_example_cluster_resource_name_policy_role" {
  name                   = "tf_example_cluster_resource_name_policy_role"
  description            = "terraform example role cluster resource name policy"
  enabled                = true
  enforcement_type       = "fixed"
  naming_pattern         = "$${userInitials.toUpperCase()}DMCLSTR$${sequence+1000}"
  auto_resolve_conflicts = true
  scope                  = "role"
  role_id                = 1
  apply_each_user        = true
}