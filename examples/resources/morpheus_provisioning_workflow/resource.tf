resource "morpheus_provisioning_workflow" "tf_example_provisioning_workflow" {
  name        = "tf_example_provisioning_workflow"
  description = "Terraform provisioning workflow example"
  platform    = "all"
  visibility  = "private"
  task {
    task_id = 18
    task_phase = "configure"
  }
}