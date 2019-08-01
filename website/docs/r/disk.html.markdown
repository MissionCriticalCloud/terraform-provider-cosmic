---
layout: "cosmic"
page_title: "Cosmic: cosmic_disk"
sidebar_current: "docs-cosmic-resource-disk"
description: |-
  Creates a disk volume from a disk offering. This disk volume will be attached to a virtual machine if the optional parameters are configured.
---

# cosmic_disk

Creates a disk volume from a disk offering. This disk volume will be attached to
a virtual machine if the optional parameters are configured.

## Example Usage

```hcl
resource "cosmic_disk" "default" {
  name               = "test-disk"
  attach             = "true"
  disk_offering      = "custom"
  size               = 50
  disk_controller    = "VIRTIO"
  virtual_machine_id = "server-1"
  zone               = "zone-1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the disk volume. Changing this forces a new
    resource to be created.

* `attach` - (Optional) Determines whether or not to attach the disk volume to a
    virtual machine (defaults false).

* `device_id` - (Optional) The device ID to map the disk volume to within the guest OS.

* `disk_offering` - (Required) The name or ID of the disk offering to use for
    this disk volume.

* `size` - (Optional) The size of the disk volume in gigabytes.

* `shrink_ok` - (Optional) Verifies if the disk volume is allowed to shrink when
    resizing (defaults false).

* `disk_controller` - (Optional) Set the disk controller type for this disk. Allowed
     options are IDE, SCSI and VIRTIO. Changing this forces a new resource to be created.

* `virtual_machine_id` - (Optional) The ID of the virtual machine to which you want
    to attach the disk volume.

* `zone` - (Required) The name or ID of the zone where this disk volume will be available.
    Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the disk volume.
* `device_id` - The device ID the disk volume is mapped to within the guest OS.

## Import (EXPERIMENTAL)

Disks can be imported; use `<DISK ID>` as the import ID. For
example:

```shell
terraform import cosmic_disk.default 6f3ee798-d417-4e7a-92bc-95ad41cf1244
```
