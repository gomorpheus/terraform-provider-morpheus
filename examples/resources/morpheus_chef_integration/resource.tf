resource "morpheus_chef_integration" "tf_example_chef_integration" {
  name                       = "tfexample chef integration"
  enabled                    = true
  url                        = "https://chef.morpheusdata.com"
  version                    = "15.9.38"
  windows_version            = "15.9.38"
  windows_msi_install_url    = "https://packages.chef.io"
  organization               = "morpheus"
  username                   = "admin"
  private_key                = "EXAMPLEPRIVATEKEY"
  organization_validator_key = "EXAMPLEPRIVATEKEY"
}