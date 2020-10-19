package cosmic

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCosmicVPC() *schema.Resource {
	return &schema.Resource{
		Create: resourceCosmicVPCCreate,
		Read:   resourceCosmicVPCRead,
		Update: resourceCosmicVPCUpdate,
		Delete: resourceCosmicVPCDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"display_text": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"cidr": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vpc_offering": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},

			"network_domain": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"source_nat_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"source_nat_ip_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"source_nat_list": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"syslog_server_list": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"zone": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				ForceNew:   true,
				Deprecated: deprecatedZoneMsg(),
			},
		},
	}
}

func resourceCosmicVPCCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	name := d.Get("name").(string)

	// Retrieve the vpc_offering ID
	vpcofferingid, e := retrieveID(client, "vpc_offering", d.Get("vpc_offering").(string))
	if e != nil {
		return e.Error()
	}

	// Retrieve the zone ID
	zoneid, e := retrieveID(client, "zone", client.ZoneName)
	if e != nil {
		return e.Error()
	}

	// Set the display text
	displaytext, ok := d.GetOk("display_text")
	if !ok {
		displaytext = name
	}

	// Create a new parameter struct
	p := client.VPC.NewCreateVPCParams(
		d.Get("cidr").(string),
		displaytext.(string),
		name,
		vpcofferingid,
		zoneid,
	)

	// If there is a network domain supplied, make sure to add it to the request
	if networkDomain, ok := d.GetOk("network_domain"); ok {
		// Set the network domain
		p.SetNetworkdomain(networkDomain.(string))
	}

	// If there is a sourcenatlist supplied, make sure to add it to the request
	if sourceNatList, ok := d.GetOk("source_nat_list"); ok {
		// Set the Source NAT list
		p.SetSourcenatlist(sourceNatList.(string))
	}

	// If there is a syslogserverlist supplied, make sure to add it to the request
	if syslogServerList, ok := d.GetOk("syslog_server_list"); ok {
		// Set the syslog server list
		p.SetSyslogserverlist(syslogServerList.(string))
	}

	// Create the new VPC
	r, err := client.VPC.CreateVPC(p)
	if err != nil {
		return fmt.Errorf("Error creating VPC %s: %s", name, err)
	}

	d.SetId(r.Id)

	return resourceCosmicVPCRead(d, meta)
}

func resourceCosmicVPCRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	// Get the VPC details
	v, count, err := client.VPC.GetVPCByID(d.Id())
	if err != nil {
		if count == 0 {
			log.Printf(
				"[DEBUG] VPC %s does no longer exist", d.Get("name").(string))
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("name", v.Name)
	d.Set("display_text", v.Displaytext)
	d.Set("cidr", v.Cidr)
	d.Set("network_domain", v.Networkdomain)
	d.Set("sourcenatlist", v.Sourcenatlist)
	d.Set("syslogserverlist", v.Syslogserverlist)

	// Get the VPC offering details
	o, _, err := client.VPC.GetVPCOfferingByID(v.Vpcofferingid)
	if err != nil {
		return err
	}

	setValueOrID(d, "vpc_offering", o.Name, v.Vpcofferingid)
	setValueOrID(d, "zone", v.Zonename, v.Zoneid)

	// Create a new parameter struct
	p := client.PublicIPAddress.NewListPublicIpAddressesParams()
	p.SetVpcid(d.Id())
	p.SetIssourcenat(true)

	// Get the source NAT IP assigned to the VPC
	l, err := client.PublicIPAddress.ListPublicIpAddresses(p)
	if err != nil {
		return err
	}

	if l.Count == 1 {
		d.Set("source_nat_ip", l.PublicIpAddresses[0].Ipaddress)
		d.Set("source_nat_ip_id", l.PublicIpAddresses[0].Id)
	}

	return nil
}

func resourceCosmicVPCUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	name := d.Get("name").(string)

	// Create a new parameter struct
	p := client.VPC.NewUpdateVPCParams(d.Id())

	// Check if the name is changed
	if d.HasChange("name") {
		// Set the new name
		p.SetName(name)
	}

	// Check if the display text is changed
	if d.HasChange("display_text") {
		// Set the display text
		displaytext, ok := d.GetOk("display_text")
		if !ok {
			displaytext = name
		}

		// Set the new display text
		p.SetDisplaytext(displaytext.(string))
	}

	// Check if the source nat list is changed
	if d.HasChange("source_nat_list") {
		// Set the source nat list
		sourcenatlist := d.Get("source_nat_list")

		// Set the new display text
		p.SetSourcenatlist(sourcenatlist.(string))
	}

	// Check if the syslog server list is changed
	if d.HasChange("syslog_server_list") {
		// Set the syslog server list
		syslogserverlist := d.Get("syslog_server_list")

		// Set the new display text
		p.SetSyslogserverlist(syslogserverlist.(string))
	}

	// Check if the VPC offering is changed
	if d.HasChange("vpc_offering") {
		// Retrieve the VPC offering ID
		o, _, err := client.VPC.GetVPCOfferingByName(d.Get("vpc_offering").(string))
		if err != nil {
			return err
		}
		// Set the new VPC offering ID
		p.SetVpcofferingid(o.Id)
	}

	// Update the VPC
	_, err := client.VPC.UpdateVPC(p)
	if err != nil {
		return fmt.Errorf("Error updating name of VPC %s: %s", name, err)
	}

	return resourceCosmicVPCRead(d, meta)
}

func resourceCosmicVPCDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	// Create a new parameter struct
	p := client.VPC.NewDeleteVPCParams(d.Id())

	// Delete the VPC
	_, err := client.VPC.DeleteVPC(p)
	if err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", d.Id())) {
			return nil
		}

		return fmt.Errorf("Error deleting VPC %s: %s", d.Get("name").(string), err)
	}

	return nil
}
