package cosmic

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCosmicAffinityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceCosmicAffinityGroupCreate,
		Read:   resourceCosmicAffinityGroupRead,
		Delete: resourceCosmicAffinityGroupDelete,
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

			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCosmicAffinityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	name := d.Get("name").(string)
	affinityGroupType := d.Get("type").(string)

	// Create a new parameter struct
	p := client.AffinityGroup.NewCreateAffinityGroupParams(name, affinityGroupType)

	// Set the description
	if description, ok := d.GetOk("description"); ok {
		p.SetDescription(description.(string))
	} else {
		p.SetDescription(name)
	}

	log.Printf("[DEBUG] Creating affinity group %s", name)
	r, err := client.AffinityGroup.CreateAffinityGroup(p)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Affinity group %s successfully created", name)
	d.SetId(r.Id)

	return resourceCosmicAffinityGroupRead(d, meta)
}

func resourceCosmicAffinityGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	log.Printf("[DEBUG] Rerieving affinity group %s", d.Get("name").(string))

	// Get the affinity group details
	ag, count, err := client.AffinityGroup.GetAffinityGroupByID(d.Id())
	if err != nil {
		if count == 0 {
			log.Printf("[DEBUG] Affinity group %s does not longer exist", d.Get("name").(string))
			d.SetId("")
			return nil
		}

		return err
	}

	// Update the config
	d.Set("name", ag.Name)
	d.Set("description", ag.Description)
	d.Set("type", ag.Type)

	return nil
}

func resourceCosmicAffinityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	// Create a new parameter struct
	p := client.AffinityGroup.NewDeleteAffinityGroupParams()
	p.SetId(d.Id())

	// Delete the affinity group
	_, err := client.AffinityGroup.DeleteAffinityGroup(p)
	if err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", d.Id())) {
			return nil
		}

		return fmt.Errorf("Error deleting affinity group: %s", err)
	}

	return nil
}
