---
page_title: "morpheus_max_storage_policy Resource - terraform-provider-morpheus"
subcategory: ""
description: |-
  Provides a Morpheus max storage policy resource
---

# morpheus_max_storage_policy

Provides a Morpheus max storage policy resource

## Example Usage

Creating the policy with a global scope:

```terraform
resource "morpheus_max_storage_policy" "tf_example_max_storage_policy_global" {
  name        = "tf_example_max_storage_policy_global"
  description = "terraform example global max storage policy"
  enabled     = true
  max_storage = 100
  scope       = "global"
}
```

Creating the policy with a cloud scope:

```terraform
resource "morpheus_max_storage_policy" "tf_example_max_storage_policy_cloud" {
  name        = "tf_example_max_storage_policy_cloud"
  description = "terraform example cloud max storage policy"
  enabled     = true
  max_storage = 100
  scope       = "cloud"
  cloud_id    = 1
}
```

Creating the policy with a group scope:

```terraform
resource "morpheus_max_storage_policy" "tf_example_max_storage_policy_group" {
  name        = "tf_example_max_storage_policy_group"
  description = "terraform example group max storage policy"
  enabled     = true
  max_storage = 100
  scope       = "group"
  group_id    = 1
}
```

Creating the policy with a role scope:

```terraform
resource "morpheus_max_storage_policy" "tf_example_max_storage_policy_role" {
  name            = "tf_example_max_storage_policy_role"
  description     = "terraform example role max storage policy"
  enabled         = true
  max_storage     = 100
  scope           = "role"
  role_id         = 1
  apply_each_user = true
}
```

Creating the policy with a user scope:

```terraform
resource "morpheus_max_storage_policy" "tf_example_max_storage_policy_user" {
  name        = "tf_example_max_storage_policy_user"
  description = "terraform example user max storage policy"
  enabled     = true
  max_storage = 100
  scope       = "user"
  user_id     = 1
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `max_storage` (Number) The maximum storage defined by the policy in GB
- `name` (String) The name of the max storage policy
- `scope` (String) The filter or scope that the policy is applied to (global, group, cloud, user, role)

### Optional

- `apply_to_each_user` (Boolean) Whether to assign the policy at the individual user level to all users assigned the associated role
- `cloud_id` (Number) The id of the cloud associated with the cloud scoped filter
- `description` (String) The description of the max storage policy
- `enabled` (Boolean) Whether the policy is enabled
- `group_id` (Number) The id of the group associated with the group scoped filter
- `role_id` (Number) The id of the role associated with the role scoped filter
- `tenant_ids` (List of Number) A list of tenant IDs to assign the policy to
- `user_id` (Number) The id of the user associated with the user scoped filter

### Read-Only

- `id` (String) The ID of the workflow policy

## Import

Import is supported using the following syntax:

```shell
terraform import morpheus_max_storage_policy.tf_example_max_storage_policy 1
```
