---
page_title: "morpheus_provision_approval_policy Resource - terraform-provider-morpheus"
subcategory: ""
description: |-
  Provides a Morpheus provision approval policy resource
---

# morpheus_provision_approval_policy

Provides a Morpheus provision approval policy resource

## Example Usage

Creating the policy with a global scope:

```terraform
resource "morpheus_provision_approval_policy" "tf_example_provision_approval_policy_global" {
  name                   = "tf_example_provision_approval_policy_global"
  description            = "terraform example global provision approval policy"
  enabled                = true
  use_internal_approvals = true
  scope                  = "global"
}
```

Creating the policy with a cloud scope:

```terraform
resource "morpheus_provision_approval_policy" "tf_example_provision_approval_policy_global" {
  name           = "tf_example_provision_approval_policy_global"
  description    = "terraform example global provision approval policy"
  enabled        = true
  integration_id = 1
  workflow_id    = 10
  scope          = "cloud"
  cloud_id       = 1
}
```

Creating the policy with a group scope:

```terraform
resource "morpheus_provision_approval_policy" "tf_example_provision_approval_policy_global" {
  name                   = "tf_example_provision_approval_policy_global"
  description            = "terraform example global provision approval policy"
  enabled                = true
  use_internal_approvals = true
  scope                  = "group"
  group_id               = 1
}
```

Creating the policy with a role scope:

```terraform
resource "morpheus_provision_approval_policy" "tf_example_provision_approval_policy_global" {
  name            = "tf_example_provision_approval_policy_global"
  description     = "terraform example global provision approval policy"
  enabled         = true
  integration_id  = 1
  workflow_id     = 10
  scope           = "role"
  role_id         = 1
  apply_each_user = true
}
```

Creating the policy with a user scope:

```terraform
resource "morpheus_provision_approval_policy" "tf_example_provision_approval_policy_global" {
  name                   = "tf_example_provision_approval_policy_global"
  description            = "terraform example global provision approval policy"
  enabled                = true
  use_internal_approvals = true
  scope                  = "user"
  user_id                = 1
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the provision approval policy
- `scope` (String) The filter or scope that the policy is applied to (global, group, cloud, user, role)

### Optional

- `apply_to_each_user` (Boolean) Whether to assign the policy at the individual user level to all users assigned the associated role
- `cloud_id` (Number) The id of the cloud associated with the cloud scoped filter
- `description` (String) The description of the provision approval policy
- `enabled` (Boolean) Whether the policy is enabled
- `group_id` (Number) The id of the group associated with the group scoped filter
- `integration_id` (Number) The ID of the approval integration used for approvals
- `role_id` (Number) The id of the role associated with the role scoped filter
- `tenant_ids` (List of Number) A list of tenant IDs to assign the policy to
- `use_internal_approvals` (Boolean) Whether the internal Morpheus approval engine is used for approvals
- `user_id` (Number) The id of the user associated with the user scoped filter
- `workflow_id` (Number) The ID of the approval workflow used for approvals

### Read-Only

- `id` (String) The ID of the provision approval policy

## Import

Import is supported using the following syntax:

```shell
terraform import morpheus_provision_approval_policy.tf_example_provision_approval_policy 1
```
