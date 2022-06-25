data "morpheus_tenant" "tf_example_tenant" {
  name = "Terraform Example Tenant"
}

resource "morpheus_vsphere_cloud" "tf_example_vsphere_cloud" {
  name                                    = "tf_example_vsphere_cloud"
  code                                    = "tfvsphere"
  location                                = "denver"
  visibility                              = "private"
  tenant_id                               = data.morpheus_tenant.tf_example_tenant.id
  enabled                                 = true
  automatically_power_on_vms              = true
  api_url                                 = "https://vcenter.morpheus.local/sdk"
  username                                = "administrator@vsphere.local"
  password                                = "password"
  api_version                             = "6.7"
  datacenter                              = "morpheusdc"
  cluster                                 = "morpheus-cluster"
  resource_pool                           = ""
  rpc_mode                                = "guestexec"
  hide_host_selection                     = true
  import_existing_vms                     = true
  enable_hypervisor_console               = true
  keyboard_layout                         = "us"
  enable_disk_type_selection              = true
  enable_storage_type_selection           = true
  enable_network_interface_type_selection = true
  storage_type                            = "thin"
  appliance_url                           = "https://demo.morpheusdata.com"
  time_zone                               = "America/Denver"
  datacenter_id                           = "12345"
  guidance                                = "manual"
  costing                                 = "costing"
  agent_install_mode                      = "cloudInit"
}