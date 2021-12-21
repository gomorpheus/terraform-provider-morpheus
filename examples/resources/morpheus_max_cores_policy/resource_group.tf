resource "morpheus_max_cores_policy" "tf_example_max_cores_policy_group" {
  name        = "tf_example_max_cores_policy_group"
  description = "Terraform example Morpheus max cores policy"
  enabled     = true
  max_cores   = 35
  scope       = "group"
  group_id    = 1
}