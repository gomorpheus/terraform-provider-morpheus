resource "morpheus_hostname_policy" "tf_example_hostname_policy_user" {
  name             = "tf_example_hostname_policy_user"
  description      = "terraform example user hostname policy"
  enabled          = true
  enforcement_type = "fixed"
  naming_pattern   = "$${userInitials.toLowerCase()}dm$${type.take(3).toLowerCase()}$${sequence+1000}"
  scope            = "user"
  user_id          = 1
}