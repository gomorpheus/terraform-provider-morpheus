resource "morpheus_operational_workflow" "operationalworkflowdemo" {
  name = "operationalworkflowdemo"
  description = "testhing"
  platform = "all"
  visibility = "private"
  allow_custom_config = true
  option_types = [
    1730
  ]
  task {
    task_id = 18
    task_phase = "operation"
  }
}