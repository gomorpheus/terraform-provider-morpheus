resource "morpheus_max_containers_policy" "tf_example_max_containers_policy_user" {
  name           = "tf_example_max_containers_policy_user"
  description    = "terraform example user max containers policy"
  enabled        = true
  max_containers = 50
  scope          = "user"
  user_id        = 1
}