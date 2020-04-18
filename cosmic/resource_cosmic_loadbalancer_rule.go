package cosmic

import (
	"fmt"
	"log"
	"strings"

	"github.com/MissionCriticalCloud/go-cosmic/v6/cosmic"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCosmicLoadBalancerRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceCosmicLoadBalancerRuleCreate,
		Read:   resourceCosmicLoadBalancerRuleRead,
		Update: resourceCosmicLoadBalancerRuleUpdate,
		Delete: resourceCosmicLoadBalancerRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ip_address_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"network_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"algorithm": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"private_port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"public_port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"protocol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					switch v {
					case "tcp", "tcp-proxy":
					default:
						errs = append(errs, fmt.Errorf("%q must be either 'tcp' or 'tcp-proxy', got: %q", key, v))
					}

					return
				},
			},

			"member_ids": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceCosmicLoadBalancerRuleCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cosmic.CosmicClient)

	d.Partial(true)

	// Create a new parameter struct
	p := cs.LoadBalancer.NewCreateLoadBalancerRuleParams(
		d.Get("algorithm").(string),
		d.Get("name").(string),
		d.Get("private_port").(int),
		d.Get("public_port").(int),
	)

	// Don't autocreate a firewall rule, use a resource if needed
	p.SetOpenfirewall(false)

	// Set the description
	if description, ok := d.GetOk("description"); ok {
		p.SetDescription(description.(string))
	} else {
		p.SetDescription(d.Get("name").(string))
	}

	if networkid, ok := d.GetOk("network_id"); ok {
		// Set the network id
		p.SetNetworkid(networkid.(string))
	}

	// Set the protocol
	if protocol, ok := d.GetOk("protocol"); ok {
		p.SetProtocol(protocol.(string))
	}

	// Set the ipaddress id
	p.SetPublicipid(d.Get("ip_address_id").(string))

	// Create the load balancer rule
	r, err := cs.LoadBalancer.CreateLoadBalancerRule(p)
	if err != nil {
		return err
	}

	// Set the load balancer rule ID and set partials
	d.SetId(r.Id)
	d.SetPartial("name")
	d.SetPartial("description")
	d.SetPartial("ip_address_id")
	d.SetPartial("network_id")
	d.SetPartial("algorithm")
	d.SetPartial("private_port")
	d.SetPartial("public_port")
	d.SetPartial("protocol")

	// Create a new parameter struct
	ap := cs.LoadBalancer.NewAssignToLoadBalancerRuleParams(r.Id)

	var mbs []string
	for _, id := range d.Get("member_ids").([]interface{}) {
		mbs = append(mbs, id.(string))
	}

	ap.SetVirtualmachineids(mbs)

	_, err = cs.LoadBalancer.AssignToLoadBalancerRule(ap)
	if err != nil {
		return err
	}

	d.SetPartial("member_ids")
	d.Partial(false)

	return resourceCosmicLoadBalancerRuleRead(d, meta)
}

func resourceCosmicLoadBalancerRuleRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cosmic.CosmicClient)

	// Get the load balancer details
	lb, count, err := cs.LoadBalancer.GetLoadBalancerRuleByID(d.Id())
	if err != nil {
		if count == 0 {
			log.Printf("[DEBUG] Load balancer rule %s does no longer exist", d.Get("name").(string))
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("algorithm", lb.Algorithm)
	d.Set("public_port", lb.Publicport)
	d.Set("private_port", lb.Privateport)
	d.Set("ip_address_id", lb.Publicipid)

	// Only set network if user specified it to avoid spurious diffs
	if _, ok := d.GetOk("network_id"); ok {
		d.Set("network_id", lb.Networkid)
	}

	return nil
}

func resourceCosmicLoadBalancerRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cosmic.CosmicClient)

	if d.HasChange("name") || d.HasChange("description") || d.HasChange("algorithm") {
		name := d.Get("name").(string)

		// Create new parameter struct
		p := cs.LoadBalancer.NewUpdateLoadBalancerRuleParams(d.Id())

		if d.HasChange("name") {
			log.Printf("[DEBUG] Name has changed for load balancer rule %s, starting update", name)

			p.SetName(name)
		}

		if d.HasChange("description") {
			log.Printf(
				"[DEBUG] Description has changed for load balancer rule %s, starting update", name)

			p.SetDescription(d.Get("description").(string))
		}

		if d.HasChange("algorithm") {
			algorithm := d.Get("algorithm").(string)

			log.Printf(
				"[DEBUG] Algorithm has changed to %s for load balancer rule %s, starting update",
				algorithm,
				name,
			)

			// Set the new Algorithm
			p.SetAlgorithm(algorithm)
		}

		_, err := cs.LoadBalancer.UpdateLoadBalancerRule(p)
		if err != nil {
			return fmt.Errorf(
				"Error updating load balancer rule %s", name)
		}
	}
	return resourceCosmicLoadBalancerRuleRead(d, meta)
}

func resourceCosmicLoadBalancerRuleDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cosmic.CosmicClient)

	// Create a new parameter struct
	p := cs.LoadBalancer.NewDeleteLoadBalancerRuleParams(d.Id())

	log.Printf("[INFO] Deleting load balancer rule: %s", d.Get("name").(string))
	if _, err := cs.LoadBalancer.DeleteLoadBalancerRule(p); err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if !strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", d.Id())) {
			return err
		}
	}

	return nil
}
