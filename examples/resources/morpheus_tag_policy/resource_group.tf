resource "morpheus_tag_policy" "tf_example_tag_policy_group" {
  name               = "tf_example_tag_policy_group"
  description        = "terraform example group tag policy"
  enabled            = true
  strict_enforcement = true
  tag_key            = "cost_center"
  tag_value          = "true"
  option_list_id     = 2
  scope              = "group"
  group_id           = 1
}