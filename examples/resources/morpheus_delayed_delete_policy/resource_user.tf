resource "morpheus_delayed_delete_policy" "tf_example_delayed_delete_policy_user" {
  name        = "tf_example_delayed_delete_policy_user"
  description = "terraform example user delayed delete policy"
  enabled     = true
  delete_days = 7
  scope       = "user"
  user_id     = 1
}