resource "morpheus_guidance_setting" "tf_example_guidance_setting" {
  power_settings_average_cpu          = 75
  power_settings_maximum_cpu          = 500
  power_settings_network_threshold    = 2000
  cpu_upsize_average_cpu              = 50
  cpu_upsize_maximum_cpu              = 99
  memory_upsize_minimum_free_memory   = 10
  memory_downsize_average_free_memory = 60
  memory_downsize_maximum_free_memory = 30
}