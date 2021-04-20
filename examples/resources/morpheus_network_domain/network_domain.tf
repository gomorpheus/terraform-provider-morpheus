resource "morpheus_network_domain" "name" {
  name        = ""
  description = ""
  public_zone = true
  visibility  = "public"
  tenant_id   = 1
  active      = true
}