---
page_title: "morpheus_security_package Data Source - terraform-provider-morpheus"
subcategory: ""
description: |-
  Provides a Morpheus security package data source.
---

# morpheus_security_package (Data Source)

Provides a Morpheus security package data source.

## Example Usage

```terraform
data "morpheus_security_package" "tf_example_security_package" {
  name = "tf_example_security_package"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `name` (String) The name of the Morpheus security package.

### Read-Only

- `id` (Number) The ID of this resource.