
data "morpheus_ansible_tower_job_template" "example_ansible_tower_job_template" {
  name = "Demo Job Template"
}

data "morpheus_integration" "tf_example_ansible_tower_integration" {
  name = "Demo Ansible Tower"
}

data "morpheus_ansible_tower_inventory" "example_ansible_tower_inventory" {
  ansible_tower_integration_id = data.morpheus_integration.tf_example_ansible_tower_integration.id
  name                         = "Demo Inventory"
}

resource "morpheus_ansible_tower_task" "tfexample_ansible_tower_task" {
  name                         = "tfexample_ansible_tower_task"
  code                         = "tfexample-ansible-tower-task"
  labels                       = ["demo", "terraform"]
  ansible_tower_integration_id = data.morpheus_integration.tf_example_ansible_tower_integration.id
  ansible_tower_inventory_id   = data.morpheus_ansible_tower_inventory.example_ansible_tower_inventory.id
  group                        = "demo"
  job_template_id              = data.morpheus_ansible_tower_job_template.example_ansible_tower_job_template.id
  scm_override                 = "main"
  execute_mode                 = "executeAll"
  execute_target               = "local"
  retryable                    = true
  retry_count                  = 5
  retry_delay_seconds          = 10
  allow_custom_config          = true
  visibility                   = "public"
}