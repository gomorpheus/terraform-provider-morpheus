resource "morpheus_operational_workflow" "tf_example_operational_workflow" {
  name                = "tf_example_operational_workflow"
  description         = "Terraform operational workflow example"
  labels              = ["demo", "terraform"]
  platform            = "all"
  visibility          = "private"
  allow_custom_config = true
  option_types = [
    1730
  ]
  task_ids = [18]
}