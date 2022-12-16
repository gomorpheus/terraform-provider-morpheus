resource "morpheus_tag_policy" "tf_example_tag_policy_cloud" {
  name               = "tf_example_tag_policy_cloud"
  description        = "terraform example cloud tag policy"
  enabled            = true
  strict_enforcement = true
  tag_key            = "cost_center"
  option_list_id     = 23
  scope              = "cloud"
  cloud_id           = 1
}