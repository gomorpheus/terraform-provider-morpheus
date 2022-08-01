resource "morpheus_max_storage_policy" "tf_example_max_storage_policy_global" {
  name        = "tf_example_max_storage_policy_global"
  description = "terraform example global max storage policy"
  enabled     = true
  max_storage = 100
  scope       = "global"
}