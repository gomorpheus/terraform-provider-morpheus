data "morpheus_cloud" "vspherecloud" {
  name = "morpheus_vsphere"
}

data "morpheus_resource_pool" "morpheus_pool" {
  name     = "morpheuspool"
  cloud_id = data.morpheus_cloud.vspherecloud.id
}

resource "morpheus_price_set" "tf_example_price_set_software" {
  name             = "terraform-test-everything"
  code             = "terraform-test-everything"
  region_code      = "us-west-2"
  cloud_id         = data.morpheus_cloud.vspherecloud.id
  resource_pool_id = data.morpheus_resource_pool.morpheus_pool.id
  price_unit       = "minute"
  type             = "fixed"
  price_ids        = [morpheus_price.tf_example_price.id]
}