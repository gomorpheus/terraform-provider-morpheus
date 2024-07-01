data "morpheus_cloud_type" "canonical_maas_cloud" {
  name = "MaaS"
}

resource "morpheus_appliance_setting" "tf_example_appliance_setting" {
  appliance_url          = "https://morpheus.test.local"
  internal_appliance_url = "https://pxemorpheus.test.local"
  api_allowed_origins = "demo"
  registration_enabled = true
  default_role_id = 5
  default_user_role_id = 10
  docker_privileged_mode = false
  smtp_from_address = "testemail@test.local"
  smtp_server = "smtp01.test.local"
  smtp_port = 465
  smtp_use_ssl = true
  smtp_use_tls = true
  smtp_username = "jsmith"
  smtp_password = "Password12"
  proxy_host = "10.0.0.100"
  proxy_port = 3128
  proxy_user = "jsmith"
  proxy_password = "Password123456"
  proxy_domain = "test.local"
  proxy_workstation = "work123"
  currency_provider = "fixer"
  currency_provider_api_key = "5a4b220e-6f9f-43da-a572-390c8f6afed8"
}