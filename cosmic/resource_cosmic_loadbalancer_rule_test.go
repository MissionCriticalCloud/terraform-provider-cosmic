package cosmic

import (
	"fmt"
	"strings"
	"testing"

	"github.com/MissionCriticalCloud/go-cosmic/cosmic"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCosmicLoadBalancerRule_basic(t *testing.T) {
	if COSMIC_SERVICE_OFFERING_1 == "" {
		t.Skip("This test requires an existing service offering (set it by exporting COSMIC_SERVICE_OFFERING_1)")
	}

	if COSMIC_TEMPLATE == "" {
		t.Skip("This test requires an existing instance template (set it by exporting COSMIC_TEMPLATE)")
	}

	if COSMIC_DEFAULT_ALLOW_ACL_ID == "" {
		t.Skip("This test requires a \"default_allow\" ACL ID (set it by exporting COSMIC_DEFAULT_ALLOW_ACL_ID)")
	}

	if COSMIC_VPC_ID == "" {
		t.Skip("This test requires an existing VPC ID (set it by exporting COSMIC_VPC_ID)")
	}

	if COSMIC_VPC_NETWORK_OFFERING == "" {
		t.Skip("This test requires an existing VPC network offering (set it by exporting COSMIC_VPC_NETWORK_OFFERING)")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicLoadBalancerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicLoadBalancerRule_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicLoadBalancerRuleExist("cosmic_loadbalancer_rule.foo", nil),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "name", "terraform-lb"),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "algorithm", "roundrobin"),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "public_port", "80"),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "private_port", "80"),
				),
			},
		},
	})
}

func TestAccCosmicLoadBalancerRule_update(t *testing.T) {
	if COSMIC_SERVICE_OFFERING_1 == "" {
		t.Skip("This test requires an existing service offering (set it by exporting COSMIC_SERVICE_OFFERING_1)")
	}

	if COSMIC_TEMPLATE == "" {
		t.Skip("This test requires an existing instance template (set it by exporting COSMIC_TEMPLATE)")
	}

	if COSMIC_DEFAULT_ALLOW_ACL_ID == "" {
		t.Skip("This test requires a \"default_allow\" ACL ID (set it by exporting COSMIC_DEFAULT_ALLOW_ACL_ID)")
	}

	if COSMIC_VPC_ID == "" {
		t.Skip("This test requires an existing VPC ID (set it by exporting COSMIC_VPC_ID)")
	}

	if COSMIC_VPC_NETWORK_OFFERING == "" {
		t.Skip("This test requires an existing VPC network offering (set it by exporting COSMIC_VPC_NETWORK_OFFERING)")
	}

	var id string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicLoadBalancerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicLoadBalancerRule_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicLoadBalancerRuleExist("cosmic_loadbalancer_rule.foo", &id),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "name", "terraform-lb"),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "algorithm", "roundrobin"),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "public_port", "80"),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "private_port", "80"),
				),
			},

			{
				Config: testAccCosmicLoadBalancerRule_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicLoadBalancerRuleExist("cosmic_loadbalancer_rule.foo", &id),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "name", "terraform-lb-update"),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "algorithm", "leastconn"),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "public_port", "443"),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "private_port", "443"),
				),
			},
		},
	})
}

func testAccCheckCosmicLoadBalancerRuleExist(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No loadbalancer rule ID is set")
		}

		if id != nil {
			if *id != "" && *id != rs.Primary.ID {
				return fmt.Errorf("Resource ID has changed!")
			}

			*id = rs.Primary.ID
		}

		cs := testAccProvider.Meta().(*cosmic.CosmicClient)
		_, count, err := cs.LoadBalancer.GetLoadBalancerRuleByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if count == 0 {
			return fmt.Errorf("Loadbalancer rule %s not found", n)
		}

		return nil
	}
}

func testAccCheckCosmicLoadBalancerRuleDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cosmic.CosmicClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cosmic_loadbalancer_rule" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Loadbalancer rule ID is set")
		}

		for k, id := range rs.Primary.Attributes {
			if !strings.Contains(k, "uuid") {
				continue
			}

			_, _, err := cs.LoadBalancer.GetLoadBalancerRuleByID(id)
			if err == nil {
				return fmt.Errorf("Loadbalancer rule %s still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

var testAccCosmicLoadBalancerRule_basic = fmt.Sprintf(`
resource "cosmic_network" "foo" {
  name             = "terraform-network"
  cidr             = "10.0.10.0/24"
  gateway          = "10.0.10.1"
  network_offering = "%s"
  vpc_id           = "%s"
  zone             = "%s"
}

resource "cosmic_ipaddress" "foo" {
  acl_id = "%s"
  vpc_id = "${cosmic_network.foo.vpc_id}"
}

resource "cosmic_instance" "foo1" {
  name             = "terraform-server1"
  display_name     = "terraform"
  service_offering = "%s"
  network_id       = "${cosmic_network.foo.id}"
  template         = "%s"
  zone             = "${cosmic_network.foo.zone}"
  expunge          = true
}

resource "cosmic_loadbalancer_rule" "foo" {
  name          = "terraform-lb"
  ip_address_id = "${cosmic_ipaddress.foo.id}"
  algorithm     = "roundrobin"
  network_id    = "${cosmic_network.foo.id}"
  public_port   = 80
  private_port  = 80
  member_ids    = ["${cosmic_instance.foo1.id}"]
}`,
	COSMIC_VPC_NETWORK_OFFERING,
	COSMIC_VPC_ID,
	COSMIC_ZONE,
	COSMIC_DEFAULT_ALLOW_ACL_ID,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_TEMPLATE)

var testAccCosmicLoadBalancerRule_update = fmt.Sprintf(`
resource "cosmic_network" "foo" {
  name             = "terraform-network"
  cidr             = "10.0.10.0/24"
  gateway          = "10.0.10.1"
  network_offering = "%s"
  vpc_id           = "%s"
  zone             = "%s"
}

resource "cosmic_ipaddress" "foo" {
  acl_id = "%s"
  vpc_id = "${cosmic_network.foo.vpc_id}"
}

resource "cosmic_instance" "foo1" {
  name             = "terraform-server1"
  display_name     = "terraform"
  service_offering = "%s"
  network_id       = "${cosmic_network.foo.id}"
  template         = "%s"
  zone             = "${cosmic_network.foo.zone}"
  expunge          = true
}

resource "cosmic_instance" "foo2" {
  name             = "terraform-server2"
  display_name     = "terraform"
  service_offering = "${cosmic_instance.foo1.service_offering}"
  network_id       = "${cosmic_network.foo.id}"
  template         = "${cosmic_instance.foo1.template}"
  zone             = "${cosmic_network.foo.zone}"
  expunge          = true
}

resource "cosmic_loadbalancer_rule" "foo" {
  name          = "terraform-lb-update"
  ip_address_id = "${cosmic_ipaddress.foo.id}"
  algorithm     = "leastconn"
  network_id    = "${cosmic_network.foo.id}"
  public_port   = 443
  private_port  = 443
  member_ids    = ["${cosmic_instance.foo1.id}", "${cosmic_instance.foo2.id}"]
}`,
	COSMIC_VPC_NETWORK_OFFERING,
	COSMIC_VPC_ID,
	COSMIC_ZONE,
	COSMIC_DEFAULT_ALLOW_ACL_ID,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_TEMPLATE)
