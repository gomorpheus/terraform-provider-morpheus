---
page_title: "morpheus_tag_policy Resource - terraform-provider-morpheus"
subcategory: ""
description: |-
  Provides a Morpheus tag policy resource
---

# morpheus_tag_policy

Provides a Morpheus tag policy resource

## Example Usage

Creating the policy with a global scope:

```terraform
resource "morpheus_tag_policy" "tf_example_tag_policy_global" {
  name               = "tf_example_tag_policy_global"
  description        = "terraform example global tag policy"
  enabled            = true
  strict_enforcement = true
  tag_key            = "cost_center"
  option_list_id     = 23
  scope              = "global"
}
```

Creating the policy with a cloud scope:

```terraform
resource "morpheus_tag_policy" "tf_example_tag_policy_cloud" {
  name               = "tf_example_tag_policy_cloud"
  description        = "terraform example cloud tag policy"
  enabled            = true
  strict_enforcement = true
  tag_key            = "cost_center"
  option_list_id     = 23
  scope              = "cloud"
  cloud_id           = 1
}
```

Creating the policy with a group scope:

```terraform
resource "morpheus_tag_policy" "tf_example_tag_policy_group" {
  name               = "tf_example_tag_policy_group"
  description        = "terraform example group tag policy"
  enabled            = true
  strict_enforcement = true
  tag_key            = "true"
  tag_value          = "true"
  option_list_id     = 2
  scope              = "group"
  group_id           = 1
}
```

Creating the policy with a user scope:

```terraform
resource "morpheus_tag_policy" "tf_example_tag_policy_user" {
  name               = "tf_example_tag_policy_user"
  description        = "terraform example user tag policy"
  enabled            = true
  strict_enforcement = true
  tag_key            = "true"
  tag_value          = "true"
  scope              = "user"
  user_id            = 1
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the tag policy
- `scope` (String) The filter or scope that the policy is applied to (global, group, cloud, user)
- `tag_key` (String) The key of the tag to enforce

### Optional

- `cloud_id` (Number) The id of the cloud associated with the cloud scoped filter
- `description` (String) The description of the tag policy
- `enabled` (Boolean) Whether the policy is enabled
- `group_id` (Number) The id of the group associated with the group scoped filter
- `option_list_id` (Number) The id of the option list associated with the policy
- `strict_enforcement` (Boolean) Whether users will be able to provision new workloads if they violate the tag policy
- `tag_value` (String) The value of the tag to enforce
- `user_id` (Number) The id of the user associated with the user scoped filter

### Read-Only

- `id` (String) The ID of the tag policy

## Import

Import is supported using the following syntax:

```shell
terraform import morpheus_tag_policy.tf_example_tag_policy
```