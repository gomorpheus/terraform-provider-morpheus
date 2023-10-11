resource "morpheus_credential" "tf_example_credential_username_password" {
  name        = "tf_example"
  description = "terraform example"
  enabled     = true
  type        = "username-password"
  username    = "admin"
  password    = "password12333"
}