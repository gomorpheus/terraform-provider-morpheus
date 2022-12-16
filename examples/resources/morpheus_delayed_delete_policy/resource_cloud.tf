resource "morpheus_delayed_delete_policy" "tf_example_delayed_delete_policy_cloud" {
  name        = "tf_example_delayed_delete_policy_cloud"
  description = "terraform example cloud delayed delete policy"
  enabled     = true
  delete_days = 7
  scope       = "cloud"
  cloud_id    = 1
}