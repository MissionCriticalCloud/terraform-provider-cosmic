---
layout: "cosmic"
page_title: "Cosmic: cosmic_affinity_group"
sidebar_current: "docs-cosmic-resource-affinity-group"
description: |-
  Creates an affinity group.
---

# cosmic_affinity_group

Creates an affinity group.

## Example Usage

```hcl
resource "cosmic_affinity_group" "default" {
  name = "test-affinity-group"
  type = "host anti-affinity"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the affinity group. Changing this
    forces a new resource to be created.

* `description` - (Optional) The description of the affinity group.

* `type` - (Required) The affinity group type. Changing this
    forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The id of the affinity group.
* `description` - The description of the affinity group.

## Import (EXPERIMENTAL)

Affinity groups can be imported; use `<AFFINITY GROUP ID>` as the import ID. For
example:

```shell
terraform import cosmic_affinity_group.default 6226ea4d-9cbe-4cc9-b30c-b9532146da5b
```
