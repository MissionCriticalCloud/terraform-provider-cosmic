---
layout: "cosmic"
page_title: "Cosmic: cosmic_secondary_ipaddress"
sidebar_current: "docs-cosmic-resource-secondary-ipaddress"
description: |-
  Assigns a secondary IP to a NIC.
---

# cosmic_secondary_ipaddress

Assigns a secondary IP to a NIC.

## Example Usage

```hcl
resource "cosmic_secondary_ipaddress" "default" {
  virtual_machine_id = "server-1"
}
```

## Argument Reference

The following arguments are supported:

* `ip_address` - (Optional) The IP address to bind the to NIC. If not supplied
    an IP address will be selected randomly. Changing this forces a new resource
    to be	created.

* `nic_id` - (Optional) The NIC ID to which you want to attach the secondary IP
    address. Changing this forces a new resource to be created (defaults to the
    ID of the primary NIC)

* `virtual_machine_id` - (Required) The ID of the virtual machine to which you
    want to attach the secondary IP address. Changing this forces a new resource
    to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The secondary IP address ID.
* `ip_address` - The IP address that was acquired and associated.

## Import (EXPERIMENTAL)

Secondary IP addresses can be imported; use `<INSTANCE ID>/<SECONDARY IP ADDRESS>` as the import ID. For
example:

```shell
terraform import cosmic_secondary_ipaddress.default e42a24d2-46cb-4b18-9d41-382582fad309/192.168.10.10
```
