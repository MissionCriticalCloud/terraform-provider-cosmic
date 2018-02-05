# terraform-provider-cosmic
Terraform provider for [Cosmic](https://github.com/MissionCriticalCloud/cosmic) framework

## Setup
Download binary specific for your release and place within your terraform repository ( or configure plugins location )

## Configuring
```
provider "cosmic" {
  api_url    = "${var.api_url}"
  api_key    = "${var.api_key}"
  secret_key = "${var.secret_key}"
}
```

## Resources
Implements following resources
- [x] cosmic_affinity_group
- [x] cosmic_disk
- [x] cosmic_egress_firewall
- [x] cosmic_firewall
- [x] cosmic_instance
- [x] cosmic_ipaddress
- [x] cosmic_loadbalancer_rule
- [x] cosmic_network
- [x] cosmic_network_acl
- [x] cosmic_network_acl_rule
- [x] cosmic_nic
- [x] cosmic_port_forward
- [x] cosmic_private_gateway
- [x] cosmic_secondary_ipaddress"
- [x] cosmic_security_group
- [x] cosmic_security_group_rule"
- [x] cosmic_ssh_keypair
- [x] cosmic_static_nat
- [x] cosmic_static_route
- [x] cosmic_template
- [x] cosmic_vpc
- [x] cosmic_vpn_connection
- [x] cosmic_vpn_customer_gateway
- [x] cosmic_vpn_gateway

## Changelog
wip