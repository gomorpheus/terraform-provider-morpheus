data "morpheus_task" "example_task" {
  name = "Deploy app"
}

resource "morpheus_task_job" "tf_example_task_job_manual" {
  name           = "TF Example Task Job Manual"
  enabled        = true
  labels         = ["aws", "demo"]
  task_id        = data.morpheus_task.example_task.id
  schedule_mode  = "manual"
  context_type   = "instance-label"
  instance_label = "demo"
}