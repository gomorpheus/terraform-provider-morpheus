resource "morpheus_cloud_formation_spec_template" "tfexample_cloud_formation_spec_template_local" {
  name         = "tf-cloud-formation-spec-example-local"
  category     = ""
  source_type  = "local"
  spec_content = <<TFEOF

TFEOF

}