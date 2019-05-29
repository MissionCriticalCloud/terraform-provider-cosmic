package cosmic

import (
	"fmt"
	"testing"

	"github.com/MissionCriticalCloud/go-cosmic/v6/cosmic"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCosmicIPAddress_basic(t *testing.T) {
	if COSMIC_DEFAULT_ALLOW_ACL_ID == "" {
		t.Skip("This test requires a \"default_allow\" ACL ID (set it by exporting COSMIC_DEFAULT_ALLOW_ACL_ID)")
	}

	if COSMIC_VPC_ID == "" {
		t.Skip("This test requires an existing VPC ID (set it by exporting COSMIC_VPC_ID)")
	}

	var ipaddr cosmic.PublicIpAddress

	resource.Test(t, resource.TestCase{
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

func testAccCheckCosmicIPAddressExists(
	n string, ipaddr *cosmic.PublicIpAddress) resource.TestCheckFunc {
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

func testAccCheckCosmicIPAddressAttributes(
	ipaddr *cosmic.PublicIpAddress) resource.TestCheckFunc {
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
resource "cosmic_ipaddress" "foo" {
  acl_id = "%s"
  vpc_id = "%s"
}`,
	COSMIC_DEFAULT_ALLOW_ACL_ID,
	COSMIC_VPC_ID)
