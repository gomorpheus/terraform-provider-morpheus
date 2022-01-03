data "morpheus_task" "example_task" {
  name = "Deploy app"
}

resource "morpheus_task_job" "tf_example_task_job_date_and_time" {
  name                    = "TF Example Task Job Date and Time"
  enabled                 = true
  task_id                 = data.morpheus_task.example_task.id
  schedule_mode           = "date_and_time"
  scheduled_date_and_time = "2022-12-30T06:00:00Z"
  context                 = "instance"
  instance_ids            = [1, 2]
}