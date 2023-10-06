data "morpheus_group" "morpheus_lab" {
  name = "MORPHEUS"
}

data "morpheus_cloud" "morpheus_vsphere" {
  name = "MORPHEUSVCENTER"
}

data "morpheus_resource_pool" "vsphere_resource_pool" {
  name     = "Morpheus-Cluster"
  cloud_id = data.morpheus_cloud.morpheus_vsphere.id
}

data "morpheus_instance_type" "ubuntu" {
  name = "Ubuntu"
}

data "morpheus_instance_layout" "ubuntu" {
  name    = "VMware VM"
  version = "22.04"
}

data "morpheus_network" "vmnetwork" {
  name = "VM Network"
}

data "morpheus_plan" "vmware" {
  name = "1 CPU, 4GB Memory"
}

resource "morpheus_vsphere_instance" "tf_example_vsphere_instance" {
  name               = "tfvsphere"
  description        = "Terraform instance example"
  cloud_id           = data.morpheus_cloud.morpheus_vsphere.id
  group_id           = data.morpheus_group.morpheus_lab.id
  instance_type_id   = data.morpheus_instance_type.ubuntu.id
  instance_layout_id = data.morpheus_instance_layout.ubuntu.id
  plan_id            = data.morpheus_plan.vmware.id
  environment        = "dev"
  resource_pool_id   = data.morpheus_resource_pool.vsphere_resource_pool.id
  labels             = ["demo", "terraform"]

  interfaces {
    network_id = data.morpheus_network.vmnetwork.id
  }

  tags = {
    name = "ubuntutf"
  }

  evar {
    name   = "application"
    value  = "demo"
    export = true
    masked = true
  }

  custom_options = {
    awsRegion = "us-east-1"
  }
}