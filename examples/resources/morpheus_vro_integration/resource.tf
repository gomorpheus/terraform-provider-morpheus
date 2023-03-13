resource "morpheus_vro_integration" "tf_example_vro_integration" {
  name      = "tfexample vro"
  enabled   = true
  url       = "https://myvro/vco/api"
  username  = "my-vro-username"
  password  = "my-vro-password"
  auth_type = "basic"
  tenant    = "vsphere.local"
}
