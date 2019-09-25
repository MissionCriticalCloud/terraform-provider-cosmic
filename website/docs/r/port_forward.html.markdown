---
layout: "cosmic"
page_title: "Cosmic: cosmic_port_forward"
sidebar_current: "docs-cosmic-resource-port-forward"
description: |-
  Creates port forwards.
---

# cosmic_port_forward

Creates port forwards.

## Example Usage

```hcl
resource "cosmic_port_forward" "default" {
  ip_address_id = "30b21801-d4b3-4174-852b-0c0f30bdbbfb"

  forward {
    protocol           = "tcp"
    private_port       = 80
    public_port        = 8080
    virtual_machine_id = "f8141e2f-4e7e-4c63-9362-986c908b7ea7"
  }
}
```

## Argument Reference

The following arguments are supported:

* `ip_address_id` - (Required) The IP address ID for which to create the port
    forwards. Changing this forces a new resource to be created.

* `managed` - (Optional) USE WITH CAUTION! If enabled all the port forwards for
    this IP address will be managed by this resource. This means it will delete
    all port forwards that are not in your config! (defaults false)

* `forward` - (Required) Can be specified multiple times. Each forward block supports
    fields documented below.

The `forward` block supports:

* `protocol` - (Required) The name of the protocol to allow. Valid options are:
    `tcp` and `udp`.

* `private_port` - (Required) The private port to forward to.

* `public_port` - (Required) The public port to forward from.

* `virtual_machine_id` - (Required) The ID of the virtual machine to forward to.

* `vm_guest_ip` - (Optional) The virtual machine IP address for the port
    forwarding rule (useful when the virtual machine has secondairy NICs
    or IP addresses).

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the IP address for which the port forwards are created.
* `vm_guest_ip` - The IP address of the virtual machine that is used
    for the port forwarding rule.

## Import (EXPERIMENTAL)

port forwards can be imported; use `<PORT FORWARD ID>` as the import ID. For
example:

```shell
terraform import cosmic_port_forward.default e42a24d2-46cb-4b18-9d41-382582fad309
```

Multiple port forwards in the same resource can be imported by specifying a
comma separated string.
