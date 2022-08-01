resource "morpheus_service_plan" "tf_example_service_plan" {
  name           = "terraform-test-sp"
  code           = "terraform-test-sp1"
  active         = true
  display_order  = 2
  provision_type = "vmware"
  // Processors
  max_cores        = 8
  custom_cores     = true
  cores_per_socket = 4

  // Memory
  max_memory       = 3145728
  memory_size_type = "mb"
  custom_memory    = true
  custom_memory_range {
    minimum = 1048576
    maximum = 3145728
  }
  // Storage
  max_storage             = 3221225472
  storage_size_type       = "mb"
  customize_root_volume   = true
  customize_extra_volumes = true
  add_volumes             = true
  max_disks_allowed       = 0
  custom_storage_range {
    minimum = 3000
    maximum = 5000
  }

  price_set_ids = [morpheus_price_set.tf_example_price_set_software.id,
    203,
    645,
  202]
}