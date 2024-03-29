---
page_title: "morpheus_resource_pool Data Source - terraform-provider-morpheus"
subcategory: ""
description: |-
  Provides a Morpheus resource pool data source.
---

# morpheus_resource_pool (Data Source)

Provides a Morpheus resource pool data source.

## Example Usage

```terraform
data "morpheus_resource_pool" "morpheus_pool" {
  name     = "morpheuspool"
  cloud_id = data.morpheus_cloud.vspherecloud.id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cloud_id` (Number) The id of the Morpheus cloud to search for the resource pool.

### Optional

- `id` (Number) The id of the resource pool
- `name` (String) The name of the Morpheus resource pool.

### Read-Only

- `active` (Boolean) Whether the resource pool is enabled or not
- `description` (String) The description of the resource pool
- `type` (String) Optional code for use with policies