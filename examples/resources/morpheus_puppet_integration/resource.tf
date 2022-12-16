resource "morpheus_puppet_integration" "tf_example_puppet_integration" {
  name                       = "tfexample puppet integration"
  enabled                    = true
  puppet_master_hostname     = "peserver01.morpheusdata.com"
  allow_immediate_execution  = true
  puppet_master_ssh_username = "root"
  puppet_master_ssh_password = "password123"
}