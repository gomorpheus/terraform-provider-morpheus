---
page_title: "morpheus_cypher_access_policy Resource - terraform-provider-morpheus"
subcategory: ""
description: |-
  Provides a Morpheus cypher access policy resource
---

# morpheus_cypher_access_policy

Provides a Morpheus cypher access policy resource

## Example Usage

Creating the policy with a global scope:

```terraform
resource "morpheus_cypher_access_policy" "tf_example_cypher_access_policy_global" {
  name          = "tf_example_cypher_access_policy_global"
  description   = "terraform example global cypher access policy"
  enabled       = true
  key_path      = ".*"
  read_access   = true
  write_access  = true
  update_access = true
  list_access   = true
  delete_access = true
  scope         = "global"
}
```

Creating the policy with a role scope:

```terraform
resource "morpheus_cypher_access_policy" "tf_example_cypher_access_policy_role" {
  name               = "tf_example_cypher_access_policy_role"
  description        = "terraform example role cypher access policy"
  enabled            = true
  key_path           = ".*"
  read_access        = true
  write_access       = true
  update_access      = true
  list_access        = true
  delete_access      = true
  scope              = "role"
  role_id            = 1
  apply_to_each_user = true
}
```

Creating the policy with a user scope:

```terraform
resource "morpheus_cypher_access_policy" "tf_example_cypher_access_policy_user" {
  name          = "tf_example_cypher_access_policy_user"
  description   = "terraform example user cypher access policy"
  enabled       = true
  key_path      = ".*"
  read_access   = true
  write_access  = true
  update_access = true
  list_access   = true
  delete_access = true
  scope         = "user"
  user_id       = 1
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `key_path` (String) The key path associated with the cypher access policy
- `name` (String) The name of the cypher access policy
- `scope` (String) The filter or scope that the policy is applied to (global, user, role)

### Optional

- `apply_to_each_user` (Boolean) Whether to assign the policy at the individual user level to all users assigned the associated role
- `delete_access` (Boolean) Whether the policy grants delete access
- `description` (String) The description of the cypher access policy
- `enabled` (Boolean) Whether the policy is enabled
- `list_access` (Boolean) Whether the policy grants list access
- `read_access` (Boolean) Whether the policy grants read access
- `role_id` (Number) The id of the role associated with the role scoped filter
- `tenant_ids` (List of Number) A list of tenant IDs to assign the policy to
- `update_access` (Boolean) Whether the policy grants update access
- `user_id` (Number) The id of the user associated with the user scoped filter
- `write_access` (Boolean) Whether the policy grants write access

### Read-Only

- `id` (String) The ID of the cypher access policy

## Import

Import is supported using the following syntax:

```shell
terraform import morpheus_cypher_access_policy.tf_example_cypher_access_policy 1
```
