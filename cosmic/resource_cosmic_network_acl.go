package cosmic

import (
	"fmt"
	"log"
	"strings"

	"github.com/MissionCriticalCloud/go-cosmic/v6/cosmic"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCosmicNetworkACL() *schema.Resource {
	return &schema.Resource{
		Create: resourceCosmicNetworkACLCreate,
		Read:   resourceCosmicNetworkACLRead,
		Delete: resourceCosmicNetworkACLDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCosmicNetworkACLCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cosmic.CosmicClient)

	name := d.Get("name").(string)

	// Create a new parameter struct
	p := cs.NetworkACL.NewCreateNetworkACLListParams(name, d.Get("vpc_id").(string))

	// Set the description
	if description, ok := d.GetOk("description"); ok {
		p.SetDescription(description.(string))
	} else {
		p.SetDescription(name)
	}

	// Create the new network ACL list
	r, err := cs.NetworkACL.CreateNetworkACLList(p)
	if err != nil {
		return fmt.Errorf("Error creating network ACL list %s: %s", name, err)
	}

	d.SetId(r.Id)

	return resourceCosmicNetworkACLRead(d, meta)
}

func resourceCosmicNetworkACLRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cosmic.CosmicClient)

	// Get the network ACL list details
	f, count, err := cs.NetworkACL.GetNetworkACLListByID(
		d.Id(),
		cosmic.WithProject(d.Get("project").(string)),
	)
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

func resourceCosmicNetworkACLDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cosmic.CosmicClient)

	// Create a new parameter struct
	p := cs.NetworkACL.NewDeleteNetworkACLListParams(d.Id())

	// Delete the network ACL list
	_, err := Retry(3, func() (interface{}, error) {
		return cs.NetworkACL.DeleteNetworkACLList(p)
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
