resource "morpheus_workflow_catalog_item" "tfexample_workflow_catalog_item" {
  name                 = "tfexample_workflow_catalog_item"
  description          = "Example Terraform workflow catalog item"
  logo_image_path      = "wordpress.png"
  logo_image_name      = "wordpress.png"
  dark_logo_image_path = "wordpressbak.png"
  dark_logo_image_name = "wordpressbak.png"
  enabled              = true
  featured             = true
  labels               = ["aws", "demo"]
  workflow_id          = 1
  option_type_ids      = [2056, 2006]
  context_type         = "appliance"
  content              = file("${path.module}/catalog-data.md")
}