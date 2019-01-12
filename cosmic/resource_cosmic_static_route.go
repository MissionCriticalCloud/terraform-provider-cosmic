package cosmic

import (
	"fmt"
	"log"
	"strings"

	"github.com/MissionCriticalCloud/go-cosmic/cosmic"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCosmicStaticRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceCosmicStaticRouteCreate,
		Read:   resourceCosmicStaticRouteRead,
		Delete: resourceCosmicStaticRouteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"cidr": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"nexthop": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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

func resourceCosmicStaticRouteCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cosmic.CosmicClient)

	// Create a new parameter struct
	p := cs.VPC.NewCreateStaticRouteParams(
		d.Get("cidr").(string),
		d.Get("nexthop").(string),
		d.Get("vpc_id").(string),
	)

	// Create the new private gateway
	r, err := cs.VPC.CreateStaticRoute(p)
	if err != nil {
		return fmt.Errorf("Error creating static route for %s: %s", d.Get("cidr").(string), err)
	}

	d.SetId(r.Id)

	return resourceCosmicStaticRouteRead(d, meta)
}

func resourceCosmicStaticRouteRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cosmic.CosmicClient)

	// Get the virtual machine details
	route, count, err := cs.VPC.GetStaticRouteByID(d.Id())
	if err != nil {
		if count == 0 {
			log.Printf("[DEBUG] Static route %s does no longer exist", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("cidr", route.Cidr)
	d.Set("nexthop", route.Nexthop)
	d.Set("vpc_id", route.Vpcid)

	return nil
}

func resourceCosmicStaticRouteDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cosmic.CosmicClient)

	// Create a new parameter struct
	p := cs.VPC.NewDeleteStaticRouteParams(d.Id())

	// Delete the private gateway
	_, err := cs.VPC.DeleteStaticRoute(p)
	if err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", d.Id())) {
			return nil
		}

		return fmt.Errorf("Error deleting static route for %s: %s", d.Get("cidr").(string), err)
	}

	return nil
}
