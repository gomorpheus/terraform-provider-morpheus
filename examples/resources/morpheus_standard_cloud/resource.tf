data "morpheus_tenant" "tf_example_tenant" {
  name = "Terraform Example Tenant"
}

resource "morpheus_standard_cloud" "tf_example_standard_cloud" {
  name                                    = "tf_example_standard_cloud"
  code                                    = "tfstandard"
  location                                = "denver"
  visibility                              = "private"
  tenant_id                               = data.morpheus_tenant.tf_example_tenant.id
  enabled                                 = true
  automatically_power_on_vms              = true
  import_existing_vms                     = true
  enable_network_interface_type_selection = true
  appliance_url                           = "https://demo.morpheusdata.com"
  time_zone                               = "America/Denver"
  datacenter_id                           = "12345"
  guidance                                = "manual"
  costing                                 = "costing"
  agent_install_mode                      = "cloudInit"
}