resource "morpheus_workflow_policy" "tf_example_workflow_policy_cloud" {
  name        = "tf_example_workflow_policy_cloud"
  description = "TF Example Workflow Policy"
  enabled     = true
  workflow_id = 1
  scope       = "cloud"
  cloud_id    = 1
}