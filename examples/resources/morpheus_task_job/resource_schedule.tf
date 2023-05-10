data "morpheus_task" "example_task" {
  name = "Deploy app"
}

data "morpheus_execute_schedule" "example_schedule" {
  name = "Run Daily at 9 AM"
}

resource "morpheus_task_job" "tf_example_task_job_schedule" {
  name                  = "TF Example Task Job Schedule"
  enabled               = true
  labels                = ["aws", "demo"]
  task_id               = data.morpheus_task.example_task.id
  schedule_mode         = "scheduled"
  execution_schedule_id = data.morpheus_execute_schedule.jobtest.id
  context_type          = "instance"
  instance_ids          = [91]
  custom_config         = "{\"test\":\"new\"}"
}