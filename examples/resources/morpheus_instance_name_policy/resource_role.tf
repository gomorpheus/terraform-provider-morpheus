resource "morpheus_instance_name_policy" "tf_example_instance_name_policy_role" {
  name                   = "tf_example_instance_name_policy_role"
  description            = "terraform example role instance name policy"
  enabled                = true
  enforcement_type       = "fixed"
  naming_pattern         = "$${userInitials.toLowerCase()}dm$${type.take(3).toLowerCase()}$${sequence+1000}"
  auto_resolve_conflicts = true
  scope                  = "role"
  role_id                = 1
  apply_each_user        = true
}