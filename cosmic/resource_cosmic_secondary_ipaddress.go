package cosmic

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCosmicSecondaryIPAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceCosmicSecondaryIPAddressCreate,
		Read:   resourceCosmicSecondaryIPAddressRead,
		Delete: resourceCosmicSecondaryIPAddressDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceCosmicSecondaryIPAddressImporter,
		},

		Schema: map[string]*schema.Schema{
			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"nic_id": {
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

func resourceCosmicSecondaryIPAddressCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	nicid, ok := d.GetOk("nic_id")
	if !ok {
		virtualmachineid := d.Get("virtual_machine_id").(string)

		// Get the virtual machine details
		vm, count, err := client.VirtualMachine.GetVirtualMachineByID(virtualmachineid)
		if err != nil {
			if count == 0 {
				log.Printf("[DEBUG] Virtual Machine %s does not exist", virtualmachineid)
				d.SetId("")
				return nil
			}
			return err
		}

		nicid = vm.Nic[0].Id
	}

	// Create a new parameter struct
	p := client.Nic.NewAddIpToNicParams(nicid.(string))

	// If there is a ipaddres supplied, add it to the parameter struct
	if ipaddress, ok := d.GetOk("ip_address"); ok {
		p.SetIpaddress(ipaddress.(string))
	}

	ip, err := client.Nic.AddIpToNic(p)
	if err != nil {
		return err
	}

	d.SetId(ip.Id)

	return nil
}

func resourceCosmicSecondaryIPAddressRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	virtualmachineid := d.Get("virtual_machine_id").(string)

	// Get the virtual machine details
	vm, count, err := client.VirtualMachine.GetVirtualMachineByID(virtualmachineid)
	if err != nil {
		if count == 0 {
			log.Printf("[DEBUG] Virtual Machine %s does not exist", virtualmachineid)
			d.SetId("")
			return nil
		}
		return err
	}

	nicid, ok := d.GetOk("nic_id")
	if !ok {
		nicid = vm.Nic[0].Id
	}

	p := client.Nic.NewListNicsParams(virtualmachineid)
	p.SetNicid(nicid.(string))

	l, err := client.Nic.ListNics(p)
	if err != nil {
		return err
	}

	if l.Count == 0 {
		log.Printf("[DEBUG] NIC %s does not exist", d.Get("nic_id").(string))
		d.SetId("")
		return nil
	}

	if l.Count > 1 {
		return fmt.Errorf("Found more then one possible result: %v", l.Nics)
	}

	for _, ip := range l.Nics[0].Secondaryip {
		if ip.Id == d.Id() {
			d.Set("ip_address", ip.Ipaddress)
			d.Set("nic_id", l.Nics[0].Id)
			d.Set("virtual_machine_id", l.Nics[0].Virtualmachineid)
			return nil
		}
	}

	log.Printf("[DEBUG] IP %s no longer exist", d.Get("ip_address").(string))
	d.SetId("")

	return nil
}

func resourceCosmicSecondaryIPAddressDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	// Create a new parameter struct
	p := client.Nic.NewRemoveIpFromNicParams(d.Id())

	log.Printf("[INFO] Removing secondary IP address: %s", d.Get("ip_address").(string))
	if _, err := client.Nic.RemoveIpFromNic(p); err != nil {
		// This is a very poor way to be told the ID does not exist :(
		if strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", d.Id())) {
			return nil
		}

		return fmt.Errorf("Error removing secondary IP address: %s", err)
	}

	return nil
}

func resourceCosmicSecondaryIPAddressImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	s := strings.Split(d.Id(), "/")
	if len(s) != 2 {
		return nil, fmt.Errorf(
			"invalid variable import format: %s (expected <INSTANCE ID>/<SECONDARY IP ADDRESS>)",
			d.Id(),
		)
	}
	vmid, sip := s[0], s[1]

	client := meta.(*CosmicClient)

	// Get the virtual machine details
	vm, count, err := client.VirtualMachine.GetVirtualMachineByID(vmid)
	if err != nil {
		if count == 0 {
			return nil, fmt.Errorf("[DEBUG] Virtual Machine %s does not exist", vmid)
		}
		return nil, err
	}

	for _, n := range vm.Nic {
		for _, ip := range n.Secondaryip {
			if ip.Ipaddress == sip {
				d.SetId(ip.Id)
				d.Set("ip_address", ip.Ipaddress)
				d.Set("nic_id", n.Id)
				d.Set("virtual_machine_id", vm.Id)
				return []*schema.ResourceData{d}, nil
			}
		}
	}

	return nil, fmt.Errorf("IP address %s does not exist", sip)
}
