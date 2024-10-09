---
page_title: "morpheus_clouds Data Source - terraform-provider-morpheus"
subcategory: ""
description: |-
  Provides a Morpheus clouds data source.
---

# morpheus_clouds (Data Source)

Provides a Morpheus clouds data source.

## Example Usage

```terraform
data "morpheus_clouds" "tf_example_clouds" {
  sort_ascending = true
  filter {
    name   = "name"
    values = ["Test*"]
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filter` (Block Set) Custom filter block as described below. (see [below for nested schema](#nestedblock--filter))
- `sort_ascending` (Boolean) Whether to sort the IDs in ascending order

### Read-Only

- `id` (String) The ID of this resource.
- `ids` (List of Number)

<a id="nestedblock--filter"></a>
### Nested Schema for `filter`

Required:

- `name` (String) The name of the filter. Filter names are case-sensitive. Valid names are (name)
- `values` (Set of String) The filter values. Filter values are case-sensitive. Filters values support the use of Golang regex and can be tested at https://regex101.com/