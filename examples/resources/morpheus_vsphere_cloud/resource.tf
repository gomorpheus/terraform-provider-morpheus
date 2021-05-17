resource "morpheus_vsphere_cloud" "morpheusvsphere" {
  name       = "tfvsphere"
  code       = "tfvsphere"
  api_url    = "https://vcenter.morpheus.local/sdk"
  username   = "administrator@vsphere.local"
  password   = "password"
  datacenter = "morpheusdc"
  cluster    = "morpheus-cluster"
  rpc_mode   = "guestexec"
  location   = "denver"
}