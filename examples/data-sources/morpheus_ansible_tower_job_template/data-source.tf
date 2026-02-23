data "morpheus_integration" "tf_example_ansible_tower_integration" {
  name = "Demo Ansible Tower"
}

data "morpheus_ansible_tower_job_template" "example_ansible_tower_job_template" {
  ansible_tower_integration_id = data.morpheus_integration.tf_example_ansible_tower_integration.id
  name                         = "Demo Job Template"
}