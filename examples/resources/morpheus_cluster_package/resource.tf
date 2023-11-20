resource "morpheus_cluster_package" "tf_example_cluster_package" {
  name              = "tf_example_cluster_package"
  code              = "tf-example-cluster-package"
  description       = "Terraform example cluster package"
  package_version   = "1.2.3"
  type              = "apps"
  package_type      = "example"
  enabled           = true
  repeat_install    = true
  spec_template_ids = [1, 2]
}