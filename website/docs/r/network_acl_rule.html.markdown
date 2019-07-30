---
layout: "cosmic"
page_title: "Cosmic: cosmic_network_acl_rule"
sidebar_current: "docs-cosmic-resource-network-acl-rule"
description: |-
  Creates network ACL rules for a given network ACL.
---

# cosmic_network_acl_rule

Creates network ACL rules for a given network ACL.

## Example Usage

```hcl
resource "cosmic_network_acl_rule" "default" {
  acl_id = "f3843ce0-334c-4586-bbd3-0c2e2bc946c6"

  rule {
    action       = "allow"
    cidr_list    = ["10.0.0.0/8"]
    protocol     = "tcp"
    ports        = ["80", "1000-2000"]
    traffic_type = "ingress"
  }
}
```

## Argument Reference

The following arguments are supported:

* `acl_id` - (Required) The network ACL ID for which to create the rules.
    Changing this forces a new resource to be created.

* `managed` - (Optional) USE WITH CAUTION! If enabled all the firewall rules for
    this network ACL will be managed by this resource. This means it will delete
    all firewall rules that are not in your config! (defaults false)

* `rule` - (Optional) Can be specified multiple times. Each rule block supports
    fields documented below. If `managed = false` at least one rule is required!

* `parallelism` (Optional) Specifies how much rules will be created or deleted
    concurrently. (defaults 2)

The `rule` block supports:

* `action` - (Optional) The action for the rule. Valid options are: `allow` and
    `deny` (defaults allow).

* `cidr_list` - (Required) A CIDR list to allow access to the given ports.

* `protocol` - (Required) The name of the protocol to allow. Valid options are:
    `tcp`, `udp`, `icmp`, `all` or a valid protocol number.

* `icmp_type` - (Optional) The ICMP type to allow, or `-1` to allow `any`. This
    can only be specified if the protocol is ICMP. (defaults 0)

* `icmp_code` - (Optional) The ICMP code to allow, or `-1` to allow `any`. This
    can only be specified if the protocol is ICMP. (defaults 0)

* `ports` - (Optional) List of ports and/or port ranges to allow. This can only
    be specified if the protocol is TCP, UDP, ALL or a valid protocol number.

* `traffic_type` - (Optional) The traffic type for the rule. Valid options are:
    `ingress` or `egress` (defaults ingress).

## Attributes Reference

The following attributes are exported:

* `id` - The ACL ID for which the rules are created.

## Import (EXPERIMENTAL)

Network ACL rules can be imported; use `<NETWORK ACL RULE ID>` as the import ID. For
example:

```shell
terraform import cosmic_network_acl_rule.default e42a24d2-46cb-4b18-9d41-382582fad309
```
