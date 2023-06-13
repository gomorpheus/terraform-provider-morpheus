data "morpheus_integration" "servicenow_prod" {
  name = "SNOW Production"
}

data "morpheus_servicenow_workflow" "morpheus_example" {
  name           = "Knowledge - Retire Knowledge"
  integration_id = data.morpheus_integration.servicenow_prod.id
}

resource "morpheus_delete_approval_policy" "tf_example_delete_approval_policy_global" {
  name           = "tf_example_delete_approval_policy_global"
  description    = "terraform example global delete approval policy"
  enabled        = true
  integration_id = data.morpheus_integration.servicenow_prod.id
  workflow_id    = data.morpheus_servicenow_workflow.morpheus_example.id
  scope          = "global"
}