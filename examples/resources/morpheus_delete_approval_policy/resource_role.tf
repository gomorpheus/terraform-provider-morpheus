resource "morpheus_delete_approval_policy" "tf_example_delete_approval_policy_global" {
  name            = "tf_example_delete_approval_policy_global"
  description     = "terraform example global delete approval policy"
  enabled         = true
  integration_id  = 1
  workflow_id     = 10
  scope           = "role"
  role_id         = 1
  apply_each_user = true
}