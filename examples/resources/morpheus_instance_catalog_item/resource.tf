resource "morpheus_instance_catalog_item" "tf_example_instance_catalog_item" {
  name        = "tfexample_instance_catalog"
  description = "terraform example instance catalog item"
  enabled     = true
  feature     = true
  content     = <<TFEOF
  {"name":"test"}
  TFEOF
  config      = <<TFEOF
  {"name":"test"}
  TFEOF
  visibility  = "private"
}