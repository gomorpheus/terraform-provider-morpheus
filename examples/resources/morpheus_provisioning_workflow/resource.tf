resource "morpheus_provisioning_workflow" "provisionworkflowdemo" {
  name = "tfdemo"
  description = "testhing"
  platform = "all"
  visibility = "private"
  task {
    task_id = 18
    task_phase = "configure"
  }
}