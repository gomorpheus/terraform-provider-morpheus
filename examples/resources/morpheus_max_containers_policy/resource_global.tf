resource "morpheus_max_containers_policy" "tf_example_max_containers_policy_global" {
  name           = "tf_example_max_containers_policy_global"
  description    = "terraform example global max containers policy"
  enabled        = true
  max_containers = 50
  scope          = "global"
}