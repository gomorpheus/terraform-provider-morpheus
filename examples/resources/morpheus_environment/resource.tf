resource "morpheus_environment" "tf_example_environment" {
  active      = true
  code        = "tfexample"
  description = "Terraform Example"
  name        = "tfexample"
}