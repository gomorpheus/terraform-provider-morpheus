resource "morpheus_cloud_formation_spec_template" "tfexample_cloud_formation_spec_template_url" {
  name                   = "tf_cloud_formation_spec_example_url"
  source_type            = "url"
  spec_path              = "http://example.com/spec.yaml"
  capability_iam         = true
  capability_named_iam   = true
  capability_auto_expand = true
}