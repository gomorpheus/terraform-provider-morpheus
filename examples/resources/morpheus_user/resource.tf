resource "morpheus_user" "tf_example_user" {
  username              = "tftest"
  first_name            = "terraform"
  last_name             = "test"
  email                 = "test@test.local"
  password              = "PmWFEAE#92331"
  role_ids              = [19, 10]
  receive_notifications = true
  linux_username        = "testuser"
  linux_password        = "PmWFEAE#92331"
  windows_username      = "testuser"
  windows_password      = "PmWFEAE#92331"
}