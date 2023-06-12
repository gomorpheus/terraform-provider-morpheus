resource "morpheus_provision_approval_policy" "tf_example_provision_approval_policy_global" {
  name                   = "tf_example_provision_approval_policy_global"
  description            = "terraform example global provision approval policy"
  enabled                = true
  use_internal_approvals = true
  scope                  = "user"
  user_id                = 1
}