package cosmic

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCosmicVPNCustomerGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceCosmicVPNCustomerGatewayCreate,
		Read:   resourceCosmicVPNCustomerGatewayRead,
		Update: resourceCosmicVPNCustomerGatewayUpdate,
		Delete: resourceCosmicVPNCustomerGatewayDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"cidr_list": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"esp_policy": {
				Type:     schema.TypeString,
				Required: true,
			},

			"gateway": {
				Type:     schema.TypeString,
				Required: true,
			},

			"ike_policy": {
				Type:     schema.TypeString,
				Required: true,
			},

			"ipsec_psk": {
				Type:     schema.TypeString,
				Required: true,
			},

			"dpd": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"esp_lifetime": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"ike_lifetime": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceCosmicVPNCustomerGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	// Create a new parameter struct
	p := client.VPN.NewCreateVpnCustomerGatewayParams(
		createCidrList(d.Get("cidr_list").(*schema.Set)),
		d.Get("esp_policy").(string),
		d.Get("gateway").(string),
		d.Get("ike_policy").(string),
		d.Get("ipsec_psk").(string),
	)

	p.SetName(d.Get("name").(string))

	if dpd, ok := d.GetOk("dpd"); ok {
		p.SetDpd(dpd.(bool))
	}

	if esplifetime, ok := d.GetOk("esp_lifetime"); ok {
		p.SetEsplifetime(int64(esplifetime.(int)))
	}

	if ikelifetime, ok := d.GetOk("ike_lifetime"); ok {
		p.SetIkelifetime(int64(ikelifetime.(int)))
	}

	// Create the new VPN Customer Gateway
	v, err := client.VPN.CreateVpnCustomerGateway(p)
	if err != nil {
		return fmt.Errorf("Error creating VPN Customer Gateway %s: %s", d.Get("name").(string), err)
	}

	d.SetId(v.Id)

	return resourceCosmicVPNCustomerGatewayRead(d, meta)
}

func resourceCosmicVPNCustomerGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	// Get the VPN Customer Gateway details
	v, count, err := client.VPN.GetVpnCustomerGatewayByID(d.Id())
	if err != nil {
		if count == 0 {
			log.Printf(
				"[DEBUG] VPN Customer Gateway %s does no longer exist", d.Get("name").(string))
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("name", v.Name)
	d.Set("cidr", v.Cidrlist)
	d.Set("esp_policy", v.Esppolicy)
	d.Set("gateway", v.Gateway)
	d.Set("ike_policy", v.Ikepolicy)
	d.Set("ipsec_psk", v.Ipsecpsk)
	d.Set("dpd", v.Dpd)
	d.Set("esp_lifetime", int(v.Esplifetime))
	d.Set("ike_lifetime", int(v.Ikelifetime))

	return nil
}

func resourceCosmicVPNCustomerGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	// Create a new parameter struct
	p := client.VPN.NewUpdateVpnCustomerGatewayParams(
		createCidrList(d.Get("cidr_list").(*schema.Set)),
		d.Get("esp_policy").(string),
		d.Get("gateway").(string),
		d.Id(),
		d.Get("ike_policy").(string),
		d.Get("ipsec_psk").(string),
	)

	p.SetName(d.Get("name").(string))

	if dpd, ok := d.GetOk("dpd"); ok {
		p.SetDpd(dpd.(bool))
	}

	if esplifetime, ok := d.GetOk("esp_lifetime"); ok {
		p.SetEsplifetime(int64(esplifetime.(int)))
	}

	if ikelifetime, ok := d.GetOk("ike_lifetime"); ok {
		p.SetIkelifetime(int64(ikelifetime.(int)))
	}

	// Update the VPN Customer Gateway
	_, err := client.VPN.UpdateVpnCustomerGateway(p)
	if err != nil {
		return fmt.Errorf("Error updating VPN Customer Gateway %s: %s", d.Get("name").(string), err)
	}

	return resourceCosmicVPNCustomerGatewayRead(d, meta)
}

func resourceCosmicVPNCustomerGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	// Create a new parameter struct
	p := client.VPN.NewDeleteVpnCustomerGatewayParams(d.Id())

	// Delete the VPN Customer Gateway
	_, err := client.VPN.DeleteVpnCustomerGateway(p)
	if err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", d.Id())) {
			return nil
		}

		return fmt.Errorf("Error deleting VPN Customer Gateway %s: %s", d.Get("name").(string), err)
	}

	return nil
}
