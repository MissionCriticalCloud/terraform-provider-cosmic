package cosmic

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCosmicNetworkACL() *schema.Resource {
	return &schema.Resource{
		Create: resourceCosmicNetworkACLCreate,
		Read:   resourceCosmicNetworkACLRead,
		Update: resourceCosmicNetworkACLUpdate,
		Delete: resourceCosmicNetworkACLDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCosmicNetworkACLCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	name := d.Get("name").(string)

	// Create a new parameter struct
	p := client.NetworkACL.NewCreateNetworkACLListParams(name, d.Get("vpc_id").(string))

	// Set the description
	if description, ok := d.GetOk("description"); ok {
		p.SetDescription(description.(string))
	} else {
		p.SetDescription(name)
	}

	// Create the new network ACL list
	r, err := client.NetworkACL.CreateNetworkACLList(p)
	if err != nil {
		return fmt.Errorf("Error creating network ACL list %s: %s", name, err)
	}

	d.SetId(r.Id)

	return resourceCosmicNetworkACLRead(d, meta)
}

func resourceCosmicNetworkACLRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	// Get the network ACL list details
	f, count, err := client.NetworkACL.GetNetworkACLListByID(d.Id())
	if err != nil {
		if count == 0 {
			log.Printf(
				"[DEBUG] Network ACL list %s does no longer exist", d.Get("name").(string))
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("name", f.Name)
	d.Set("description", f.Description)
	d.Set("vpc_id", f.Vpcid)

	return nil
}

func resourceCosmicNetworkACLUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)
	name := d.Get("name").(string)

	// Create a new parameter struct
	p := client.NetworkACL.NewUpdateNetworkACLListParams(d.Id())

	// Check if the name or description is changed
	if d.HasChange("name") || d.HasChange("description") {
		p.SetName(name)

		// Compute/set the display text
		description := d.Get("description").(string)
		if description == "" {
			description = name
		}
		p.SetDescription(description)
	}

	// Update the network ACL
	_, err := client.NetworkACL.UpdateNetworkACLList(p)
	if err != nil {
		return fmt.Errorf(
			"Error updating network ACL %s: %s", name, err)
	}

	return resourceCosmicNetworkACLRead(d, meta)
}

func resourceCosmicNetworkACLDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	// Create a new parameter struct
	p := client.NetworkACL.NewDeleteNetworkACLListParams(d.Id())

	// Delete the network ACL list
	_, err := Retry(3, func() (interface{}, error) {
		return client.NetworkACL.DeleteNetworkACLList(p)
	})
	if err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", d.Id())) {
			return nil
		}

		return fmt.Errorf("Error deleting network ACL list %s: %s", d.Get("name").(string), err)
	}

	return nil
}
