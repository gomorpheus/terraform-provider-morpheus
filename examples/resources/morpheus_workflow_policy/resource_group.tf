resource "morpheus_workflow_policy" "tf_example_workflow_policy_group" {
  name        = "tf_example_workflow_policy_group"
  description = "TF Example Workflow Policy"
  enabled     = true
  workflow_id = 1
  scope       = "group"
  group_id    = 1
}