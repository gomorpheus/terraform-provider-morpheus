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
  name           = "1 CPU, 4GB Memory"
  provision_type = "vmware"
}

resource "morpheus_vsphere_mks_cluster" "tf_example_vsphere_instance" {
  name          = "tfvsphere"
  description   = "Terraform instance example"
  visibility    = "private"
  labels        = ["demo", "terraform"]
  cloud_id      = data.morpheus_cloud.morpheus_vsphere.id
  group_id      = data.morpheus_group.morpheus_lab.id

  cluster_layout_id = data.morpheus_instance_layout.ubuntu.id

  master_node_pool {
    plan_id          = data.morpheus_plan.vmware.id
    resource_pool_id = data.morpheus_resource_pool.vsphere_resource_pool.id
    host_id          = data.xxx
    folder_id        = data.xxx
    create_user      = true
    user_group_id    = 4


    network_interface {
      network_id                = data.morpheus_network.vmnetwork.id
      network_interface_type_id = 5
    }

    storage_volume {
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

  worker_node_pool {
    number_of_workers = 3
    plan_id           = data.morpheus_plan.vmware.id
    resource_pool_id  = data.morpheus_resource_pool.vsphere_resource_pool.id
    host_id           = data.xxx
    folder_id         = data.xxx
    create_user       = true
    user_group_id     = 4

    network_interface {
      network_id                = data.morpheus_network.vmnetwork.id
      network_interface_type_id = 5
    }

    storage_volume {
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
}