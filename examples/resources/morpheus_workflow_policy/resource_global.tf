resource "morpheus_workflow_policy" "tf_example_workflow_policy_global" {
  name        = "tf_example_workflow_policy_global"
  description = "tfvsphere"
  enabled     = true
  workflow_id = 1
  scope       = "global"
}