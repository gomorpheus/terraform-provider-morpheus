resource "morpheus_max_storage_policy" "tf_example_max_storage_policy_group" {
  name        = "tf_example_max_storage_policy_group"
  description = "terraform example group max storage policy"
  enabled     = true
  max_storage = 100
  scope       = "group"
  group_id    = 1
}