---
page_title: "morpheus_budget_policy Resource - terraform-provider-morpheus"
subcategory: ""
description: |-
  Provides a Morpheus budget policy resource
---

# morpheus_budget_policy

Provides a Morpheus budget policy resource

## Example Usage

Creating the policy with a global scope:

```terraform
resource "morpheus_budget_policy" "tf_example_budget_policy_global" {
  name         = "tf_example_budget_policy_global"
  description  = "terraform example global budget policy"
  enabled      = true
  max_price    = "4000"
  currency     = "USD"
  unit_of_time = "hour"
  scope        = "global"
}
```

Creating the policy with a cloud scope:

```terraform
resource "morpheus_budget_policy" "tf_example_budget_policy_cloud" {
  name         = "tf_example_budget_policy_cloud"
  description  = "terraform example cloud budget policy"
  enabled      = true
  max_price    = "4000"
  currency     = "USD"
  unit_of_time = "hour"
  scope        = "cloud"
  cloud_id     = 1
}
```

Creating the policy with a group scope:

```terraform
resource "morpheus_budget_policy" "tf_example_budget_policy_group" {
  name         = "tf_example_budget_policy_group"
  description  = "terraform example group budget policy"
  enabled      = true
  max_price    = "4000"
  currency     = "USD"
  unit_of_time = "hour"
  scope        = "group"
  group_id     = 1
}
```

Creating the policy with a role scope:

```terraform
resource "morpheus_budget_policy" "tf_example_budget_policy_role" {
  name            = "tf_example_budget_policy_role"
  description     = "terraform example role budget policy"
  enabled         = true
  max_price       = "4000"
  currency        = "USD"
  unit_of_time    = "hour"
  scope           = "role"
  role_id         = 1
  apply_each_user = true
}
```

Creating the policy with a user scope:

```terraform
resource "morpheus_budget_policy" "tf_example_budget_policy_user" {
  name         = "tf_example_budget_policy_user"
  description  = "terraform example user budget policy"
  enabled      = true
  max_price    = "4000"
  currency     = "USD"
  unit_of_time = "hour"
  scope        = "user"
  user_id      = 1
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `currency` (String) The budget currency
- `max_price` (String) The max budget price
- `name` (String) The name of the budget policy
- `scope` (String) The filter or scope that the policy is applied to (global, group, cloud, user, role)
- `unit_of_time` (String) The unit of time to measure the budget (hour or month)

### Optional

- `apply_to_each_user` (Boolean) Whether to assign the policy at the individual user level to all users assigned the associated role
- `cloud_id` (Number) The id of the cloud associated with the cloud scoped filter
- `description` (String) The description of the budget policy
- `enabled` (Boolean) Whether the policy is enabled
- `group_id` (Number) The id of the group associated with the group scoped filter
- `role_id` (Number) The id of the role associated with the role scoped filter
- `tenant_ids` (List of Number) A list of tenant IDs to assign the policy to
- `user_id` (Number) The id of the user associated with the user scoped filter

### Read-Only

- `id` (String) The ID of the budget policy

## Import

Import is supported using the following syntax:

```shell
terraform import morpheus_budget_policy.tf_example_budget_policy 1
```
