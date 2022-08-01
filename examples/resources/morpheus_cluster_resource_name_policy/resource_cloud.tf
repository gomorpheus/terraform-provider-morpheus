resource "morpheus_cluster_resource_name_policy" "tf_example_cluster_resource_name_policy_cloud" {
  name                   = "tf_example_cluster_resource_name_policy_cloud"
  description            = "tfvsphere"
  enabled                = true
  enforcement_type       = "fixed"
  naming_pattern         = "$${userInitials.toUpperCase()}DMCLSTR$${sequence+1000}"
  auto_resolve_conflicts = true
  scope                  = "cloud"
  cloud_id               = 1
}