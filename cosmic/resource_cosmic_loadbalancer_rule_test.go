package cosmic

import (
	"fmt"
	"strings"
	"testing"

	"github.com/MissionCriticalCloud/go-cosmic/v6/cosmic"
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

	if COSMIC_VPC_NETWORK_OFFERING == "" {
		t.Skip("This test requires an existing VPC network offering (set it by exporting COSMIC_VPC_NETWORK_OFFERING)")
	}

	if COSMIC_VPC_OFFERING == "" {
		t.Skip("This test requires an existing VPC offering (set it by exporting COSMIC_VPC_OFFERING)")
	}

	var rule cosmic.LoadBalancerRule

	createAttributes := &testAccCheckCosmicLoadBalancerRuleExpectedAttributes{
		Name:        "terraform-lb",
		Algorithm:   "roundrobin",
		PrivatePort: "8080",
		PublicPort:  "80",
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicLoadBalancerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicLoadBalancerRule_basic(createAttributes),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosmicLoadBalancerRuleExist("cosmic_loadbalancer_rule.foo", nil, &rule, false),
					testAccCheckCosmicLoadBalancerRuleAttributes(&rule, createAttributes),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "name", createAttributes.Name),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "algorithm", createAttributes.Algorithm),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "private_port", createAttributes.PrivatePort),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "public_port", createAttributes.PublicPort),
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

	if COSMIC_VPC_NETWORK_OFFERING == "" {
		t.Skip("This test requires an existing VPC network offering (set it by exporting COSMIC_VPC_NETWORK_OFFERING)")
	}

	if COSMIC_VPC_OFFERING == "" {
		t.Skip("This test requires an existing VPC offering (set it by exporting COSMIC_VPC_OFFERING)")
	}

	var id string
	var rule cosmic.LoadBalancerRule

	createAttributes := &testAccCheckCosmicLoadBalancerRuleExpectedAttributes{
		Name:        "terraform-lb",
		Algorithm:   "roundrobin",
		PrivatePort: "8080",
		PublicPort:  "80",
	}

	updateAttributes := &testAccCheckCosmicLoadBalancerRuleExpectedAttributes{
		Name:        "terraform-lb-update",
		Algorithm:   "leastconn",
		PrivatePort: "8080",
		PublicPort:  "80",
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicLoadBalancerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicLoadBalancerRule_basic(createAttributes),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosmicLoadBalancerRuleExist("cosmic_loadbalancer_rule.foo", &id, &rule, false),
					testAccCheckCosmicLoadBalancerRuleAttributes(&rule, createAttributes),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "name", createAttributes.Name),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "algorithm", createAttributes.Algorithm),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "private_port", createAttributes.PrivatePort),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "public_port", createAttributes.PublicPort),
				),
			},

			{
				Config: testAccCosmicLoadBalancerRule_update(updateAttributes),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosmicLoadBalancerRuleExist("cosmic_loadbalancer_rule.foo", &id, &rule, false),
					testAccCheckCosmicLoadBalancerRuleAttributes(&rule, updateAttributes),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "name", updateAttributes.Name),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "algorithm", updateAttributes.Algorithm),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "private_port", updateAttributes.PrivatePort),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "public_port", updateAttributes.PublicPort),
				),
			},
		},
	})
}

func testAccCheckCosmicLoadBalancerRuleExist(n string, id *string, rule *cosmic.LoadBalancerRule, shouldChange bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No loadbalancer rule ID is set")
		}

		if id != nil {
			if shouldChange {
				if *id != "" && *id == rs.Primary.ID {
					return fmt.Errorf("Resource ID has not changed!")
				}
			} else {
				if *id != "" && *id != rs.Primary.ID {
					return fmt.Errorf("Resource ID has changed!")
				}
			}

			*id = rs.Primary.ID
		}

		cs := testAccProvider.Meta().(*cosmic.CosmicClient)
		lbrule, count, err := cs.LoadBalancer.GetLoadBalancerRuleByID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if count == 0 {
			return fmt.Errorf("Loadbalancer rule %s not found", n)
		}

		*rule = *lbrule

		return nil
	}
}

type testAccCheckCosmicLoadBalancerRuleExpectedAttributes struct {
	Algorithm   string
	Name        string
	PrivatePort string
	PublicPort  string
}

func testAccCheckCosmicLoadBalancerRuleAttributes(rule *cosmic.LoadBalancerRule, want *testAccCheckCosmicLoadBalancerRuleExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if rule.Name != want.Name {
			return fmt.Errorf("Bad name: got %s; want %s", rule.Name, want.Name)
		}

		if rule.Algorithm != want.Algorithm {
			return fmt.Errorf("Bad algorithm: got %s; want %s", rule.Algorithm, want.Algorithm)
		}

		if rule.Privateport != want.PrivatePort {
			return fmt.Errorf("Bad private port: got %s; want %s", rule.Privateport, want.PrivatePort)
		}

		if rule.Publicport != want.PublicPort {
			return fmt.Errorf("Bad public port: got %s; want %s", rule.Publicport, want.PublicPort)
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

func testAccCosmicLoadBalancerRule_basic(attr *testAccCheckCosmicLoadBalancerRuleExpectedAttributes) string {
	return fmt.Sprintf(`
resource "cosmic_vpc" "foo" {
  name           = "terraform-vpc"
  display_text   = "terraform-vpc"
  cidr           = "10.0.10.0/22"
  network_domain = "terraform-domain"
  vpc_offering   = "%s"
  zone           = "%s"
}

resource "cosmic_network" "foo" {
  name             = "terraform-network"
  cidr             = "10.0.10.0/24"
  gateway          = "10.0.10.1"
  network_offering = "%s"
  vpc_id           = "${cosmic_vpc.foo.id}"
  zone             = "${cosmic_vpc.foo.zone}"
}

data "cosmic_network_acl" "default_allow" {
  filter {
    name  = "name"
    value = "default_allow"
  }
}

resource "cosmic_ipaddress" "foo" {
  acl_id = "${data.cosmic_network_acl.default_allow.id}"
  vpc_id = "${cosmic_vpc.foo.id}"
}

resource "cosmic_instance" "foo1" {
  name             = "terraform-server1"
  display_name     = "terraform-server1"
  service_offering = "%s"
  network_id       = "${cosmic_network.foo.id}"
  template         = "%s"
  zone             = "${cosmic_vpc.foo.zone}"
  expunge          = true
}

resource "cosmic_loadbalancer_rule" "foo" {
  name          = "%s"
  ip_address_id = "${cosmic_ipaddress.foo.id}"
  algorithm     = "%s"
  network_id    = "${cosmic_network.foo.id}"
  public_port   = "%s"
  private_port  = "%s"
  member_ids    = ["${cosmic_instance.foo1.id}"]
}`,
		COSMIC_VPC_OFFERING,
		COSMIC_ZONE,
		COSMIC_VPC_NETWORK_OFFERING,
		COSMIC_SERVICE_OFFERING_1,
		COSMIC_TEMPLATE,
		attr.Name,
		attr.Algorithm,
		attr.PublicPort,
		attr.PrivatePort,
	)
}

func testAccCosmicLoadBalancerRule_update(attr *testAccCheckCosmicLoadBalancerRuleExpectedAttributes) string {
	return fmt.Sprintf(`
resource "cosmic_vpc" "foo" {
  name           = "terraform-vpc"
  display_text   = "terraform-vpc"
  cidr           = "10.0.10.0/22"
  network_domain = "terraform-domain"
  vpc_offering   = "%s"
  zone           = "%s"
}

resource "cosmic_network" "foo" {
  name             = "terraform-network"
  cidr             = "10.0.10.0/24"
  gateway          = "10.0.10.1"
  network_offering = "%s"
  vpc_id           = "${cosmic_vpc.foo.id}"
  zone             = "${cosmic_vpc.foo.zone}"
}

data "cosmic_network_acl" "default_allow" {
  filter {
    name  = "name"
    value = "default_allow"
  }
}

resource "cosmic_ipaddress" "foo" {
  acl_id = "${data.cosmic_network_acl.default_allow.id}"
  vpc_id = "${cosmic_vpc.foo.id}"
}

resource "cosmic_instance" "foo1" {
  name             = "terraform-server1"
  display_name     = "terraform-server1"
  service_offering = "%s"
  network_id       = "${cosmic_network.foo.id}"
  template         = "%s"
  zone             = "${cosmic_vpc.foo.zone}"
  expunge          = true
}

resource "cosmic_instance" "foo2" {
  name             = "terraform-server2"
  display_name     = "terraform-server2"
  service_offering = "${cosmic_instance.foo1.service_offering}"
  network_id       = "${cosmic_network.foo.id}"
  template         = "${cosmic_instance.foo1.template}"
  zone             = "${cosmic_vpc.foo.zone}"
  expunge          = true
}

resource "cosmic_loadbalancer_rule" "foo" {
  name          = "%s"
  ip_address_id = "${cosmic_ipaddress.foo.id}"
  algorithm     = "%s"
  network_id    = "${cosmic_network.foo.id}"
  public_port   = "%s"
  private_port  = "%s"
  member_ids    = ["${cosmic_instance.foo1.id}", "${cosmic_instance.foo2.id}"]
}`,
		COSMIC_VPC_OFFERING,
		COSMIC_ZONE,
		COSMIC_VPC_NETWORK_OFFERING,
		COSMIC_SERVICE_OFFERING_1,
		COSMIC_TEMPLATE,
		attr.Name,
		attr.Algorithm,
		attr.PublicPort,
		attr.PrivatePort,
	)
}
