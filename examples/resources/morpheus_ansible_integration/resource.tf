data "morpheus_key_pair" "tf_example_key_pair" {
  name = "morpheusgit"
}

resource "morpheus_ansible_integration" "tf_example_ansible_integration" {
  name                          = "tfexample ansible"
  enabled                       = true
  url                           = "https://github.com/gomorpheus/morpheus-ansible.git"
  default_branch                = "master"
  playbooks_path                = "/"
  roles_path                    = "/roles"
  group_variables_path          = "/vars"
  host_variables_path           = "/vars"
  enable_ansible_galaxy_install = true
  enable_verbose_logging        = true
  enable_agent_command_bus      = true
  key_pair_id                   = data.morpheus_key_pair.tf_example_key_pair.id
  enable_git_caching            = true
}