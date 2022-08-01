resource "morpheus_max_containers_policy" "tf_example_max_containers_policy_cloud" {
  name           = "tf_example_max_containers_policy_cloud"
  description    = "terraform example cloud max containers policy"
  enabled        = true
  max_containers = 50
  scope          = "cloud"
  cloud_id       = 1
}