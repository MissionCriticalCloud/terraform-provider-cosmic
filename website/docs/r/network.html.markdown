---
layout: "cosmic"
page_title: "Cosmic: cosmic_network"
sidebar_current: "docs-cosmic-resource-network"
description: |-
  Creates a network.
---

# cosmic_network

Creates a network.

## Example Usage

Basic usage:

```hcl
resource "cosmic_network" "default" {
  name             = "test-network"
  cidr             = "10.0.0.0/16"
  network_offering = "Default Network"
  zone             = "zone-1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the network.

* `display_text` - (Optional) The display text of the network.

* `cidr` - (Required) The CIDR block for the network. Changing this forces a new
    resource to be created.

* `gateway` - (Optional) Gateway that will be provided to the instances in this
    network. Defaults to the first usable IP in the range.

* `startip` - (Optional) Start of the IP block that will be available on the
    network. Defaults to the second available IP in the range.

* `endip` - (Optional) End of the IP block that will be available on the
    network. Defaults to the last available IP in the range.

* `network_domain` - (Optional) DNS domain for the network.

* `network_offering` - (Required) The name or ID of the network offering to use
    for this network.

* `vlan` - (Optional) The VLAN number (1-4095) the network will use. This might be
    required by the Network Offering if specifyVlan=true is set. Only the ROOT
    admin can set this value.

* `vpc_id` - (Optional) The VPC ID in which to create this network. Changing
    this forces a new resource to be created.

* `acl_id` - (Optional) The ACL ID that should be attached to the network or
    `none` if you do not want to attach an ACL. You can dynamically attach and
    swap ACL's, but if you want to detach an attached ACL and revert to using
    `none`, this will force a new resource to be created. (defaults `none`)

* `ip_exclusion_list` - (Optional) list of ip addresses and/or ranges of
    addresses to be excluded from the network for assignment to instances
    in this network (eg. 10.0.0.2,10.0.0.4-10.0.0.7).

* `zone` - (Required) The name or ID of the zone where this network will be
    available. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the network.
* `display_text` - The display text of the network.
* `network_domain` - DNS domain for the network.
* `source_nat_ip_id` - The ID of the associated source NAT IP.

## Import (EXPERIMENTAL)

Networks can be imported; use `<NETWORK ID>` as the import ID. For
example:

```shell
terraform import cosmic_network.default 36619b20-5584-43bf-9a84-e242bacd5582
```
