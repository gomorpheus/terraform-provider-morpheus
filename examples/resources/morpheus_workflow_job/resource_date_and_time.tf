data "morpheus_workflow" "example_workflow" {
  name = "Deploy app"
}

resource "morpheus_workflow_job" "tf_example_workflow_job_date_and_time" {
  name                    = "TF Example Workflow Job Date and Time"
  enabled                 = true
  labels                  = ["aws", "demo"]
  workflow_id             = data.morpheus_workflow.example_workflow.id
  schedule_mode           = "date_and_time"
  scheduled_date_and_time = "2022-12-30T06:00:00Z"
  context_type            = "instance"
  instance_ids            = [1, 2]
  custom_options = {
    "demo" = "testing"
  }
}