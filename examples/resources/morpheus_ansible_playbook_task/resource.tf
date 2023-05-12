resource "morpheus_ansible_playbook_task" "ansplaybook" {
  name                = "tfansibletest"
  code                = "tfansibletest"
  labels              = ["demo", "terraform"]
  ansible_repo_id     = "5"
  git_ref             = "master"
  playbook            = "mongo_install"
  tags                = "mongo"
  skip_tags           = "web"
  command_options     = "-b"
  execute_target      = "local"
  retryable           = true
  retry_count         = 1
  retry_delay_seconds = 10
  allow_custom_config = true
}