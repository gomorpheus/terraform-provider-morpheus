---
page_title: "morpheus_integration Data Source - terraform-provider-morpheus"
subcategory: ""
description: |-
  Provides a Morpheus integration data source.
---

# morpheus_integration (Data Source)

Provides a Morpheus integration data source.

## Example Usage

```terraform
data "morpheus_integration" "tf_example_integration" {
  name = "ansible dev"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `name` (String) The name of the integration

### Read-Only

- `id` (Number) The ID of this resource.