package cosmic

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/xanzy/go-cosmic/cosmic"
)

func resourceCosmicVPC() *schema.Resource {
	return &schema.Resource{
		Create: resourceCosmicVPCCreate,
		Read:   resourceCosmicVPCRead,
		Update: resourceCosmicVPCUpdate,
		Delete: resourceCosmicVPCDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"display_text": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"cidr": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vpc_offering": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"network_domain": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"source_nat_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"source_nat_list": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"syslog_server_list": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"zone": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCosmicVPCCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cosmic.CosmicClient)

	name := d.Get("name").(string)

	// Retrieve the vpc_offering ID
	vpcofferingid, e := retrieveID(cs, "vpc_offering", d.Get("vpc_offering").(string))
	if e != nil {
		return e.Error()
	}

	// Retrieve the zone ID
	zoneid, e := retrieveID(cs, "zone", d.Get("zone").(string))
	if e != nil {
		return e.Error()
	}

	// Set the display text
	displaytext, ok := d.GetOk("display_text")
	if !ok {
		displaytext = name
	}

	// Create a new parameter struct
	p := cs.VPC.NewCreateVPCParams(
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

	// If there is a project supplied, we retrieve and set the project id
	if err := setProjectid(p, cs, d); err != nil {
		return err
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
	r, err := cs.VPC.CreateVPC(p)
	if err != nil {
		return fmt.Errorf("Error creating VPC %s: %s", name, err)
	}

	d.SetId(r.Id)

	return resourceCosmicVPCRead(d, meta)
}

func resourceCosmicVPCRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cosmic.CosmicClient)

	// Get the VPC details
	v, count, err := cs.VPC.GetVPCByID(
		d.Id(),
		cosmic.WithProject(d.Get("project").(string)),
	)
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
	o, _, err := cs.VPC.GetVPCOfferingByID(v.Vpcofferingid)
	if err != nil {
		return err
	}

	setValueOrID(d, "vpc_offering", o.Name, v.Vpcofferingid)
	setValueOrID(d, "project", v.Project, v.Projectid)
	setValueOrID(d, "zone", v.Zonename, v.Zoneid)

	// Create a new parameter struct
	p := cs.PublicIPAddress.NewListPublicIpAddressesParams()
	p.SetVpcid(d.Id())
	p.SetIssourcenat(true)

	// If there is a project supplied, we retrieve and set the project id
	if err := setProjectid(p, cs, d); err != nil {
		return err
	}

	// Get the source NAT IP assigned to the VPC
	l, err := cs.PublicIPAddress.ListPublicIpAddresses(p)
	if err != nil {
		return err
	}

	if l.Count == 1 {
		d.Set("source_nat_ip", l.PublicIpAddresses[0].Ipaddress)
	}

	return nil
}

func resourceCosmicVPCUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cosmic.CosmicClient)

	name := d.Get("name").(string)

	// Create a new parameter struct
	p := cs.VPC.NewUpdateVPCParams(d.Id())

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

	// Update the VPC
	_, err := cs.VPC.UpdateVPC(p)
	if err != nil {
		return fmt.Errorf("Error updating name of VPC %s: %s", name, err)
	}

	return resourceCosmicVPCRead(d, meta)
}

func resourceCosmicVPCDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cosmic.CosmicClient)

	// Create a new parameter struct
	p := cs.VPC.NewDeleteVPCParams(d.Id())

	// Delete the VPC
	_, err := cs.VPC.DeleteVPC(p)
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
