---
layout: "cosmic"
page_title: "Cosmic: cosmic_ipaddress"
sidebar_current: "docs-cosmic-resource-ipaddress"
description: |-
  Acquires and associates a public IP.
---

# cosmic_ipaddress

Acquires and associates a public IP.

## Example Usage

```hcl
resource "cosmic_ipaddress" "default" {
  network_id = "6eb22f91-7454-4107-89f4-36afcdf33021"
}
```

## Argument Reference

The following arguments are supported:

* `is_portable` - (Optional) This determines if the IP address should be transferable
    across zones (defaults false)

* `network_id` - (Optional) The ID of the network for which an IP address should
    be acquired and associated. Changing this forces a new resource to be created.

* `vpc_id` - (Optional) The ID of the VPC for which an IP address should be
   acquired and associated. Changing this forces a new resource to be created.

* `zone` - (Optional) The name or ID of the zone for which an IP address should be
   acquired and associated. Changing this forces a new resource to be created.

*NOTE: `network_id` and/or `zone` should have a value when `is_portable` is `false`!*
*NOTE: Either `network_id` or `vpc_id` should have a value when `is_portable` is `true`!*

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the acquired and associated IP address.
* `ip_address` - The IP address that was acquired and associated.

## Import (EXPERIMENTAL)

IP addresses can be imported; use `<IP ADDRESS ID>` as the import ID. For
example:

```shell
terraform import cosmic_ipaddress.default e42a24d2-46cb-4b18-9d41-382582fad309
```
