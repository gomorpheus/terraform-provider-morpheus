resource "morpheus_tag_policy" "tf_example_tag_policy_global" {
  name               = "tf_example_tag_policy_global"
  description        = "terraform example global tag policy"
  enabled            = true
  strict_enforcement = true
  tag_key            = "cost_center"
  option_list_id     = 23
  scope              = "global"
}