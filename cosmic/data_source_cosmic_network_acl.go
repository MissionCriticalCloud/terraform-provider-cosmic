package cosmic

import (
	"fmt"
	"log"

	"github.com/MissionCriticalCloud/go-cosmic/v6/cosmic"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceCosmicNetworkACL() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCosmicNetworkACLRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			// Computed values
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCosmicNetworkACLRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*cosmic.CosmicClient)

	p := conn.NetworkACL.NewListNetworkACLListsParams()
	p.SetListall(true)

	cosmicACLLists, err := conn.NetworkACL.ListNetworkACLLists(p)
	if err != nil {
		return fmt.Errorf("Failed to list ACL Lists: %s", err)
	}

	filters := d.Get("filter")
	var ACLLists []*cosmic.NetworkACLList

	for _, l := range cosmicACLLists.NetworkACLLists {
		match, err := matchFilters(l, filters.(*schema.Set))
		if err != nil {
			return err
		}

		if match {
			ACLLists = append(ACLLists, l)
		}
	}

	if len(ACLLists) == 0 {
		return fmt.Errorf("No Network ACL List matched with the specified filter")
	}

	if len(ACLLists) > 1 {
		return fmt.Errorf("More than one Network ACL List found, use a more specific filter")
	}

	log.Printf("[DEBUG] Selected Network ACL List: %s\n", ACLLists[0].Id)

	d.SetId(ACLLists[0].Id)
	d.Set("name", ACLLists[0].Name)
	d.Set("description", ACLLists[0].Description)
	d.Set("vpc_id", ACLLists[0].Vpcid)

	return nil
}
