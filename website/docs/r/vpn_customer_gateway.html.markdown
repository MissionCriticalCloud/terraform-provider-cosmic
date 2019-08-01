---
layout: "cosmic"
page_title: "Cosmic: cosmic_vpn_customer_gateway"
sidebar_current: "docs-cosmic-resource-vpn-customer-gateway"
description: |-
  Creates a site to site VPN local customer gateway.
---

# cosmic_vpn_customer_gateway

Creates a site to site VPN local customer gateway.

## Example Usage

Basic usage:

```hcl
resource "cosmic_vpn_customer_gateway" "default" {
  name       = "test-vpc"
  cidr       = "10.0.0.0/8"
  esp_policy = "aes256-sha1;modp1024"
  gateway    = "192.168.0.1"
  ike_policy = "aes256-sha1;modp1024"
  ipsec_psk  = "terraform"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the VPN Customer Gateway.

* `cidr_list` - (Required) a list of CIDRs that needs to be routed through this gateway.

* `esp_policy` - (Required) The ESP policy to use for this VPN Customer Gateway.

* `gateway` - (Required) The public IP address of the related VPN Gateway.

* `ike_policy` - (Required) The IKE policy to use for this VPN Customer Gateway.

* `ipsec_psk` - (Required) The IPSEC pre-shared key used for this gateway.

* `dpd` - (Optional) If DPD is enabled for the related VPN connection (defaults false)

* `esp_lifetime` - (Optional) The ESP lifetime of phase 2 VPN connection to this
    VPN Customer Gateway in seconds (defaults 86400)

* `ike_lifetime` - (Optional) The IKE lifetime of phase 2 VPN connection to this
    VPN Customer Gateway in seconds (defaults 86400)

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the VPN Customer Gateway.
* `dpd` - Enable or disable DPD is enabled for the related VPN connection.
* `esp_lifetime` - The ESP lifetime of phase 2 VPN connection to this VPN Customer Gateway.
* `ike_lifetime` - The IKE lifetime of phase 2 VPN connection to this VPN Customer Gateway.

## Import (EXPERIMENTAL)

VPN customer gateways can be imported; use `<VPN CUSTOMER GATEWAY ID>` as the import ID. For
example:

```shell
terraform import cosmic_vpn_customer_gateway.default 741a7fca-1d05-4bb6-9290-1008300f0e5a
```
