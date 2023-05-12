data "morpheus_vro_workflow" "tf_example_vro_workflow" {
  name = "My vRO Workflow Name"
}

resource "morpheus_vro_task" "tf_example_vro_task" {
  name               = "tfexample vro-task"
  code               = "tfexample-vro-task"
  labels             = ["demo", "terraform"]
  vro_integration_id = morhpeus_vro_integration.tf_example_vro_integration.id
  vro_workflow_value = data.morpheus_vro_workflow.tf_example_vro_workflow.value
  body               = <<EOF
{
    "parameters": [
        {
            "name": "vmName",
            "type": "string",
            "value": {
                "string": {
                    "value": "<%=instance.hostname%>"
                }
            }
        }
    ]
}
EOF
  execute_target     = "local"
  retryable          = false
}
