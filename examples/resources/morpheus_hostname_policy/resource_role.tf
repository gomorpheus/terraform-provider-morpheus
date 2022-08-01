resource "morpheus_hostname_policy" "tf_example_hostname_policy_role" {
  name             = "tf_example_hostname_policy_role"
  description      = "terraform example role hostname policy"
  enabled          = true
  enforcement_type = "fixed"
  naming_pattern   = "$${userInitials.toLowerCase()}dm$${type.take(3).toLowerCase()}$${sequence+1000}"
  scope            = "role"
  role_id          = 1
  apply_each_user  = true
}