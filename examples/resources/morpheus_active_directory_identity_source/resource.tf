resource "morpheus_active_directory_identity_source" "addemo" {
  tenant_id               = 1
  name                    = "addemo"
  description             = "TF example AD identity source"
  ad_server               = "dc01.contoso.com"
  domain                  = "contoso.com"
  use_ssl                 = false
  binding_username        = "administrator"
  binding_password        = "Password123"
  required_group          = "administrators"
  search_member_groups    = true
  default_account_role_id = 7

  role_mapping {
    role_id                     = 2
    role_name                   = "developers"
    active_directory_group_name = "developers"
    active_directory_group_fqn  = "CN=developers,CN=Users,DC=contoso,DC=com"
  }
}