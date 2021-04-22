resource "morpheus_network_domain" "tf_example_network_domain" {
  name        = "tfexampledomain"
  description = "Terraform example network domain"
  public_zone = true
  visibility  = "private"
  tenant_id   = 1
  active      = true
}