resource "morpheus_helm_spec_template" "tfexample_helm_spec_template_url" {
  name        = "tf-helm-spec-example-url"
  source_type = "url"
  spec_path   = "http://example.com/chart.yaml"
}