resource "morpheus_hostname_policy" "tf_example_hostname_policy_group" {
  name             = "tf_example_hostname_policy_group"
  description      = "terraform example group hostname policy"
  enabled          = true
  enforcement_type = "fixed"
  naming_pattern   = "$${userInitials.toLowerCase()}dm$${type.take(3).toLowerCase()}$${sequence+1000}"
  scope            = "group"
  group_id         = 1
}