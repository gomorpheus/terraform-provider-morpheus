resource "morpheus_max_storage_policy" "tf_example_max_storage_policy_user" {
  name        = "tf_example_max_storage_policy_user"
  description = "terraform example user max storage policy"
  enabled     = true
  max_storage = 100
  scope       = "user"
  user_id     = 1
}