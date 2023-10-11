resource "morpheus_credential" "tf_example_credential_tenant_username_keypair" {
  name        = "tf_example_credential_tenant_username_keypair"
  description = "terraform credential example for tenant username keypair"
  enabled     = true
  type        = "tenant-username-keypair"
  tenant      = "tenant123"
  username    = "admin"
  key_pair_id = 22
}