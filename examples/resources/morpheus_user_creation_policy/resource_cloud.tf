resource "morpheus_user_creation_policy" "tf_example_user_creation_policy_cloud" {
  name             = "tf_example_user_creation_policy_cloud"
  description      = "terraform example cloud user creation policy"
  enabled          = true
  enforcement_type = "fixed"
  create_user      = true
  scope            = "cloud"
  cloud_id         = 1
}