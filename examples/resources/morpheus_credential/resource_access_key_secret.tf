resource "morpheus_credential" "tf_example_credential_access_key_secret" {
  name        = "tf_example_credential_access_key_secret"
  description = "terraform credential example for access key and secret key"
  enabled     = true
  type        = "access-key-secret"
  access_key  = "FIEFMIQNQ"
  secret_key  = "MFMWEIIEIFENF"
}