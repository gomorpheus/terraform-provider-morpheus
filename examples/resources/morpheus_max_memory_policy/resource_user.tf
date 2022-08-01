resource "morpheus_max_memory_policy" "tf_example_max_memory_policy_user" {
  name        = "tf_example_max_memory_policy_user"
  description = "terraform example user max memory policy"
  enabled     = true
  max_memory  = 256
  scope       = "user"
  user_id     = 1
}