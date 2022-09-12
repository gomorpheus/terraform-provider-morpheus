data "morpheus_cluster_type" "kubernetes_layout" {
  name = "Kubernetes Cluster"
}

data "morpheus_provision_type" "provision_layout" {
  name = "VMware"
}

resource "morpheus_cluster_layout" "example_kubernetes_layout" {
  name              = "tfexample cluster layout"
  description       = "Terraform example cluster layout"
  version           = "1.0"
  creatable         = false
  minimum_memory    = 4294967296
  workflow_id       = 2
  cluster_type_id   = data.morpheus_cluster_type.kubernetes_layout.id
  provision_type_id = data.morpheus_provision_type.provision_layout.id
  enable_scaling    = false
  option_type_ids = [
    1910,
    2037
  ]

  evar {
    name   = "application"
    value  = "first"
    export = true
  }

  master_node_pool {
    count          = 1
    node_type_id   = 3
    priority_order = 0
  }

  worker_node_pool {
    count          = 4
    node_type_id   = 4
    priority_order = 1
  }

  worker_node_pool {
    count          = 4
    node_type_id   = 3
    priority_order = 2
  }
}