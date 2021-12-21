resource "morpheus_max_cores_policy" "tf_example_max_cores_policy_global" {
  name        = "tf_example_max_cores_policy_global"
  description = "Terraform example Morpheus max cores policy"
  enabled     = true
  max_cores   = 35
  scope       = "global"
}