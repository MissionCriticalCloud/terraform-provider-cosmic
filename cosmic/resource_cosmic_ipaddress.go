package cosmic

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCosmicIPAddress() *schema.Resource {
	aclidSchema := &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Default:  none,
	}

	aclidSchema.StateFunc = func(v interface{}) string {
		value := v.(string)

		if value == none {
			aclidSchema.ForceNew = true
		} else {
			aclidSchema.ForceNew = false
		}

		return value
	}

	return &schema.Resource{
		Create: resourceCosmicIPAddressCreate,
		Read:   resourceCosmicIPAddressRead,
		Update: resourceCosmicIPAddressUpdate,
		Delete: resourceCosmicIPAddressDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCosmicIPAddressImporter,
		},

		Schema: map[string]*schema.Schema{
			"network_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"acl_id": aclidSchema,
		},
	}
}

func resourceCosmicIPAddressCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	if err := verifyIPAddressParams(d); err != nil {
		return err
	}

	// Create a new parameter struct
	p := client.PublicIPAddress.NewAssociateIpAddressParams()

	if networkid, ok := d.GetOk("network_id"); ok {
		// Set the networkid
		p.SetNetworkid(networkid.(string))
	}

	if vpcid, ok := d.GetOk("vpc_id"); ok {
		// Set the vpcid
		p.SetVpcid(vpcid.(string))
	}

	// Associate a new IP address
	r, err := client.PublicIPAddress.AssociateIpAddress(p)
	if err != nil {
		return fmt.Errorf("Error associating a new IP address: %s", err)
	}

	d.SetId(r.Id)

	// Set the ACL if we are on a VPC and acl_id is supplied
	if _, ok := d.GetOk("vpc_id"); ok {
		if aclid, ok := d.GetOk("acl_id"); ok && aclid.(string) != none {
			p := client.NetworkACL.NewReplaceNetworkACLListParams(aclid.(string))
			p.SetPublicipid(d.Id())

			_, err := client.NetworkACL.ReplaceNetworkACLList(p)
			if err != nil {
				return fmt.Errorf("Error replacing ACL: %s", err)
			}
		}
	}

	return resourceCosmicIPAddressRead(d, meta)
}

func resourceCosmicIPAddressRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	// Get the IP address details
	ip, count, err := client.PublicIPAddress.GetPublicIpAddressByID(d.Id())
	if err != nil {
		if count == 0 {
			log.Printf(
				"[DEBUG] IP address with ID %s is no longer associated", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	// Updated the IP address
	d.Set("ip_address", ip.Ipaddress)

	if _, ok := d.GetOk("network_id"); ok {
		d.Set("network_id", ip.Associatednetworkid)
	}

	if _, ok := d.GetOk("vpc_id"); ok {
		d.Set("vpc_id", ip.Vpcid)
	}

	if ip.Aclid == "" {
		ip.Aclid = none
	}
	d.Set("acl_id", ip.Aclid)

	return nil
}

func resourceCosmicIPAddressUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	// Replace the ACL if the ID has changed
	if d.HasChange("acl_id") {
		p := client.NetworkACL.NewReplaceNetworkACLListParams(d.Get("acl_id").(string))
		p.SetPublicipid(d.Id())

		_, err := client.NetworkACL.ReplaceNetworkACLList(p)
		if err != nil {
			return fmt.Errorf("Error replacing ACL: %s", err)
		}
	}

	return resourceCosmicIPAddressRead(d, meta)
}

func resourceCosmicIPAddressDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	// Create a new parameter struct
	p := client.PublicIPAddress.NewDisassociateIpAddressParams(d.Id())

	// Disassociate the IP address
	if _, err := client.PublicIPAddress.DisassociateIpAddress(p); err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", d.Id())) {
			return nil
		}

		return fmt.Errorf("Error disassociating IP address %s: %s", d.Get("name").(string), err)
	}

	return nil
}

func resourceCosmicIPAddressImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*CosmicClient)
	ip, _, _ := client.PublicIPAddress.GetPublicIpAddressByID(d.Id())

	// Set the vpc_id if the IP is attached to a VPC.
	if ip.Vpcid != "" {
		d.Set("vpc_id", ip.Vpcid)
	}

	return []*schema.ResourceData{d}, nil
}

func verifyIPAddressParams(d *schema.ResourceData) error {
	_, network := d.GetOk("network_id")
	_, vpc := d.GetOk("vpc_id")

	if (network && vpc) || (!network && !vpc) {
		return fmt.Errorf(
			"You must supply a value for either (so not both) the 'network_id' or 'vpc_id' parameter")
	}

	return nil
}
