resource "morpheus_user_group" "tf_example_user_group" {
  name         = "tftest"
  description  = "terraform"
  sudo_access  = true
  server_group = "test"
  user_ids     = [19, 10]
}