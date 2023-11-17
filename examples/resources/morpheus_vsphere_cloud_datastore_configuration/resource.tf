resource "morpheus_vsphere_cloud_datastore_configuration" "tf_example_datastore" {
  cloud_id         = 2
  name             = "Example_Datastore"
  active           = true
  group_access_all = true
  group_access     = [1, 2]
  visibility       = "public"
  tenant_access {
    id           = 1
    default      = true
    image_target = false
  }

  tenant_access {
    id           = 2
    default      = true
    image_target = true
  }
}