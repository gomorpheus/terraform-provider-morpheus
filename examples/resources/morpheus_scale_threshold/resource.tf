resource "morpheus_scale_threshold" "tf_example_scale_threshold" {
  name                    = "example_scale_threshold"
  auto_upscale            = true
  auto_downscale          = true
  min_count               = 1
  max_count               = 3
  enable_cpu_threshold    = true
  min_cpu_percentage      = 30.0
  max_cpu_percentage      = 75.0
  enable_memory_threshold = true
  min_memory_percentage   = 20.0
  max_memory_percentage   = 60.0
  enable_disk_threshold   = true
  min_disk_percentage     = 25.0
  max_disk_percentage     = 80.0
}