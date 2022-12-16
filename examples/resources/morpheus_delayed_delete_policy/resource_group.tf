resource "morpheus_delayed_delete_policy" "tf_example_delayed_delete_policy_group" {
  name        = "tf_example_delayed_delete_policy_group"
  description = "terraform example group delayed delete policy"
  enabled     = true
  delete_days = 7
  scope       = "group"
  group_id    = 1
}