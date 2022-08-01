resource "morpheus_max_memory_policy" "tf_example_max_memory_policy_cloud" {
  name        = "tf_example_max_memory_policy_cloud"
  description = "terraform example cloud max memory policy"
  enabled     = true
  max_memory  = 256
  scope       = "cloud"
  cloud_id    = 1
}