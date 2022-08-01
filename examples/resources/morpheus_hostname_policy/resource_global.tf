resource "morpheus_hostname_policy" "tf_example_hostname_policy_global" {
  name             = "tf_example_hostname_policy_global"
  description      = "terraform example global hostname policy"
  enabled          = true
  enforcement_type = "fixed"
  naming_pattern   = "$${userInitials.toLowerCase()}dm$${type.take(3).toLowerCase()}$${sequence+1000}"
  scope            = "global"
}