data "morpheus_group" "morpheus_lab" {
  name = "MORPHEUS"
}

data "morpheus_cloud" "morpheus_aws" {
  name = "MORPHEUSAWS"
}

data "morpheus_resource_pool" "aws_resource_pool" {
  name     = "Morpheus-VPC (us-east-2)"
  cloud_id = data.morpheus_cloud.morpheus_aws.id
}

data "morpheus_instance_type" "ubuntu" {
  name = "Ubuntu"
}

data "morpheus_instance_layout" "ubuntu" {
  name    = "Amazon VM"
  version = "22.04"
}

data "morpheus_network" "vmnetwork" {
  name = "Morpheus-Subnet (subnet-0ed95648b7e27a375)"
}

data "morpheus_plan" "aws" {
  name = "T3 Small - 2 Core, 2GB Memory"
}

resource "morpheus_aws_instance" "tf_example_aws_instance" {
  name               = "tfaws"
  description        = "Terraform instance example"
  cloud_id           = data.morpheus_cloud.morpheus_aws.id
  group_id           = data.morpheus_group.morpheus_lab.id
  instance_type_id   = data.morpheus_instance_type.ubuntu.id
  instance_layout_id = data.morpheus_instance_layout.ubuntu.id
  plan_id            = data.morpheus_plan.aws.id
  environment        = "dev"
  resource_pool_id   = data.morpheus_resource_pool.aws_resource_pool.id
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