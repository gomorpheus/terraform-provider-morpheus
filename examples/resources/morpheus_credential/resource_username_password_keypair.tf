resource "morpheus_credential" "tf_example_credential_username_password_keypair" {
  name        = "tf_example_credential_username_password_keypair"
  description = "terraform credential example for username password key pair"
  enabled     = true
  type        = "username-password-keypair"
  username    = "admin"
  password    = "password123"
  key_pair_id = 22
}