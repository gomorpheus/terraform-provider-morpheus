resource "morpheus_max_cores_policy" "tf_example_max_cores_policy_cloud" {
  name        = "tf_example_max_cores_policy_cloud"
  description = "Terraform example Morpheus max cores policy"
  enabled     = true
  max_cores   = 35
  scope       = "cloud"
  cloud_id    = 1
}