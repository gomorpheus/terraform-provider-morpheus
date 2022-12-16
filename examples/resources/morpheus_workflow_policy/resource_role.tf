resource "morpheus_workflow_policy" "tf_example_workflow_policy_role" {
  name            = "tf_example_workflow_policy_role"
  description     = "TF Example Workflow Policy"
  enabled         = true
  workflow_id     = 1
  scope           = "role"
  role_id         = 1
  apply_each_user = true
}