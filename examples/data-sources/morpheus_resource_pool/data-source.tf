data "morpheus_resource_pool" "morpheus_pool" {
  name     = "morpheuspool"
  cloud_id = data.morpheus_cloud.vspherecloud.id
}