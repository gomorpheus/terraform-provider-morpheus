data "morpheus_task" "example_task" {
  name = "Deploy app"
}

resource "morpheus_task_job" "tf_example_task_job_manual" {
  name          = "TF Example Task Job Manual"
  enabled       = true
  task_id       = data.morpheus_task.example_task.id
  schedule_mode = "manual"
  context       = "instance"
  instance_ids  = [1, 2]
}