resource "morpheus_instance_layout" "tf_example_instance_layout" {
  instance_type_id = morpheus_instance_type.tf_example_instance_type.id
  name             = "todo_app_frontend"
  version          = "1.0"
  technology       = "vmware"
  node_type_ids = [
     morpheus_node_type.ubuntu_base.id
  ]
  workflow_id = morpheus_provisioning_workflow.tfexample_workflow.id
}