resource "morpheus_delayed_delete_policy" "tf_example_delayed_delete_policy_global" {
  name        = "tf_example_delayed_delete_policy_global"
  description = "terraform example global delayed delete policy"
  enabled     = true
  delete_days = 7
  scope       = "global"
}