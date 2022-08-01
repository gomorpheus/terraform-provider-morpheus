resource "morpheus_max_memory_policy" "tf_example_max_memory_policy_group" {
  name        = "tf_example_max_memory_policy_group"
  description = "terraform example group max memory policy"
  enabled     = true
  max_memory  = 256
  scope       = "group"
  group_id    = 1
}