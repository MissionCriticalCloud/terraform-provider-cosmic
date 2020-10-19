package cosmic

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCosmicVPNConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceCosmicVPNConnectionCreate,
		Read:   resourceCosmicVPNConnectionRead,
		Delete: resourceCosmicVPNConnectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"customer_gateway_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vpn_gateway_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCosmicVPNConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	// Create a new parameter struct
	p := client.VPN.NewCreateVpnConnectionParams(
		d.Get("customer_gateway_id").(string),
		d.Get("vpn_gateway_id").(string),
	)

	// Create the new VPN Connection
	v, err := client.VPN.CreateVpnConnection(p)
	if err != nil {
		return fmt.Errorf("Error creating VPN Connection: %s", err)
	}

	d.SetId(v.Id)

	return resourceCosmicVPNConnectionRead(d, meta)
}

func resourceCosmicVPNConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	// Get the VPN Connection details
	v, count, err := client.VPN.GetVpnConnectionByID(d.Id())
	if err != nil {
		if count == 0 {
			log.Printf("[DEBUG] VPN Connection does no longer exist")
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("customer_gateway_id", v.S2scustomergatewayid)
	d.Set("vpn_gateway_id", v.S2svpngatewayid)

	return nil
}

func resourceCosmicVPNConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	// Create a new parameter struct
	p := client.VPN.NewDeleteVpnConnectionParams(d.Id())

	// Delete the VPN Connection
	_, err := client.VPN.DeleteVpnConnection(p)
	if err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", d.Id())) {
			return nil
		}

		return fmt.Errorf("Error deleting VPN Connection: %s", err)
	}

	return nil
}
