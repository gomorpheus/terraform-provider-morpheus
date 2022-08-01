resource "morpheus_max_storage_policy" "tf_example_max_storage_policy_cloud" {
  name        = "tf_example_max_storage_policy_cloud"
  description = "terraform example cloud max storage policy"
  enabled     = true
  max_storage = 100
  scope       = "cloud"
  cloud_id    = 1
}