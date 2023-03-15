data "morpheus_workflow" "example_workflow" {
  name = "Example Workflow"
}

resource "morpheus_nested_workflow_task" "tfexample_nested_workflow" {
  name                      = "tfexample_nested_workflow"
  code                      = "tfexample_nested_workflow"
  operational_workflow_id   = data.morpheus_workflow.example_workflow.id
  operational_workflow_name = data.morpheus_workflow.example_workflow.name
}