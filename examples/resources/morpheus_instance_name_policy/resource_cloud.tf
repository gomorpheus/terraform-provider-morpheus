resource "morpheus_instance_name_policy" "tf_example_instance_name_policy_cloud" {
  name                   = "tf_example_instance_name_policy_cloud"
  description            = "terraform example cloud instance name policy"
  enabled                = true
  enforcement_type       = "fixed"
  naming_pattern         = "$${userInitials.toLowerCase()}dm$${type.take(3).toLowerCase()}$${sequence+1000}"
  auto_resolve_conflicts = true
  scope                  = "cloud"
  cloud_id               = 1
}