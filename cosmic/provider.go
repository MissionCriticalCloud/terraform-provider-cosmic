package cosmic

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("COSMIC_API_URL", nil),
			},

			"api_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("COSMIC_API_KEY", nil),
			},

			"secret_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("COSMIC_SECRET_KEY", nil),
			},

			"http_get_only": &schema.Schema{
				Type:        schema.TypeBool,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("COSMIC_HTTP_GET_ONLY", false),
			},

			"timeout": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("COSMIC_TIMEOUT", 900),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"cosmic_affinity_group":       resourceCosmicAffinityGroup(),
			"cosmic_disk":                 resourceCosmicDisk(),
			"cosmic_egress_firewall":      resourceCosmicEgressFirewall(),
			"cosmic_firewall":             resourceCosmicFirewall(),
			"cosmic_instance":             resourceCosmicInstance(),
			"cosmic_ipaddress":            resourceCosmicIPAddress(),
			"cosmic_loadbalancer_rule":    resourceCosmicLoadBalancerRule(),
			"cosmic_network":              resourceCosmicNetwork(),
			"cosmic_network_acl":          resourceCosmicNetworkACL(),
			"cosmic_network_acl_rule":     resourceCosmicNetworkACLRule(),
			"cosmic_nic":                  resourceCosmicNIC(),
			"cosmic_port_forward":         resourceCosmicPortForward(),
			"cosmic_private_gateway":      resourceCosmicPrivateGateway(),
			"cosmic_secondary_ipaddress":  resourceCosmicSecondaryIPAddress(),
			"cosmic_security_group":       resourceCosmicSecurityGroup(),
			"cosmic_security_group_rule":  resourceCosmicSecurityGroupRule(),
			"cosmic_ssh_keypair":          resourceCosmicSSHKeyPair(),
			"cosmic_static_nat":           resourceCosmicStaticNAT(),
			"cosmic_static_route":         resourceCosmicStaticRoute(),
			"cosmic_template":             resourceCosmicTemplate(),
			"cosmic_vpc":                  resourceCosmicVPC(),
			"cosmic_vpn_connection":       resourceCosmicVPNConnection(),
			"cosmic_vpn_customer_gateway": resourceCosmicVPNCustomerGateway(),
			"cosmic_vpn_gateway":          resourceCosmicVPNGateway(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		APIURL:      d.Get("api_url").(string),
		APIKey:      d.Get("api_key").(string),
		SecretKey:   d.Get("secret_key").(string),
		HTTPGETOnly: d.Get("http_get_only").(bool),
		Timeout:     int64(d.Get("timeout").(int)),
	}

	return config.NewClient()
}
