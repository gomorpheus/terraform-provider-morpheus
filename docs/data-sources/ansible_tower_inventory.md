---
page_title: "morpheus_ansible_tower_inventory Data Source - terraform-provider-morpheus"
subcategory: ""
description: |-
  Provides a Morpheus ansible tower inventory data source.
---

# morpheus_ansible_tower_inventory (Data Source)

Provides a Morpheus ansible tower inventory data source.

## Example Usage

```terraform
data "morpheus_ansible_tower_inventory" "example_ansible_tower_inventory" {
  name = "Demo Inventory"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `name` (String) The name of the ansible tower inventory

### Read-Only

- `id` (Number) The ID of this resource.