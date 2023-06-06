data "morpheus_workflow" "example_workflow" {
  name = "Deploy app"
}

data "morpheus_execute_schedule" "example_schedule" {
  name = "Run Daily at 9 AM"
}

resource "morpheus_workflow_job" "tf_example_workflow_job_date_and_time" {
  name                  = "TF Example Workflow Job Schedule"
  enabled               = true
  labels                = ["aws", "demo"]
  workflow_id           = data.morpheus_workflow.example_workflow.id
  schedule_mode         = "scheduled"
  execution_schedule_id = data.morpheus_execute_schedule.example_schedule.id
  context_type          = "instance"
  instance_ids          = [91]
  custom_options = {
    "demo" = "testing"
  }
}