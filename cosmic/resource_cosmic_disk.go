package cosmic

import (
	"fmt"
	"strings"

	"github.com/MissionCriticalCloud/go-cosmic/v6/cosmic"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCosmicDisk() *schema.Resource {
	return &schema.Resource{
		Create: resourceCosmicDiskCreate,
		Read:   resourceCosmicDiskRead,
		Update: resourceCosmicDiskUpdate,
		Delete: resourceCosmicDiskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"disk_offering": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},

			"attach": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"device_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"size": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"shrink_ok": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"disk_controller": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					switch v {
					case "IDE", "SCSI", "VIRTIO":
					default:
						errs = append(errs, fmt.Errorf("%q must be either 'IDE', 'SCSI' or 'VIRTIO', got: %q", key, v))
					}

					return
				},
			},

			"virtual_machine_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"zone": &schema.Schema{
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				ForceNew:   true,
				Deprecated: deprecatedZoneMsg(),
			},
		},
	}
}

func resourceCosmicDiskCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)
	d.Partial(true)

	name := d.Get("name").(string)

	// Create a new parameter struct
	p := client.Volume.NewCreateVolumeParams()
	p.SetName(name)

	// Retrieve the disk_offering ID
	diskofferingid, e := retrieveID(client, "disk_offering", d.Get("disk_offering").(string))
	if e != nil {
		return e.Error()
	}
	// Set the disk_offering ID
	p.SetDiskofferingid(diskofferingid)

	if d.Get("size").(int) != 0 {
		// Set the volume size
		p.SetSize(int64(d.Get("size").(int)))
	}

	// Set the disk_controller type if configured
	if diskcontroller, ok := d.GetOk("disk_controller"); ok {
		p.SetDiskcontroller(diskcontroller.(string))
	}

	// Retrieve the zone ID
	zoneid, e := retrieveID(client, "zone", client.ZoneName)
	if e != nil {
		return e.Error()
	}
	// Set the zone ID
	p.SetZoneid(zoneid)

	// Create the new volume
	r, err := client.Volume.CreateVolume(p)
	if err != nil {
		return fmt.Errorf("Error creating the new disk %s: %s", name, err)
	}

	// Set the volume ID and partials
	d.SetId(r.Id)
	d.SetPartial("name")
	d.SetPartial("device_id")
	d.SetPartial("disk_offering")
	d.SetPartial("size")
	d.SetPartial("disk_controller")
	d.SetPartial("virtual_machine_id")
	d.SetPartial("zone")

	if d.Get("attach").(bool) {
		err := resourceCosmicDiskAttach(d, meta)
		if err != nil {
			return fmt.Errorf("Error attaching the new disk %s to virtual machine: %s", name, err)
		}

		// Set the additional partial
		d.SetPartial("attach")
	}

	d.Partial(false)
	return resourceCosmicDiskRead(d, meta)
}

func resourceCosmicDiskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	// Get the volume details
	v, count, err := client.Volume.GetVolumeByID(d.Id())
	if err != nil {
		if count == 0 {
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("name", v.Name)
	d.Set("attach", v.Attached != "")           // If attached this contains a timestamp when attached
	d.Set("size", int(v.Size/(1024*1024*1024))) // Needed to get GB's again
	d.Set("disk_controller", v.Diskcontroller)

	setValueOrID(d, "disk_offering", v.Diskofferingname, v.Diskofferingid)
	setValueOrID(d, "zone", v.Zonename, v.Zoneid)

	if v.Attached != "" {
		d.Set("device_id", int(v.Deviceid))
		d.Set("virtual_machine_id", v.Virtualmachineid)
	}

	return nil
}

func resourceCosmicDiskUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)
	d.Partial(true)

	name := d.Get("name").(string)

	if d.HasChange("disk_offering") || d.HasChange("size") {
		// Create a new parameter struct
		p := client.Volume.NewResizeVolumeParams(d.Id())

		// Retrieve the disk_offering ID
		diskofferingid, e := retrieveID(client, "disk_offering", d.Get("disk_offering").(string))
		if e != nil {
			return e.Error()
		}

		// Set the disk_offering ID
		p.SetDiskofferingid(diskofferingid)

		if d.Get("size").(int) != 0 {
			// Set the size
			p.SetSize(int64(d.Get("size").(int)))
		}

		// Set the shrink bit
		p.SetShrinkok(d.Get("shrink_ok").(bool))

		// Change the disk_offering
		r, err := client.Volume.ResizeVolume(p)
		if err != nil {
			return fmt.Errorf("Error changing disk offering/size for disk %s: %s", name, err)
		}

		// Update the volume ID and set partials
		d.SetId(r.Id)
		d.SetPartial("disk_offering")
		d.SetPartial("size")
	}

	// If the device ID changed, just detach here so we can re-attach the
	// volume at the end of this function
	if d.HasChange("device_id") || d.HasChange("virtual_machine") {
		// Detach the volume
		if err := resourceCosmicDiskDetach(d, meta); err != nil {
			return fmt.Errorf("Error detaching disk %s from virtual machine: %s", name, err)
		}
	}

	if d.Get("attach").(bool) {
		// Attach the volume
		err := resourceCosmicDiskAttach(d, meta)
		if err != nil {
			return fmt.Errorf("Error attaching disk %s to virtual machine: %s", name, err)
		}

		// Set the additional partials
		d.SetPartial("attach")
		d.SetPartial("device_id")
		d.SetPartial("virtual_machine_id")
	} else {
		// Detach the volume
		if err := resourceCosmicDiskDetach(d, meta); err != nil {
			return fmt.Errorf("Error detaching disk %s from virtual machine: %s", name, err)
		}
	}

	d.Partial(false)
	return resourceCosmicDiskRead(d, meta)
}

func resourceCosmicDiskDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	// Detach the volume
	if err := resourceCosmicDiskDetach(d, meta); err != nil {
		return err
	}

	// Create a new parameter struct
	p := client.Volume.NewDeleteVolumeParams(d.Id())

	// Delete the voluem
	if _, err := client.Volume.DeleteVolume(p); err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", d.Id())) {
			return nil
		}

		return err
	}

	return nil
}

func resourceCosmicDiskAttach(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	if virtualmachineid, ok := d.GetOk("virtual_machine_id"); ok {
		// First check if the disk isn't already attached
		if attached, err := isAttached(d, meta); err != nil || attached {
			return err
		}

		// Create a new parameter struct
		p := client.Volume.NewAttachVolumeParams(d.Id(), virtualmachineid.(string))

		if deviceid, ok := d.GetOk("device_id"); ok {
			p.SetDeviceid(int64(deviceid.(int)))
		}

		// Attach the new volume
		r, err := Retry(10, retryableAttachVolumeFunc(client, p))
		if err != nil {
			return fmt.Errorf("Error attaching volume to VM: %s", err)
		}

		d.SetId(r.(*cosmic.AttachVolumeResponse).Id)
	}

	return nil
}

func resourceCosmicDiskDetach(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	// Check if the volume is actually attached, before detaching
	if attached, err := isAttached(d, meta); err != nil || !attached {
		return err
	}

	// Create a new parameter struct
	p := client.Volume.NewDetachVolumeParams()

	// Set the volume ID
	p.SetId(d.Id())

	// Detach the currently attached volume
	_, err := client.Volume.DetachVolume(p)
	if err != nil {
		if virtualmachineid, ok := d.GetOk("virtual_machine_id"); ok {
			// Create a new parameter struct
			pd := client.VirtualMachine.NewStopVirtualMachineParams(virtualmachineid.(string))

			// Stop the virtual machine in order to be able to detach the disk
			if _, err := client.VirtualMachine.StopVirtualMachine(pd); err != nil {
				return err
			}

			// Try again to detach the currently attached volume
			if _, err := client.Volume.DetachVolume(p); err != nil {
				return err
			}

			// Create a new parameter struct
			pu := client.VirtualMachine.NewStartVirtualMachineParams(virtualmachineid.(string))

			// Start the virtual machine again
			if _, err := client.VirtualMachine.StartVirtualMachine(pu); err != nil {
				return err
			}
		}
	}

	return err
}

func isAttached(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*CosmicClient)

	// Get the volume details
	v, _, err := client.Volume.GetVolumeByID(d.Id())
	if err != nil {
		return false, err
	}

	return v.Attached != "", nil
}

func retryableAttachVolumeFunc(client *CosmicClient, p *cosmic.AttachVolumeParams) func() (interface{}, error) {
	return func() (interface{}, error) {

		r, err := client.Volume.AttachVolume(p)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
}
