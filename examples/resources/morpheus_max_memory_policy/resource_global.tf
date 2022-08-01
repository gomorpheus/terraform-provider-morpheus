resource "morpheus_max_memory_policy" "tf_example_max_memory_policy_global" {
  name        = "tf_example_max_memory_policy_global"
  description = "terraform example global max memory policy"
  enabled     = true
  max_memory  = 256
  scope       = "global"
}