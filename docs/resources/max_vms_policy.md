---
page_title: "morpheus_max_vms_policy Resource - terraform-provider-morpheus"
subcategory: ""
description: |-
  Provides a Morpheus max vms policy resource
---

# morpheus_max_vms_policy

Provides a Morpheus max vms policy resource

## Example Usage

```terraform
resource "morpheus_max_vms_policy" "tf_example_max_vms_policy_global" {
  name        = "tf_example_max_vms_policy_global"
  description = "Terraform example Morpheus max vms policy"
  enabled     = true
  max_vms     = 35
  scope       = "global"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **max_vms** (Number) The maximum vms defined by the policy
- **name** (String) The name of the workflow policy
- **scope** (String) The filter or scope that the policy is applied to (global, group, cloud, user, role)

### Optional

- **apply_to_each_user** (Boolean) Whether to assign the policy at the individual user level to all users assigned the associated role
- **cloud_id** (Number) The id of the cloud associated with the cloud scoped filter
- **description** (String) The description of the workflow policy
- **enabled** (Boolean) Whether the policy is enabled
- **group_id** (Number) The id of the group associated with the gropu scoped filter
- **role_id** (Number) The id of the role associated with the role scoped filter
- **user_id** (Number) The id of the user associated with the user scoped filter

### Read-Only

- **id** (String) The ID of the workflow policy

## Import

Import is supported using the following syntax:

```shell
terraform import morpheus_max_vms_policy.tf_example_max_vms_policy 1
```