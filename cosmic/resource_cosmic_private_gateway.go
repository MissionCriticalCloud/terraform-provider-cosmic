package cosmic

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCosmicPrivateGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceCosmicPrivateGatewayCreate,
		Read:   resourceCosmicPrivateGatewayRead,
		Delete: resourceCosmicPrivateGatewayDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"ip_address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"acl_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCosmicPrivateGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	ipaddress := d.Get("ip_address").(string)

	// Create a new parameter struct
	p := client.VPC.NewCreatePrivateGatewayParams(
		ipaddress,
		d.Get("network_id").(string),
		d.Get("vpc_id").(string),
	)

	// Set the acl ID
	p.SetAclid(d.Get("acl_id").(string))

	// Create the new private gateway
	r, err := client.VPC.CreatePrivateGateway(p)
	if err != nil {
		return fmt.Errorf("Error creating private gateway for %s: %s", ipaddress, err)
	}

	d.SetId(r.Id)

	return resourceCosmicPrivateGatewayRead(d, meta)
}

func resourceCosmicPrivateGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	// Get the private gateway details
	gw, count, err := client.VPC.GetPrivateGatewayByID(d.Id())
	if err != nil {
		if count == 0 {
			log.Printf("[DEBUG] Private gateway %s does no longer exist", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("ip_address", gw.Ipaddress)
	d.Set("network_id", gw.Networkid)
	d.Set("acl_id", gw.Aclid)
	d.Set("vpc_id", gw.Vpcid)

	return nil
}

func resourceCosmicPrivateGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	// Create a new parameter struct
	p := client.VPC.NewDeletePrivateGatewayParams(d.Id())

	// Delete the private gateway
	_, err := client.VPC.DeletePrivateGateway(p)
	if err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", d.Id())) {
			return nil
		}

		return fmt.Errorf("Error deleting private gateway %s: %s", d.Id(), err)
	}

	return nil
}
