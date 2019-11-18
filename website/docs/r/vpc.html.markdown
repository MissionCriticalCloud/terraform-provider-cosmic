---
layout: "cosmic"
page_title: "Cosmic: cosmic_vpc"
sidebar_current: "docs-cosmic-resource-vpc"
description: |-
  Creates a VPC.
---

# cosmic_vpc

Creates a VPC.

## Example Usage

Basic usage:

```hcl
resource "cosmic_vpc" "default" {
  name         = "test-vpc"
  cidr         = "10.0.0.0/16"
  vpc_offering = "Default VPC Offering"
  zone         = "zone-1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the VPC.

* `display_text` - (Optional) The display text of the VPC.

* `cidr` - (Required) The CIDR block for the VPC. Changing this forces a new
    resource to be created.

* `vpc_offering` - (Required) The name or ID of the VPC offering to use for this VPC.
    Changing this forces a new resource to be created.

* `network_domain` - (Optional) The default DNS domain for networks created in
    this VPC. Changing this forces a new resource to be created.

* `source_nat_list` - Source Nat CIDR list for used to allow other CIDRs to be
    source NATted by the VPC over the public interface.

* `syslog_server_list` - Comma separated list of IP addresses to configure as syslog
    servers on the VPC to forward IP tables logging.

* `zone` - (Required) The name or ID of the zone where this disk volume will be
    available. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the VPC.
* `display_text` - The display text of the VPC.
* `source_nat_ip` - The source NAT IP assigned to the VPC.
* `source_nat_ip_id` - The ID of the source NAT IP assigned to the VPC.

## Import (EXPERIMENTAL)

VPCs can be imported; use `<VPC ID>` as the import ID. For
example:

```shell
terraform import cosmic_vpc.default 84b23264-917a-4712-b8bf-cd7604db43b0
```
