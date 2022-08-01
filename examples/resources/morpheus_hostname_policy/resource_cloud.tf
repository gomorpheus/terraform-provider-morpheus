resource "morpheus_hostname_policy" "tf_example_hostname_policy_cloud" {
  name             = "tf_example_hostname_policy_cloud"
  description      = "terraform example cloud hostname policy"
  enabled          = true
  enforcement_type = "fixed"
  naming_pattern   = "$${userInitials.toLowerCase()}dm$${type.take(3).toLowerCase()}$${sequence+1000}"
  scope            = "cloud"
  cloud_id         = 1
}