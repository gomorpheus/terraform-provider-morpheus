data "morpheus_tenant" "demo_tenant" {
  name = "Demo"
}

resource "morpheus_saml_identity_source" "addemo" {
  tenant_id                      = morpheus_tenant.demo_tenant.id
  name                           = "samldemo"
  description                    = "TF example SAML identity source"
  login_redirect_url             = "https://tfexamplesaml.test.local:8443/realms/master/protocol/saml"
  logout_redirect_url            = "https://tfexamplesaml.test.local:8443/realms/master/protocol/saml"
  include_saml_request_parameter = true
  saml_request                   = "SelfSigned"
  validate_assertion_signature   = false
  given_name_attribute           = "givenName"
  surname_attribute              = "surname"
  email_attribute                = "email"
  default_account_role_id        = 4
  role_attribute_name            = "memberOf"
  required_role_attribute_value  = "test"
  role_mapping {
    role_id             = 4
    role_name           = "Demo"
    assertion_attribute = "developers"
  }

  role_mapping {
    role_id             = 5
    role_name           = "tf-example-user-role"
    assertion_attribute = "developers"
  }
  enable_role_mapping_permission = false
}