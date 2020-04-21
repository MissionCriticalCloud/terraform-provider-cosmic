package cosmic

import (
	"fmt"
	"strings"
	"testing"

	"github.com/MissionCriticalCloud/go-cosmic/v6/cosmic"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCosmicNetwork_basic(t *testing.T) {
	if COSMIC_VPC_OFFERING == "" {
		t.Skip("This test requires an existing VPC offering (set it by exporting COSMIC_VPC_OFFERING)")
	}

	if COSMIC_VPC_NETWORK_OFFERING == "" {
		t.Skip("This test requires an existing VPC network offering (set it by exporting COSMIC_VPC_NETWORK_OFFERING)")
	}

	var network cosmic.Network

	resource.ParallelTest(t, resource.TestCase{
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
					resource.TestCheckResourceAttr(
						"cosmic_network.foo", "cidr", "10.0.10.0/24"),
					resource.TestCheckResourceAttr(
						"cosmic_network.foo", "gateway", "10.0.10.1"),
					testAccCheckNetworkTags(&network, "terraform-tag", "true"),
				),
			},
		},
	})
}

func TestAccCosmicNetwork_update(t *testing.T) {
	if COSMIC_VPC_OFFERING == "" {
		t.Skip("This test requires an existing VPC offering (set it by exporting COSMIC_VPC_OFFERING)")
	}

	if COSMIC_VPC_NETWORK_OFFERING == "" {
		t.Skip("This test requires an existing VPC network offering (set it by exporting COSMIC_VPC_NETWORK_OFFERING)")
	}

	var network cosmic.Network

	resource.ParallelTest(t, resource.TestCase{
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
					resource.TestCheckResourceAttr(
						"cosmic_network.foo", "cidr", "10.0.10.0/24"),
					resource.TestCheckResourceAttr(
						"cosmic_network.foo", "gateway", "10.0.10.1"),
				),
			},

			{
				Config: testAccCosmicNetwork_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicNetworkExists(
						"cosmic_network.foo", &network),
					testAccCheckCosmicNetworkBasicAttributes(&network),
					resource.TestCheckResourceAttr(
						"cosmic_network.foo", "cidr", "10.0.10.0/25"),
					// resource.TestCheckResourceAttr(
					// 	"cosmic_network.foo", "dns", "10.0.10.10"),
					resource.TestCheckResourceAttr(
						"cosmic_network.foo", "gateway", "10.0.10.100"),
				),
			},
		},
	})
}

func TestAccCosmicNetwork_dns(t *testing.T) {
	if COSMIC_VPC_OFFERING == "" {
		t.Skip("This test requires an existing VPC offering (set it by exporting COSMIC_VPC_OFFERING)")
	}

	if COSMIC_VPC_NETWORK_OFFERING == "" {
		t.Skip("This test requires an existing VPC network offering (set it by exporting COSMIC_VPC_NETWORK_OFFERING)")
	}

	var network cosmic.Network

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicNetwork_dns,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicNetworkExists(
						"cosmic_network.foo", &network),
					testAccCheckCosmicNetworkBasicAttributes(&network),
					resource.TestCheckResourceAttr(
						"cosmic_network.foo", "cidr", "10.0.10.0/24"),
					resource.TestCheckResourceAttr(
						"cosmic_network.foo", "dns.#", "1"),
					resource.TestCheckResourceAttr(
						"cosmic_network.foo", "dns.0", "10.10.10.10"),
					resource.TestCheckResourceAttr(
						"cosmic_network.foo", "gateway", "10.0.10.1"),
				),
			},
		},
	})
}

func TestAccCosmicNetwork_updateACL(t *testing.T) {
	if COSMIC_VPC_OFFERING == "" {
		t.Skip("This test requires an existing VPC offering (set it by exporting COSMIC_VPC_OFFERING)")
	}

	if COSMIC_VPC_NETWORK_OFFERING == "" {
		t.Skip("This test requires an existing VPC network offering (set it by exporting COSMIC_VPC_NETWORK_OFFERING)")
	}

	var network cosmic.Network

	resource.ParallelTest(t, resource.TestCase{
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

func testAccCheckCosmicNetworkExists(n string, network *cosmic.Network) resource.TestCheckFunc {
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
resource "cosmic_vpc" "foo" {
  name           = "terraform-vpc"
  display_text   = "terraform-vpc-text"
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

  tags = {
    terraform-tag = "true"
  }
}`,
	strings.ToLower(COSMIC_VPC_OFFERING),
	COSMIC_ZONE,
	COSMIC_VPC_NETWORK_OFFERING,
)

var testAccCosmicNetwork_update = fmt.Sprintf(`
resource "cosmic_vpc" "foo" {
  name           = "terraform-vpc"
  display_text   = "terraform-vpc-text"
  cidr           = "10.0.10.0/22"
  network_domain = "terraform-domain"
  vpc_offering   = "%s"
  zone           = "%s"
}

resource "cosmic_network" "foo" {
  name             = "terraform-network"
  cidr             = "10.0.10.0/25"
  gateway          = "10.0.10.100"
  dns              = ["10.10.10.10"]
  network_offering = "%s"
  vpc_id           = "${cosmic_vpc.foo.id}"
  zone             = "${cosmic_vpc.foo.zone}"

  tags = {
    terraform-tag = "true"
  }
}`,
	COSMIC_VPC_OFFERING,
	COSMIC_ZONE,
	COSMIC_VPC_NETWORK_OFFERING,
)

var testAccCosmicNetwork_dns = fmt.Sprintf(`
resource "cosmic_vpc" "foo" {
  name           = "terraform-vpc"
  display_text   = "terraform-vpc-text"
  cidr           = "10.0.10.0/22"
  network_domain = "terraform-domain"
  vpc_offering   = "%s"
  zone           = "%s"
}

resource "cosmic_network" "foo" {
  name             = "terraform-network"
  cidr             = "10.0.10.0/24"
  gateway          = "10.0.10.1"
  dns              = ["10.10.10.10"]
  network_offering = "%s"
  vpc_id           = "${cosmic_vpc.foo.id}"
  zone             = "${cosmic_vpc.foo.zone}"

  tags = {
    terraform-tag = "true"
  }
}`,
	COSMIC_VPC_OFFERING,
	COSMIC_ZONE,
	COSMIC_VPC_NETWORK_OFFERING,
)

var testAccCosmicNetwork_acl = fmt.Sprintf(`
resource "cosmic_vpc" "foo" {
  name           = "terraform-vpc"
  display_text   = "terraform-vpc-text"
  cidr           = "10.0.10.0/22"
  network_domain = "terraform-domain"
  vpc_offering   = "%s"
  zone           = "%s"
}

resource "cosmic_network_acl" "foo" {
  name   = "foo"
  vpc_id = "${cosmic_vpc.foo.id}"
}

resource "cosmic_network" "foo" {
  name             = "terraform-network"
  cidr             = "10.0.10.0/24"
  gateway          = "10.0.10.1"
  network_offering = "%s"
  vpc_id           = "${cosmic_vpc.foo.id}"
  acl_id           = "${cosmic_network_acl.foo.id}"
  zone             = "${cosmic_vpc.foo.zone}"
}`,
	COSMIC_VPC_OFFERING,
	COSMIC_ZONE,
	COSMIC_VPC_NETWORK_OFFERING,
)

var testAccCosmicNetwork_updateACL = fmt.Sprintf(`
resource "cosmic_vpc" "foo" {
  name           = "terraform-vpc"
  display_text   = "terraform-vpc-text"
  cidr           = "10.0.10.0/22"
  network_domain = "terraform-domain"
  vpc_offering   = "%s"
  zone           = "%s"
}

resource "cosmic_network_acl" "bar" {
  name   = "bar"
  vpc_id = "${cosmic_vpc.foo.id}"
}

resource "cosmic_network" "foo" {
  name             = "terraform-network"
  cidr             = "10.0.10.0/24"
  gateway          = "10.0.10.1"
  network_offering = "%s"
  vpc_id           = "${cosmic_vpc.foo.id}"
  acl_id           = "${cosmic_network_acl.bar.id}"
  zone             = "${cosmic_vpc.foo.zone}"
}`,
	COSMIC_VPC_OFFERING,
	COSMIC_ZONE,
	COSMIC_VPC_NETWORK_OFFERING,
)
