package cosmic

import (
	"fmt"
	"testing"

	"github.com/MissionCriticalCloud/go-cosmic/v6/cosmic"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCosmicPrivateGateway_basic(t *testing.T) {
	if COSMIC_VPC_OFFERING == "" {
		t.Skip("This test requires an existing VPC offering (set it by exporting COSMIC_VPC_OFFERING)")
	}

	var gateway cosmic.PrivateGateway

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicPrivateGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicPrivateGateway_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicPrivateGatewayExists(
						"cosmic_private_gateway.foo", &gateway),
					testAccCheckCosmicPrivateGatewayAttributes(&gateway),
				),
			},
		},
	})
}

func testAccCheckCosmicPrivateGatewayExists(n string, gateway *cosmic.PrivateGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Private Gateway ID is set")
		}

		cs := testAccProvider.Meta().(*cosmic.CosmicClient)
		pgw, _, err := cs.VPC.GetPrivateGatewayByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if pgw.Id != rs.Primary.ID {
			return fmt.Errorf("Private Gateway not found")
		}

		*gateway = *pgw

		return nil
	}
}

func testAccCheckCosmicPrivateGatewayAttributes(gateway *cosmic.PrivateGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if gateway.Ipaddress != "10.0.252.254" {
			return fmt.Errorf("Bad Gateway: %s", gateway.Ipaddress)
		}

		return nil
	}
}

func testAccCheckCosmicPrivateGatewayDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cosmic.CosmicClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cosmic_private_gateway" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No private gateway ID is set")
		}

		gateway, _, err := cs.VPC.GetPrivateGatewayByID(rs.Primary.ID)
		if err == nil && gateway.Id != "" {
			return fmt.Errorf("Private gateway %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

var testAccCosmicPrivateGateway_basic = fmt.Sprintf(`
resource "cosmic_vpc" "foo" {
  name           = "terraform-vpc"
  display_text   = "terraform-vpc"
  cidr           = "10.0.10.0/22"
  vpc_offering   = "%s"
  network_domain = "terraform-domain"
  zone           = "%s"
}

resource "cosmic_network" "foo" {
  name             = "terraform-network"
  cidr             = "10.0.252.0/24"
  network_offering = "DefaultPrivateGatewayNetworkOffering"
  zone             = "${cosmic_vpc.foo.zone}"
}

resource "cosmic_network_acl" "foo" {
  name   = "terraform-acl"
  vpc_id = "${cosmic_vpc.foo.id}"
}

resource "cosmic_private_gateway" "foo" {
  ip_address = "10.0.252.254"
  network_id = "${cosmic_network.foo.id}"
  acl_id     = "${cosmic_network_acl.foo.id}"
  vpc_id     = "${cosmic_network_acl.foo.vpc_id}"
}`,
	COSMIC_VPC_OFFERING,
	COSMIC_ZONE,
)
