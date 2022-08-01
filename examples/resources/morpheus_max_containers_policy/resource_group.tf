resource "morpheus_max_containers_policy" "tf_example_max_containers_policy_group" {
  name           = "tf_example_max_containers_policy_group"
  description    = "terraform example group max containers policy"
  enabled        = true
  max_containers = 50
  scope          = "group"
  group_id       = 1
}