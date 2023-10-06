data "morpheus_credential" "azure_credentials" {
  name = "azuredemo"
}

resource "morpheus_azure_cloud" "tf_example_azure_cloud" {
  name                       = "tf-azure-demo"
  code                       = "tf-azure-demo"
  location                   = "colorado"
  visibility                 = "public"
  tenant_id                  = 1
  enabled                    = true
  automatically_power_on_vms = true
  cloud_type                 = "global"
  azure_subscription_id      = "12345"
  azure_tenant_id            = "12345"
  credential_id              = data.morpheus_credential.azure_credentials.id
  region                     = "centralus"
  resource_group             = "all"
  import_existing_instances  = true
  rpc_mode                   = "guestexec"
  appliance_url              = "https://morpheus.local"
  time_zone                  = "America/Denver"
  datacenter_id              = "tfazuredemo"
  guidance                   = "manual"
  costing                    = "full"
  agent_install_mode         = "cloudInit"
}