data "morpheus_group" "morpheus_lab" {
  name = "platform engineering"
}

data "morpheus_cloud" "morpheus_cloud" {
  name = "MVM Cloud"
}

data "morpheus_instance_type" "ubuntu" {
  name = "Ubuntu"
}

data "morpheus_instance_layout" "ubuntu" {
  name    = "Single KVM VM"
  version = "22.04"
}

data "morpheus_network" "vmnetwork" {
  name = "Compute"
}

data "morpheus_plan" "mvm" {
  name           = "2 CPU, 16GB Memory"
  provision_type = "KVM"
}

resource "morpheus_mvm_instance" "tf_example_mvm_instance" {
  for_each = toset([
    "tfdemo",
  ])
  name               = each.key
  description        = "Terraform instance example"
  cloud_id           = data.morpheus_cloud.morpheus_cloud.id
  group_id           = data.morpheus_group.morpheus_lab.id
  instance_type_id   = data.morpheus_instance_type.ubuntu.id
  instance_layout_id = data.morpheus_instance_layout.ubuntu.id
  plan_id            = data.morpheus_plan.mvm.id
  environment        = "dev"
  resource_pool_name = "mvmcluster01"
  labels             = ["demo", "terraform"]
  create_user        = true
  skip_agent_install = true
  workflow_name      = "Nginx Install"
  interfaces {
    network_id                = data.morpheus_network.vmnetwork.id
    network_interface_type_id = 5
  }

  volumes {
    root         = true
    size         = 30
    name         = "root"
    storage_type = 1
    datastore_id = 36
  }

  tags = {
    name = "ubuntutf"
  }
}