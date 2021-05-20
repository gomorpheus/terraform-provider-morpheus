---
subcategory: ""
page_title: "Create a Morpheus vSphere instance - Morpheus Provider"
description: |-
    An example of creating a Morpheus vSphere instance with optional fields defaulted.
---

# Create a Morpheus instance using the `morpheus_vsphere_instance` resource

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