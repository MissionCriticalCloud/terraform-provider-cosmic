package cosmic

import (
	"fmt"
	"testing"

	"github.com/MissionCriticalCloud/go-cosmic/v6/cosmic"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCosmicIPAddress_basic(t *testing.T) {
	if COSMIC_VPC_OFFERING == "" {
		t.Skip("This test requires an existing VPC offering (set it by exporting COSMIC_VPC_OFFERING)")
	}

	var ipaddr cosmic.PublicIpAddress

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicIPAddressDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicIPAddress_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicIPAddressExists(
						"cosmic_ipaddress.foo", &ipaddr),
					testAccCheckCosmicIPAddressAttributes(&ipaddr),
				),
			},
		},
	})
}

func testAccCheckCosmicIPAddressExists(n string, ipaddr *cosmic.PublicIpAddress) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No IP address ID is set")
		}

		cs := testAccProvider.Meta().(*cosmic.CosmicClient)
		pip, _, err := cs.PublicIPAddress.GetPublicIpAddressByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if pip.Id != rs.Primary.ID {
			return fmt.Errorf("IP address not found")
		}

		*ipaddr = *pip

		return nil
	}
}

func testAccCheckCosmicIPAddressAttributes(ipaddr *cosmic.PublicIpAddress) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "cosmic_vpc" {
				continue
			}

			if ipaddr.Vpcid != rs.Primary.ID {
				return fmt.Errorf("Bad network ID: %s", ipaddr.Associatednetworkid)
			}
		}

		return nil
	}
}

func testAccCheckCosmicIPAddressDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cosmic.CosmicClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cosmic_ipaddress" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No IP address ID is set")
		}

		ip, _, err := cs.PublicIPAddress.GetPublicIpAddressByID(rs.Primary.ID)
		if err == nil && ip.Associatednetworkid != "" {
			return fmt.Errorf("Public IP %s still associated", rs.Primary.ID)
		}
	}

	return nil
}

var testAccCosmicIPAddress_basic = fmt.Sprintf(`
resource "cosmic_vpc" "foo" {
  name           = "terraform-vpc"
  display_text   = "terraform-vpc"
  cidr           = "10.0.10.0/22"
  vpc_offering   = "%s"
  network_domain = "terraform-domain"
  zone           = "%s"
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
}`,
	COSMIC_VPC_OFFERING,
	COSMIC_ZONE,
)
