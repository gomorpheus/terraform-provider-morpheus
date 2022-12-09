resource "morpheus_tag_policy" "tf_example_tag_policy_user" {
  name               = "tf_example_tag_policy_user"
  description        = "terraform example user tag policy"
  enabled            = true
  strict_enforcement = true
  tag_key            = "cost_center"
  tag_value          = "true"
  option_list_id     = 2
  scope              = "user"
  user_id            = 1
}