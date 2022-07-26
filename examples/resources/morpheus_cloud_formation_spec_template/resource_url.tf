resource "morpheus_cloud_formation_spec_template" "tfexample_cloud_formation_spec_template_url" {
  name        = "tf-cloud-formation-spec-example-url"
  source_type = "url"
  spec_path   = "http://example.com/spec.yaml"
}