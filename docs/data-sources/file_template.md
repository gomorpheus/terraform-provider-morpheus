---
page_title: "morpheus_file_template Data Source - terraform-provider-morpheus"
subcategory: ""
description: |-
  Provides a Morpheus file template data source.
---

# morpheus_file_template (Data Source)

Provides a Morpheus file template data source.

## Example Usage

```terraform
data "morpheus_file_template" "example_file_template" {
  name = "Terraform Example File Template"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `name` (String) The name of the Morpheus file template.

### Read-Only

- `id` (Number) The ID of this resource.