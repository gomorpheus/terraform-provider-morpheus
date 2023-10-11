resource "morpheus_credential" "tf_example_credential_api_key" {
  name        = "tf_example_credential_api_key"
  description = "terraform credential example for api key"
  enabled     = true
  type        = "api-key"
  api_key     = "FIEFMIQNQ"
}