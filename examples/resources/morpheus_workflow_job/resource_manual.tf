data "morpheus_workflow" "example_workflow" {
  name = "Deploy app"
}

resource "morpheus_workflow_job" "tf_example_workflow_job_date_and_time" {
  name           = "TF Example Workflow Job Manual"
  enabled        = true
  labels         = ["aws", "demo"]
  workflow_id    = data.morpheus_workflow.example_workflow.id
  schedule_mode  = "manual"
  context_type   = "instance-label"
  instance_label = "demo"
  custom_options = {
    "demo" = "testing"
  }
}