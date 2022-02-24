data "morpheus_tenant_role" "example" {
  name = "Tenant Admin"
}

resource "morpheus_tenant" "tf_example_tenant" {
  name            = "tftenant"
  description     = "Terraform example tenant"
  enabled         = true
  subdomain       = "tfexample"
  base_role_id    = data.morpheus_tenant_role.example.id
  currency        = "USD"
  account_number  = "12345"
  account_name    = "tenant 12345"
  customer_number = "12345"
}