---
subcategory: ""
page_title: "Authenticate to Morpheus"
description: |-
    A guide to obtaining client credentials and adding them to provider configuration.
---

# Getting Started with the Morpheus Provider

## Before you begin

* [Install Terraform](https://www.terraform.io/intro/getting-started/install.html)
and read the Terraform getting started guide that follows. This guide will
assume basic proficiency with Terraform
* Build the Morpheus Terraform provider using the `make dev` command. 

## Configuring the provider

```terraform
terraform {
  required_providers {
    morpheus = {
      source  = "morpheus/morpheus"
      version = "~> 0.1"
    }
  }
}

# Configure the provider
provider "morpheus" {
  url      = "${var.morpheus_url}"
  username = "${var.morpheus_username}"
  password = "${var.morpheus_password}"       
}
```

## Create a Morpheus vSphere instance using the `morpheus_vsphere_instance` resource

```terraform
resource "morpheus_vsphere_instance" "name" {
  description   = "Terraform instance example"
  cloud_id      = data.morpheus_cloud.vsphere.id
  group_id      = data.morpheus_group.all.id
  type          = "centos"
  layout        = "centos"
  plan          = "1 CPU, 4GB Memory"
  environment   = "dev"
  resource_pool = "democluster"
  labels        = ["demo","terraform"]

  interfaces  {
    network   = "VM Network"
  }

  tags = {
    name  = "tfdemo"
  }

  evar {
    name   = "application"
    value  = "demo"
    export = true
    masked = true
  }
}
```