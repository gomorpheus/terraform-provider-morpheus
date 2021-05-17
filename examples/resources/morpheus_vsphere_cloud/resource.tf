resource "morpheus_vsphere_cloud" "tf_example_vsphere_cloud" {
  name       = "tf_example_vsphere_cloud"
  code       = "tfvsphere"
  api_url    = "https://vcenter.morpheus.local/sdk"
  username   = "administrator@vsphere.local"
  password   = "password"
  datacenter = "morpheusdc"
  cluster    = "morpheus-cluster"
  rpc_mode   = "guestexec"
  location   = "denver"
}