---
page_title: "morpheus_tenant_role Data Source - terraform-provider-morpheus"
subcategory: ""
description: |-
  Provides a Morpheus tenant role data source.
---

# morpheus_tenant_role (Data Source)

Provides a Morpheus tenant role data source.

## Example Usage

```terraform
data "morpheus_tenant_role" "example_tenant_role" {
  name = "Tenant Admin"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `name` (String) The name of the Morpheus tenant role.

### Read-Only

- `id` (Number) The ID of this resource.