---
page_title: "morpheus_group Data Source - terraform-provider-morpheus"
subcategory: ""
description: |-
  Provides a Morpheus group data source.
---

# morpheus_group (Data Source)

Provides a Morpheus group data source.

## Example Usage

```terraform
data "morpheus_group" "morpheusgroup" {
  name = "Morpheus"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `name` (String) The name of the Morpheus group.

### Read-Only

- `code` (String) Optional code for use with policies
- `id` (Number) The ID of this resource.
- `location` (String) Optional location argument for your group