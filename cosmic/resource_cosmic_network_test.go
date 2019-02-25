package cosmic

import (
	"fmt"
	"testing"

	"github.com/MissionCriticalCloud/go-cosmic/cosmic"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCosmicNetwork_basic(t *testing.T) {
	var network cosmic.Network

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicNetwork_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicNetworkExists(
						"cosmic_network.foo", &network),
					testAccCheckCosmicNetworkBasicAttributes(&network),
					testAccCheckNetworkTags(&network, "terraform-tag", "true"),
				),
			},
		},
	})
}

func TestAccCosmicNetwork_updateACL(t *testing.T) {
	var network cosmic.Network

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicNetwork_acl,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicNetworkExists(
						"cosmic_network.foo", &network),
					testAccCheckCosmicNetworkBasicAttributes(&network),
				),
			},

			{
				Config: testAccCosmicNetwork_updateACL,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicNetworkExists(
						"cosmic_network.foo", &network),
					testAccCheckCosmicNetworkBasicAttributes(&network),
				),
			},
		},
	})
}

func testAccCheckCosmicNetworkExists(
	n string, network *cosmic.Network) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No network ID is set")
		}

		cs := testAccProvider.Meta().(*cosmic.CosmicClient)
		ntwrk, _, err := cs.Network.GetNetworkByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if ntwrk.Id != rs.Primary.ID {
			return fmt.Errorf("Network not found")
		}

		*network = *ntwrk

		return nil
	}
}

func testAccCheckCosmicNetworkBasicAttributes(network *cosmic.Network) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if network.Name != "terraform-network" {
			return fmt.Errorf("Bad name: %s", network.Name)
		}

		if network.Displaytext != "terraform-network" {
			return fmt.Errorf("Bad display name: %s", network.Displaytext)
		}

		if network.Cidr != "10.0.10.0/24" {
			return fmt.Errorf("Bad CIDR: %s", network.Cidr)
		}

		if network.Networkofferingname != COSMIC_VPC_NETWORK_OFFERING {
			return fmt.Errorf("Bad network offering: %s", network.Networkofferingname)
		}

		return nil
	}
}

func testAccCheckNetworkTags(n *cosmic.Network, key string, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		tags := make(map[string]string)
		for item := range n.Tags {
			tags[n.Tags[item].Key] = n.Tags[item].Value
		}
		return testAccCheckTags(tags, key, value)
	}
}

func testAccCheckCosmicNetworkDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cosmic.CosmicClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cosmic_network" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No network ID is set")
		}

		_, _, err := cs.Network.GetNetworkByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Network %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

var testAccCosmicNetwork_basic = fmt.Sprintf(`
resource "cosmic_network" "foo" {
  name             = "terraform-network"
  cidr             = "10.0.10.0/24"
  gateway          = "10.0.10.1"
  network_offering = "%s"
  vpc_id           = "%s"
  zone             = "%s"

  tags = {
    terraform-tag = "true"
  }
}`,
	COSMIC_VPC_NETWORK_OFFERING,
	COSMIC_VPC_ID,
	COSMIC_ZONE)

var testAccCosmicNetwork_acl = fmt.Sprintf(`
resource "cosmic_network_acl" "foo" {
  name   = "foo"
  vpc_id = "%s"
}

resource "cosmic_network" "foo" {
  name             = "terraform-network"
  cidr             = "10.0.10.0/24"
  gateway          = "10.0.10.1"
  network_offering = "%s"
  vpc_id           = "${cosmic_network_acl.foo.vpc_id}"
  acl_id           = "${cosmic_network_acl.foo.id}"
  zone             = "%s"
}`,
	COSMIC_VPC_ID,
	COSMIC_VPC_NETWORK_OFFERING,
	COSMIC_ZONE)

var testAccCosmicNetwork_updateACL = fmt.Sprintf(`
resource "cosmic_network_acl" "bar" {
  name   = "bar"
  vpc_id = "%s"
}

resource "cosmic_network" "foo" {
  name             = "terraform-network"
  cidr             = "10.0.10.0/24"
  gateway          = "10.0.10.1"
  network_offering = "%s"
  vpc_id           = "${cosmic_network_acl.bar.vpc_id}"
  acl_id           = "${cosmic_network_acl.bar.id}"
  zone             = "%s"
}`,
	COSMIC_VPC_ID,
	COSMIC_VPC_NETWORK_OFFERING,
	COSMIC_ZONE)
