resource "morpheus_workflow_catalog_item" "tfexample_workflow_catalog_item" {
  name         = "tfexample_workflow_catalog_item"
  description  = "Example Terraform workflow catalog item"
  enabled      = true
  featured     = true
  workflow_id  = 1
  context_type = "appliance"
  content      = <<TFEOF
Testing
TFEOF
}