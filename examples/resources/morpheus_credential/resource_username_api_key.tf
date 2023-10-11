resource "morpheus_credential" "tf_example_credential_username_api_key" {
  name        = "tf_example_credential_username_api_key"
  description = "terraform credential example for username api key"
  enabled     = true
  type        = "username-api-key"
  username    = "admin"
  api_key     = "MFIEIWEIFINEF"
}