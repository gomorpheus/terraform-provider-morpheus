---
page_title: "morpheus_cloud_type Data Source - terraform-provider-morpheus"
subcategory: ""
description: |-
  Provides a Morpheus cloud type data source.
---

# morpheus_cloud_type (Data Source)

Provides a Morpheus cloud type data source.

## Example Usage

```terraform
data "morpheus_cloud_type" "canonical_maas_cloud" {
  name = "MaaS"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the Morpheus cloud type

### Read-Only

- `id` (Number) The ID of this resource.