resource "morpheus_credential" "tf_example_credential_client_id_secret" {
  name          = "tf_example_credential_client_id_secret"
  description   = "terraform credential example for client id secret"
  enabled       = true
  type          = "client-id-secret"
  client_id     = "FIEFMIQNQ"
  client_secret = "MMEWMIFINWEINFINE"
}