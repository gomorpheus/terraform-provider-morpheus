resource "morpheus_resource_pool_group" "tfexample_resource_pool_group" {
  name              = "TFExample Resource Pool Group"
  description       = "TFExample Resource Pool Group"
  mode              = "roundRobin"
  resource_pool_ids = [1, 2, 3]
  all_group_access  = true
  group_access {
    group_id = 2
    default  = true
  }
  visibility = "public"
  tenant_ids = [1, 2]
}