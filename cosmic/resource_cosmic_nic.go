package cosmic

import (
	"fmt"
	"log"
	"strings"

	"github.com/MissionCriticalCloud/go-cosmic/v6/cosmic"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCosmicNIC() *schema.Resource {
	return &schema.Resource{
		Create: resourceCosmicNICCreate,
		Read:   resourceCosmicNICRead,
		Delete: resourceCosmicNICDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"virtual_machine_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCosmicNICCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	// Create a new parameter struct
	p := client.VirtualMachine.NewAddNicToVirtualMachineParams(
		d.Get("network_id").(string),
		d.Get("virtual_machine_id").(string),
	)

	// If there is a ipaddres supplied, add it to the parameter struct
	if ipaddress, ok := d.GetOk("ip_address"); ok {
		p.SetIpaddress(ipaddress.(string))
	}

	// Create and attach the new NIC
	r, err := Retry(10, retryableAddNicFunc(client, p))
	if err != nil {
		return fmt.Errorf("Error creating the new NIC: %s", err)
	}

	found := false
	for _, n := range r.(*cosmic.AddNicToVirtualMachineResponse).Nic {
		if n.Networkid == d.Get("network_id").(string) {
			d.SetId(n.Id)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("Could not find NIC ID for network ID: %s", d.Get("network_id").(string))
	}

	return resourceCosmicNICRead(d, meta)
}

func resourceCosmicNICRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	// Get the virtual machine details
	vm, count, err := client.VirtualMachine.GetVirtualMachineByID(d.Get("virtual_machine_id").(string))
	if err != nil {
		if count == 0 {
			log.Printf("[DEBUG] Instance %s does no longer exist", d.Get("virtual_machine_id").(string))
			d.SetId("")
			return nil
		}

		return err
	}

	// Read NIC info
	found := false
	for _, n := range vm.Nic {
		if n.Id == d.Id() {
			d.Set("ip_address", n.Ipaddress)
			d.Set("network_id", n.Networkid)
			d.Set("virtual_machine_id", vm.Id)
			found = true
			break
		}
	}

	if !found {
		log.Printf("[DEBUG] NIC for network ID %s does no longer exist", d.Get("network_id").(string))
		d.SetId("")
	}

	return nil
}

func resourceCosmicNICDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	// Create a new parameter struct
	p := client.VirtualMachine.NewRemoveNicFromVirtualMachineParams(
		d.Id(),
		d.Get("virtual_machine_id").(string),
	)

	// Remove the NIC
	_, err := client.VirtualMachine.RemoveNicFromVirtualMachine(p)
	if err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", d.Id())) {
			return nil
		}

		return fmt.Errorf("Error deleting NIC: %s", err)
	}

	return nil
}

func retryableAddNicFunc(client *CosmicClient, p *cosmic.AddNicToVirtualMachineParams) func() (interface{}, error) {
	return func() (interface{}, error) {
		r, err := client.VirtualMachine.AddNicToVirtualMachine(p)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
}
