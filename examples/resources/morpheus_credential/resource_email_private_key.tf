resource "morpheus_credential" "tf_example_credential_email_private_key" {
  name        = "tf_example_credential_email_private_key"
  description = "terraform credential example for email private key"
  enabled     = true
  type        = "email-private-key"
  email       = "test@example.local"
  key_pair_id = 33
}