# CHANGELOG

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/) and this project adheres to [Semantic Versioning](http://semver.org/).

## Unreleased

- Remove project support

## 0.3.0 (2019-06-17)

- Support Terraform 0.12
- Add `dns` option to `cosmic_network` to allow configuring DNS resolver per tier
- Add `disk_controller` option to `cosmic_instance`
- Add `disk_controller` option to `cosmic_disk`

## 0.2.0 (2019-02-28)

- Add option to configure provider using `COSMIC_CONFIG` and `COSMIC_PROFILE` environment variables
- Add ability to use protocol numbers for `cosmic_network_acl_rule`'s `protocol` option (instead of only `icmp`, `tcp`, `udp`, `all`)
- Changing `cosmic_loadbalancer_rule`'s `member_ids`, `private_port`, `public_port` or `protocol` options no longer recreates the resource
- Changing `cosmic_network`'s `ip_exclusion_list` option no longer recreates the resource
- Changing `cosmic_vpc`'s `vpc_offering` option no longer recreates the resource
- Removed `cosmic_egress_firewall` and `cosmic_firewall` resources; no longer implemented by the Cosmic API

## 0.1.0 (2019-01-27)

First versioned release.

Recent additions include:

- Add `config` and `profile` options to configure cosmic provider using a cloudmonkey config
- Add `optimise_for` option for `cosmic_instance`
- Add `protocol` option for `cosmic_loadbalancer_rule`
- Add `terraform import ...` support to cosmic resources to import existing infrastructure
