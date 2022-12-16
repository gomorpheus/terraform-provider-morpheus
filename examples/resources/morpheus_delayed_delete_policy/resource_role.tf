resource "morpheus_delayed_delete_policy" "tf_example_delayed_delete_policy_role" {
  name               = "tf_example_delayed_delete_policy_role"
  description        = "terraform example role delayed delete policy"
  enabled            = true
  delete_days        = 7
  scope              = "role"
  role_id            = 1
  apply_to_each_user = true
}