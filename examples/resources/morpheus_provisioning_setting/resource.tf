resource "morpheus_provisioning_setting" "tf_example_provisioning_setting" {
  allow_zone_selection         = false
  allow_host_selection         = false
  require_environments         = false
  show_pricing                 = true
  hide_datastore_stats         = true
  cross_tenant_naming_policies = false
  cloudinit_username           = "cloudinit"
  cloudinit_password           = "Pa55w0rd!"
  windows_password             = "Pa55w0rd!"
  pxe_root_password            = "Pa55w0rd!"
}
