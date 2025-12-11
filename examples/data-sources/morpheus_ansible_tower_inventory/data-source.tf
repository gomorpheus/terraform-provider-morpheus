data "morpheus_integration" "tf_example_ansible_tower_integration" {
  name = "Demo Ansible Tower"
}

data "morpheus_ansible_tower_inventory" "example_ansible_tower_inventory" {
  ansible_tower_integration_id = data.morpheus_integration.tf_example_ansible_tower_integration.id
  name                         = "Demo Inventory"
}