package cosmic

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/MissionCriticalCloud/go-cosmic/v6/cosmic"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCosmicPortForward() *schema.Resource {
	return &schema.Resource{
		Create: resourceCosmicPortForwardCreate,
		Read:   resourceCosmicPortForwardRead,
		Update: resourceCosmicPortForwardUpdate,
		Delete: resourceCosmicPortForwardDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCosmicPortForwardImporter,
		},

		Schema: map[string]*schema.Schema{
			"ip_address_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"managed": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"forward": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"private_port": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},

						"private_end_port": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},

						"public_port": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},

						"public_end_port": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},

						"virtual_machine_id": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"vm_guest_ip": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},

						"uuid": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceCosmicPortForwardCreate(d *schema.ResourceData, meta interface{}) error {
	// We need to set this upfront in order to be able to save a partial state
	d.SetId(d.Get("ip_address_id").(string))

	// Create all forwards that are configured
	if nrs := d.Get("forward").(*schema.Set); nrs.Len() > 0 {
		// Create an empty schema.Set to hold all forwards
		forwards := resourceCosmicPortForward().Schema["forward"].ZeroValue().(*schema.Set)

		err := createPortForwards(d, meta, forwards, nrs)

		// We need to update this first to preserve the correct state
		d.Set("forward", forwards)

		if err != nil {
			return err
		}
	}

	return resourceCosmicPortForwardRead(d, meta)
}

func createPortForwards(d *schema.ResourceData, meta interface{}, forwards *schema.Set, nrs *schema.Set) error {
	var errs *multierror.Error

	var wg sync.WaitGroup
	wg.Add(nrs.Len())

	sem := make(chan struct{}, 10)
	for _, forward := range nrs.List() {
		// Put in a tiny sleep here to avoid DoS'ing the API
		time.Sleep(500 * time.Millisecond)

		go func(forward map[string]interface{}) {
			defer wg.Done()
			sem <- struct{}{}

			// Create a single forward
			err := createPortForward(d, meta, forward)

			// If we have a UUID, we need to save the forward
			if forward["uuid"].(string) != "" {
				forwards.Add(forward)
			}

			if err != nil {
				errs = multierror.Append(errs, err)
			}

			<-sem
		}(forward.(map[string]interface{}))
	}

	wg.Wait()

	return errs.ErrorOrNil()
}

func createPortForward(d *schema.ResourceData, meta interface{}, forward map[string]interface{}) error {
	client := meta.(*CosmicClient)

	// Make sure all required parameters are there
	if err := verifyPortForwardParams(d, forward); err != nil {
		return err
	}

	vm, _, err := client.VirtualMachine.GetVirtualMachineByID(forward["virtual_machine_id"].(string))
	if err != nil {
		return err
	}

	// Create a new parameter struct
	p := client.Firewall.NewCreatePortForwardingRuleParams(d.Id(), forward["private_port"].(int),
		forward["protocol"].(string), forward["public_port"].(int), vm.Id)

	if privateEndPort, ok := forward["private_end_port"]; ok && privateEndPort.(int) != 0 {
		p.SetPrivateendport(privateEndPort.(int))
	}

	if publicEndPort, ok := forward["public_end_port"]; ok && publicEndPort.(int) != 0 {
		p.SetPublicendport(publicEndPort.(int))
	}

	if vmGuestIP, ok := forward["vm_guest_ip"]; ok && vmGuestIP.(string) != "" {
		p.SetVmguestip(vmGuestIP.(string))

		// Set the network ID based on the guest IP, needed when the public IP address
		// is not associated with any network yet
	NICS:
		for _, nic := range vm.Nic {
			if vmGuestIP.(string) == nic.Ipaddress {
				p.SetNetworkid(nic.Networkid)
				break NICS
			}
			for _, ip := range nic.Secondaryip {
				if vmGuestIP.(string) == ip.Ipaddress {
					p.SetNetworkid(nic.Networkid)
					break NICS
				}
			}
		}
	} else {
		// If no guest IP is configured, use the primary NIC
		p.SetNetworkid(vm.Nic[0].Networkid)
	}

	// Do not open the firewall automatically in any case
	p.SetOpenfirewall(false)

	r, err := client.Firewall.CreatePortForwardingRule(p)
	if err != nil {
		return err
	}

	forward["uuid"] = r.Id

	return nil
}

func resourceCosmicPortForwardRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	// First check if the IP address is still associated
	_, count, err := client.PublicIPAddress.GetPublicIpAddressByID(d.Id())
	if err != nil {
		if count == 0 {
			log.Printf(
				"[DEBUG] IP address with ID %s is no longer associated", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	// Get all the forwards from the running environment
	p := client.Firewall.NewListPortForwardingRulesParams()
	p.SetIpaddressid(d.Id())
	p.SetListall(true)

	l, err := client.Firewall.ListPortForwardingRules(p)
	if err != nil {
		return err
	}

	// Make a map of all the forwards so we can easily find a forward
	forwardMap := make(map[string]*cosmic.PortForwardingRule, l.Count)
	for _, f := range l.PortForwardingRules {
		forwardMap[f.Id] = f
	}

	// Create an empty schema.Set to hold all forwards
	forwards := resourceCosmicPortForward().Schema["forward"].ZeroValue().(*schema.Set)

	// Read all forwards that are configured
	if rs := d.Get("forward").(*schema.Set); rs.Len() > 0 {
		for _, forward := range rs.List() {
			forward := forward.(map[string]interface{})

			id, ok := forward["uuid"]
			if !ok || id.(string) == "" {
				continue
			}

			// Get the forward
			f, ok := forwardMap[id.(string)]
			if !ok {
				forward["uuid"] = ""
				continue
			}

			// Delete the known rule so only unknown rules remain in the ruleMap
			delete(forwardMap, id.(string))

			privPort, err := strconv.Atoi(f.Privateport)
			if err != nil {
				return err
			}

			privEndPort, err := strconv.Atoi(f.Privateendport)
			if err != nil {
				return err
			}

			pubPort, err := strconv.Atoi(f.Publicport)
			if err != nil {
				return err
			}

			pubEndPort, err := strconv.Atoi(f.Publicendport)
			if err != nil {
				return err
			}

			// Update the values
			forward["protocol"] = f.Protocol
			forward["private_port"] = privPort
			forward["public_port"] = pubPort
			forward["virtual_machine_id"] = f.Virtualmachineid
			// When specifying a single private or public port the same value is the end port, so
			// we only save the end port to state when it differs from the start port.
			if privEndPort != privPort {
				forward["private_end_port"] = privEndPort
			}
			if pubEndPort != pubPort {
				forward["public_end_port"] = pubEndPort
			}

			// This one is a bit tricky. We only want to update this optional value
			// if we've set one ourselves. If not this would become a computed value
			// and that would mess up the calculated hash of the set item.
			if forward["vm_guest_ip"].(string) != "" {
				forward["vm_guest_ip"] = f.Vmguestip
			}

			forwards.Add(forward)
		}
	}

	// If this is a managed resource, add all unknown forwards to dummy forwards
	managed := d.Get("managed").(bool)
	if managed && len(forwardMap) > 0 {
		for uuid := range forwardMap {
			// Make a dummy forward to hold the unknown UUID
			forward := map[string]interface{}{
				"protocol":           uuid,
				"private_port":       0,
				"private_end_port":   0,
				"public_port":        0,
				"public_end_port":    0,
				"virtual_machine_id": uuid,
				"uuid":               uuid,
			}

			// Add the dummy forward to the forwards set
			forwards.Add(forward)
		}
	}

	if forwards.Len() > 0 {
		d.Set("forward", forwards)
	} else if !managed {
		d.SetId("")
	}

	return nil
}

func resourceCosmicPortForwardUpdate(d *schema.ResourceData, meta interface{}) error {
	// Check if the forward set as a whole has changed
	if d.HasChange("forward") {
		o, n := d.GetChange("forward")
		ors := o.(*schema.Set).Difference(n.(*schema.Set))
		nrs := n.(*schema.Set).Difference(o.(*schema.Set))

		// We need to start with a rule set containing all the rules we
		// already have and want to keep. Any rules that are not deleted
		// correctly and any newly created rules, will be added to this
		// set to make sure we end up in a consistent state
		forwards := o.(*schema.Set).Intersection(n.(*schema.Set))

		// First loop through all the old forwards and delete them
		if ors.Len() > 0 {
			err := deletePortForwards(d, meta, forwards, ors)

			// We need to update this first to preserve the correct state
			d.Set("forward", forwards)

			if err != nil {
				return err
			}
		}

		// Then loop through all the new forwards and create them
		if nrs.Len() > 0 {
			err := createPortForwards(d, meta, forwards, nrs)

			// We need to update this first to preserve the correct state
			d.Set("forward", forwards)

			if err != nil {
				return err
			}
		}
	}

	return resourceCosmicPortForwardRead(d, meta)
}

func resourceCosmicPortForwardDelete(d *schema.ResourceData, meta interface{}) error {
	// Create an empty rule set to hold all rules that where
	// not deleted correctly
	forwards := resourceCosmicPortForward().Schema["forward"].ZeroValue().(*schema.Set)

	// Delete all forwards
	if ors := d.Get("forward").(*schema.Set); ors.Len() > 0 {
		err := deletePortForwards(d, meta, forwards, ors)

		// We need to update this first to preserve the correct state
		d.Set("forward", forwards)

		if err != nil {
			return err
		}
	}

	return nil
}

func deletePortForwards(d *schema.ResourceData, meta interface{}, forwards *schema.Set, ors *schema.Set) error {
	var errs *multierror.Error

	var wg sync.WaitGroup
	wg.Add(ors.Len())

	sem := make(chan struct{}, 10)
	for _, forward := range ors.List() {
		// Put a sleep here to avoid DoS'ing the API
		time.Sleep(500 * time.Millisecond)

		go func(forward map[string]interface{}) {
			defer wg.Done()
			sem <- struct{}{}

			// Delete a single forward
			err := deletePortForward(d, meta, forward)

			// If we have a UUID, we need to save the forward
			if forward["uuid"].(string) != "" {
				forwards.Add(forward)
			}

			if err != nil {
				errs = multierror.Append(errs, err)
			}

			<-sem
		}(forward.(map[string]interface{}))
	}

	wg.Wait()

	return errs.ErrorOrNil()
}

func deletePortForward(d *schema.ResourceData, meta interface{}, forward map[string]interface{}) error {
	client := meta.(*CosmicClient)

	// Create the parameter struct
	p := client.Firewall.NewDeletePortForwardingRuleParams(forward["uuid"].(string))

	// Delete the forward
	if _, err := client.Firewall.DeletePortForwardingRule(p); err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if !strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", forward["uuid"].(string))) {
			return err
		}
	}

	// Empty the UUID of this rule
	forward["uuid"] = ""

	return nil
}

func verifyPortForwardParams(d *schema.ResourceData, forward map[string]interface{}) error {
	protocol := forward["protocol"].(string)
	if protocol != "tcp" && protocol != "udp" {
		return fmt.Errorf(
			"%s is not a valid protocol. Valid options are 'tcp' and 'udp'", protocol)
	}
	return nil
}

func resourceCosmicPortForwardImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*CosmicClient)

	forwards := d.Get("forward").(*schema.Set)
	ipid := ""

	// As we can specify multiple forward {} blocks we should allow importing multiple
	// port forwards by iterating over a comma separated string
	for _, id := range strings.Split(d.Id(), ",") {
		l, count, err := client.Firewall.GetPortForwardingRuleByID(id)
		if err != nil {
			return nil, err
		}
		if count == 0 {
			return nil, fmt.Errorf("Port forwarding rule %s does not exist", id)
		}

		// Check each port forward rule we're importing is mapped to the same IP address
		if ipid != "" && l.Ipaddressid != ipid {
			return nil, fmt.Errorf("Port forwarding rule %s is not attached to expected IP address. Expected: %s, got: %s", id, ipid, l.Ipaddressid)
		}
		ipid = l.Ipaddressid

		privPort, err := strconv.Atoi(l.Privateport)
		if err != nil {
			return nil, err
		}

		privEndPort, err := strconv.Atoi(l.Privateendport)
		if err != nil {
			return nil, err
		}

		pubPort, err := strconv.Atoi(l.Publicport)
		if err != nil {
			return nil, err
		}

		pubEndPort, err := strconv.Atoi(l.Publicendport)
		if err != nil {
			return nil, err
		}

		forward := map[string]interface{}{
			"protocol":           l.Protocol,
			"private_port":       privPort,
			"public_port":        pubPort,
			"virtual_machine_id": l.Virtualmachineid,
			"uuid":               l.Id,
		}
		// Set end port value if it differs from start port.
		if privEndPort != privPort {
			forward["private_end_port"] = privEndPort
		}
		if pubEndPort != pubPort {
			forward["public_end_port"] = pubEndPort
		}

		// This is a tricky one: we're going to assume that vm_guest_ip isn't
		// set when using the instance's default IP, so we'll only set it during
		// import if the IP isn't the IP attached to the default NIC. If we're
		// wrong we'll end up needing to create all port fowards in the resource.
		vm, _, err := client.VirtualMachine.GetVirtualMachineByID(l.Virtualmachineid)
		if err != nil {
			return nil, err
		}
		// Loop through NICs and any configured secondary IPs
	NICs:
		for _, n := range vm.Nic {
			for _, sip := range n.Secondaryip {
				if sip.Ipaddress == l.Vmguestip {
					forward["vm_guest_ip"] = l.Vmguestip
					break NICs
				}
			}

			if n.Ipaddress == l.Vmguestip {
				if n.Isdefault != true {
					forward["vm_guest_ip"] = l.Vmguestip
					break NICs
				}
			}
		}

		log.Printf("[DEBUG] Importing forward: %v", forward)
		forwards.Add(forward)
	}

	d.SetId(ipid)
	d.Set("ip_address_id", ipid)
	d.Set("forward", forwards)
	d.Set("managed", false)
	return []*schema.ResourceData{d}, nil
}
