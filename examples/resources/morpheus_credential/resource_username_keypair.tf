resource "morpheus_credential" "tf_example_credential_username_keypair" {
  name        = "tf_example_credential_username_keypair"
  description = "terraform credential example for username key pair"
  enabled     = true
  type        = "username-keypair"
  username    = "admin"
  key_pair_id = 22
}