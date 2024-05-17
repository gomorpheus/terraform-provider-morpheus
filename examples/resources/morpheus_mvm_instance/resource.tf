resource "morpheus_mvm_instance" "test" {
  name               = "demo"
  cloud_id           = 4
  group_id           = 4
  plan_id            = 2
  instance_layout_id = 2
  instance_type_id   = 4
}