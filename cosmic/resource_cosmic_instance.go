package cosmic

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCosmicInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceCosmicInstanceCreate,
		Read:   resourceCosmicInstanceRead,
		Update: resourceCosmicInstanceUpdate,
		Delete: resourceCosmicInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"service_offering": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},

			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"template": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"root_disk_size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"group": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"affinity_group_ids": {
				Type:          schema.TypeSet,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Set:           schema.HashString,
				ConflictsWith: []string{"affinity_group_names"},
			},

			"affinity_group_names": {
				Type:          schema.TypeSet,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Set:           schema.HashString,
				ConflictsWith: []string{"affinity_group_ids"},
			},

			"zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"keypair": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				StateFunc: func(v interface{}) string {
					switch v.(type) {
					case string:
						hash := sha1.Sum([]byte(v.(string)))
						return hex.EncodeToString(hash[:])
					default:
						return ""
					}
				},
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

			"optimise_for": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				StateFunc: func(val interface{}) string {
					return strings.Title(strings.ToLower(val.(string)))
				},
			},

			"expunge": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceCosmicInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	// Retrieve the service_offering ID
	serviceofferingid, e := retrieveID(client, "service_offering", d.Get("service_offering").(string))
	if e != nil {
		return e.Error()
	}

	// Retrieve the zone ID
	zoneid, e := retrieveID(client, "zone", d.Get("zone").(string))
	if e != nil {
		return e.Error()
	}

	// Retrieve the zone object
	zone, _, err := client.Zone.GetZoneByID(zoneid)
	if err != nil {
		return err
	}

	// Retrieve the template ID
	templateid, e := retrieveTemplateID(client, zone.Id, d.Get("template").(string))
	if e != nil {
		return e.Error()
	}

	// Create a new parameter struct
	p := client.VirtualMachine.NewDeployVirtualMachineParams(serviceofferingid, templateid, zone.Id)

	// Set the name
	name, hasName := d.GetOk("name")
	if hasName {
		p.SetName(name.(string))
	}

	// Set the display name
	if displayname, ok := d.GetOk("display_name"); ok {
		p.SetDisplayname(displayname.(string))
	} else if hasName {
		p.SetDisplayname(name.(string))
	}

	// If there is a root_disk_size supplied, add it to the parameter struct
	if rootdisksize, ok := d.GetOk("root_disk_size"); ok {
		p.SetRootdisksize(int64(rootdisksize.(int)))
	}

	if zone.Networktype == "Advanced" {
		// Set the default network ID
		p.SetNetworkids([]string{d.Get("network_id").(string)})
	}

	// If there is a ipaddres supplied, add it to the parameter struct
	if ipaddress, ok := d.GetOk("ip_address"); ok {
		p.SetIpaddress(ipaddress.(string))
	}

	// Set the disk_controller type if configured
	if diskcontroller, ok := d.GetOk("disk_controller"); ok {
		p.SetDiskcontroller(diskcontroller.(string))
	}

	// If optimise_for is supplied add it to the parameter struct
	if optimise_for, ok := d.GetOk("optimise_for"); ok {
		// Param must be lowercase and capitalized
		p.SetOptimisefor(strings.Title(strings.ToLower(optimise_for.(string))))
	}

	// If there is a group supplied, add it to the parameter struct
	if group, ok := d.GetOk("group"); ok {
		p.SetGroup(group.(string))
	}

	// If there are affinity group IDs supplied, add them to the parameter struct
	if agIDs := d.Get("affinity_group_ids").(*schema.Set); agIDs.Len() > 0 {
		var groups []string
		for _, group := range agIDs.List() {
			groups = append(groups, group.(string))
		}
		p.SetAffinitygroupids(groups)
	}

	// If there are affinity group names supplied, add them to the parameter struct
	if agNames := d.Get("affinity_group_names").(*schema.Set); agNames.Len() > 0 {
		var groups []string
		for _, group := range agNames.List() {
			groups = append(groups, group.(string))
		}
		p.SetAffinitygroupnames(groups)
	}

	// If a keypair is supplied, add it to the parameter struct
	if keypair, ok := d.GetOk("keypair"); ok {
		p.SetKeypair(keypair.(string))
	}

	if userData, ok := d.GetOk("user_data"); ok {
		ud, err := getUserData(userData.(string), client.HTTPGETOnly)
		if err != nil {
			return err
		}

		p.SetUserdata(ud)
	}

	// Create the new instance
	r, err := client.VirtualMachine.DeployVirtualMachine(p)
	if err != nil {
		return fmt.Errorf("Error creating the new instance %s: %s", name, err)
	}

	d.SetId(r.Id)

	// Set the connection info for any configured provisioners
	d.SetConnInfo(map[string]string{
		"host":     r.Nic[0].Ipaddress,
		"password": r.Password,
	})

	return resourceCosmicInstanceRead(d, meta)
}

func resourceCosmicInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	// Get the virtual machine details
	vm, count, err := client.VirtualMachine.GetVirtualMachineByID(d.Id())
	if err != nil {
		if count == 0 {
			log.Printf("[DEBUG] Instance %s does no longer exist", d.Get("name").(string))
			d.SetId("")
			return nil
		}

		return err
	}

	// Update the config
	d.Set("name", vm.Name)
	d.Set("display_name", vm.Displayname)
	d.Set("group", vm.Group)
	d.Set("disk_controller", vm.Rootdevicecontroller)
	d.Set("optimise_for", vm.Optimisefor)

	// In some rare cases (when destroying a machine failes) it can happen that
	// an instance does not have any attached NIC anymore.
	if len(vm.Nic) > 0 {
		d.Set("network_id", vm.Nic[0].Networkid)
		d.Set("ip_address", vm.Nic[0].Ipaddress)
	}

	if _, ok := d.GetOk("affinity_group_ids"); ok {
		groups := &schema.Set{F: schema.HashString}
		for _, group := range vm.Affinitygroup {
			groups.Add(group.Id)
		}
		d.Set("affinity_group_ids", groups)
	}

	if _, ok := d.GetOk("affinity_group_names"); ok {
		groups := &schema.Set{F: schema.HashString}
		for _, group := range vm.Affinitygroup {
			groups.Add(group.Name)
		}
		d.Set("affinity_group_names", groups)
	}

	setValueOrID(d, "service_offering", vm.Serviceofferingname, vm.Serviceofferingid)
	setValueOrID(d, "template", vm.Templatename, vm.Templateid)
	setValueOrID(d, "zone", vm.Zonename, vm.Zoneid)

	return nil
}

func resourceCosmicInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)
	d.Partial(true)

	name := d.Get("name").(string)

	// Check if the display name is changed and if so, update the virtual machine
	if d.HasChange("display_name") {
		log.Printf("[DEBUG] Display name changed for %s, starting update", name)

		// Create a new parameter struct
		p := client.VirtualMachine.NewUpdateVirtualMachineParams(d.Id())

		// Set the new display name
		p.SetDisplayname(d.Get("display_name").(string))

		// Update the display name
		_, err := client.VirtualMachine.UpdateVirtualMachine(p)
		if err != nil {
			return fmt.Errorf(
				"Error updating the display name for instance %s: %s", name, err)
		}

		d.SetPartial("display_name")
	}

	// Check if the group is changed and if so, update the virtual machine
	if d.HasChange("group") {
		log.Printf("[DEBUG] Group changed for %s, starting update", name)

		// Create a new parameter struct
		p := client.VirtualMachine.NewUpdateVirtualMachineParams(d.Id())

		// Set the new group
		p.SetGroup(d.Get("group").(string))

		// Update the display name
		_, err := client.VirtualMachine.UpdateVirtualMachine(p)
		if err != nil {
			return fmt.Errorf(
				"Error updating the group for instance %s: %s", name, err)
		}

		d.SetPartial("group")
	}

	// Attributes that require reboot to update
	if d.HasChange("name") || d.HasChange("service_offering") || d.HasChange("affinity_group_ids") ||
		d.HasChange("affinity_group_names") || d.HasChange("keypair") || d.HasChange("user_data") ||
		d.HasChange("optimise_for") {
		// Before we can actually make these changes, the virtual machine must be stopped
		_, err := client.VirtualMachine.StopVirtualMachine(
			client.VirtualMachine.NewStopVirtualMachineParams(d.Id()))
		if err != nil {
			return fmt.Errorf(
				"Error stopping instance %s before making changes: %s", name, err)
		}

		// Check if the name has changed and if so, update the name
		if d.HasChange("name") {
			log.Printf("[DEBUG] Name for %s changed to %s, starting update", d.Id(), name)

			// Create a new parameter struct
			p := client.VirtualMachine.NewUpdateVirtualMachineParams(d.Id())

			// Set the new name
			p.SetName(name)

			// Update the display name
			_, err := client.VirtualMachine.UpdateVirtualMachine(p)
			if err != nil {
				return fmt.Errorf(
					"Error updating the name for instance %s: %s", name, err)
			}

			d.SetPartial("name")
		}

		// Check if the service offering is changed and if so, update the offering
		if d.HasChange("service_offering") {
			log.Printf("[DEBUG] Service offering changed for %s, starting update", name)

			// Retrieve the service_offering ID
			serviceofferingid, e := retrieveID(client, "service_offering", d.Get("service_offering").(string))
			if e != nil {
				return e.Error()
			}

			// Create a new parameter struct
			p := client.VirtualMachine.NewChangeServiceForVirtualMachineParams(d.Id(), serviceofferingid)

			// Change the service offering
			_, err = client.VirtualMachine.ChangeServiceForVirtualMachine(p)
			if err != nil {
				return fmt.Errorf(
					"Error changing the service offering for instance %s: %s", name, err)
			}
			d.SetPartial("service_offering")
		}

		// Check if the affinity group IDs have changed and if so, update the IDs
		if d.HasChange("affinity_group_ids") {
			p := client.AffinityGroup.NewUpdateVMAffinityGroupParams(d.Id())
			var groups []string

			if agIDs := d.Get("affinity_group_ids").(*schema.Set); agIDs.Len() > 0 {
				for _, group := range agIDs.List() {
					groups = append(groups, group.(string))
				}
			}

			// Set the new groups
			p.SetAffinitygroupids(groups)

			// Update the affinity groups
			_, err = client.AffinityGroup.UpdateVMAffinityGroup(p)
			if err != nil {
				return fmt.Errorf(
					"Error updating the affinity groups for instance %s: %s", name, err)
			}
			d.SetPartial("affinity_group_ids")
		}

		// Check if the affinity group names have changed and if so, update the names
		if d.HasChange("affinity_group_names") {
			p := client.AffinityGroup.NewUpdateVMAffinityGroupParams(d.Id())
			var groups []string

			if agNames := d.Get("affinity_group_names").(*schema.Set); agNames.Len() > 0 {
				for _, group := range agNames.List() {
					groups = append(groups, group.(string))
				}
			}

			// Set the new groups
			p.SetAffinitygroupnames(groups)

			// Update the affinity groups
			_, err = client.AffinityGroup.UpdateVMAffinityGroup(p)
			if err != nil {
				return fmt.Errorf(
					"Error updating the affinity groups for instance %s: %s", name, err)
			}
			d.SetPartial("affinity_group_names")
		}

		// Check if the keypair has changed and if so, update the keypair
		if d.HasChange("keypair") {
			log.Printf("[DEBUG] SSH keypair changed for %s, starting update", name)

			p := client.SSH.NewResetSSHKeyForVirtualMachineParams(d.Id(), d.Get("keypair").(string))

			// Change the ssh keypair
			_, err = client.SSH.ResetSSHKeyForVirtualMachine(p)
			if err != nil {
				return fmt.Errorf(
					"Error changing the SSH keypair for instance %s: %s", name, err)
			}
			d.SetPartial("keypair")
		}

		// Check if the user data has changed and if so, update the user data
		if d.HasChange("user_data") {
			log.Printf("[DEBUG] user_data changed for %s, starting update", name)

			ud, err := getUserData(d.Get("user_data").(string), client.HTTPGETOnly)
			if err != nil {
				return err
			}

			p := client.VirtualMachine.NewUpdateVirtualMachineParams(d.Id())
			p.SetUserdata(ud)
			_, err = client.VirtualMachine.UpdateVirtualMachine(p)
			if err != nil {
				return fmt.Errorf(
					"Error updating user_data for instance %s: %s", name, err)
			}
			d.SetPartial("user_data")
		}

		// Check if the user data has changed and if so, update the user data
		if d.HasChange("optimise_for") {
			log.Printf("[DEBUG] optimise_for changed for %s, starting update", name)

			ud, err := getUserData(d.Get("optimise_for").(string), client.HTTPGETOnly)
			if err != nil {
				return err
			}

			p := client.VirtualMachine.NewUpdateVirtualMachineParams(d.Id())
			p.SetOptimisefor(ud)
			_, err = client.VirtualMachine.UpdateVirtualMachine(p)
			if err != nil {
				return fmt.Errorf(
					"Error updating optimise_for for instance %s: %s", name, err)
			}
			d.SetPartial("optimise_for")
		}

		// Start the virtual machine again
		_, err = client.VirtualMachine.StartVirtualMachine(
			client.VirtualMachine.NewStartVirtualMachineParams(d.Id()))
		if err != nil {
			return fmt.Errorf(
				"Error starting instance %s after making changes", name)
		}
	}

	d.Partial(false)

	return resourceCosmicInstanceRead(d, meta)
}

func resourceCosmicInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	// Create a new parameter struct
	p := client.VirtualMachine.NewDestroyVirtualMachineParams(d.Id())

	if d.Get("expunge").(bool) {
		p.SetExpunge(true)
	}

	log.Printf("[INFO] Destroying instance: %s", d.Get("name").(string))
	if _, err := client.VirtualMachine.DestroyVirtualMachine(p); err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", d.Id())) {
			return nil
		}

		return fmt.Errorf("Error destroying instance: %s", err)
	}

	return nil
}

// getUserData returns the user data as a base64 encoded string
func getUserData(userData string, httpGetOnly bool) (string, error) {
	ud := base64.StdEncoding.EncodeToString([]byte(userData))

	// deployVirtualMachine uses POST by default, so max userdata is 32K
	maxUD := 32768

	if httpGetOnly {
		// deployVirtualMachine using GET instead, so max userdata is 2K
		maxUD = 2048
	}

	if len(ud) > maxUD {
		return "", fmt.Errorf(
			"The supplied user_data contains %d bytes after encoding, "+
				"this exeeds the limit of %d bytes", len(ud), maxUD)
	}

	return ud, nil
}
