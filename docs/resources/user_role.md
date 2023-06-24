---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "morpheus_user_role Resource - terraform-provider-morpheus"
subcategory: ""
description: |-
  Provides a Morpheus user role resource.
---

# morpheus_user_role (Resource)

Provides a Morpheus user role resource.

## Example Usage

```terraform
resource "morpheus_user_role" "tfexample_resource_user_role" {
  name               = "tf-example-user-role"
  description        = "Terraform provider example user role"
  multitenant_role   = false
  multitenant_locked = false
  permission_set     = data.morpheus_permission_set.base_permission_set.json
}

data "morpheus_group" "demo" {
  name = "Demo"
}

data "morpheus_instance_type" "demo" {
  name = "Demo"
}

data "morpheus_blueprint" "demo" {
  name = "Demo"
}

data "morpheus_catalog_item_type" "demo" {
  name = "Demo"
}

data "morpheus_vdi_pool" "demo" {
  name = "Demo"
}

data "morpheus_task" "demo" {
  name = "Demo"
}

data "morpheus_workflow" "demo" {
  name = "Demo"
}

data "morpheus_permission_set" "base_permission_set" {
  override_permission_sets = [
    data.morpheus_permission_set.override_set.json,
  ]
  default_group_permission             = "full"
  default_instance_type_permission     = "none"
  default_blueprint_permission         = "none"
  default_report_type_permission       = "full"
  default_persona                      = "vdi"
  default_catalog_item_type_permission = "full"
  default_vdi_pool_permission          = "full"
  default_workflow_permission          = "full"
  default_task_permission              = "full"

  feature_permission {
    code   = "provisioning-admin"
    access = "full"
  }

  group_permission {
    id     = data.morpheus_group.demo.id
    access = "full"
  }

  instance_type_permission {
    id     = data.morpheus_instance_type.demo.id
    access = "full"
  }

  blueprint_permission {
    id     = data.morpheus_blueprint.demo.id
    access = "full"
  }

  report_type_permission {
    code   = "guidance"
    access = "full"
  }

  persona_permission {
    code   = "standard"
    access = "full"
  }

  persona_permission {
    code   = "serviceCatalog"
    access = "none"
  }

  catalog_item_type_permission {
    id     = data.morpheus_catalog_item_type.demo.id
    access = "full"
  }

  vdi_pool_permission {
    id     = data.morpheus_vdi_pool.demo.id
    access = "full"
  }

  workflow_permission {
    id     = data.morpheus_workflow.demo.id
    access = "full"
  }

  task_permission {
    id     = data.morpheus_task.demo.id
    access = "none"
  }
}

data "morpheus_permission_set" "override_set" {
  default_task_permission = "none"
  default_persona         = "standard"
  workflow_permission {
    id     = 2
    access = "full"
  }
  workflow_permission {
    id     = 11
    access = "full"
  }
  group_permission {
    id     = 1
    access = "read"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the user role

### Optional

- `description` (String) The description of the user role
- `multitenant_locked` (Boolean) Whether subtenants are allowed to branch off or modify this role.
- `multitenant_role` (Boolean) Whether the user role is automatically copied into all existing subtenants as well as placed into a subtenant when created
- `permission_set` (String)

### Read-Only

- `id` (String) The ID of the user role

## Import

Import is supported using the following syntax:

```shell
terraform import morpheus_user_role.tf_example_user_role 1
```