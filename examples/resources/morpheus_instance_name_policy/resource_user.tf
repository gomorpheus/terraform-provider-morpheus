resource "morpheus_instance_name_policy" "tf_example_instance_name_policy_user" {
  name                   = "tf_example_instance_name_policy_user"
  description            = "terraform example user instance name policy"
  enabled                = true
  enforcement_type       = "fixed"
  naming_pattern         = "$${userInitials.toLowerCase()}dm$${type.take(3).toLowerCase()}$${sequence+1000}"
  auto_resolve_conflicts = true
  scope                  = "user"
  user_id                = 1
}