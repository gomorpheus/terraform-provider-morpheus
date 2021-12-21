resource "morpheus_max_hosts_policy" "tf_example_max_hosts_policy_user" {
  name        = "tf_example_max_hosts_policy_user"
  description = "Terraform example Morpheus max hosts policy"
  enabled     = true
  max_hosts   = 35
  scope       = "user"
  user_id     = 1
}