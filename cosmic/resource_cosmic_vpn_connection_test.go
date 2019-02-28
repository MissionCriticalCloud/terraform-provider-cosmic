package cosmic

import (
	"fmt"
	"testing"

	"github.com/MissionCriticalCloud/go-cosmic/cosmic"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCosmicVPNConnection_basic(t *testing.T) {
	if COSMIC_VPC_ID == "" {
		t.Skip("This test requires an existing VPC ID (set it by exporting COSMIC_VPC_ID)")
	}

	var vpnConnection cosmic.VpnConnection

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicVPNConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicVPNConnection_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicVPNConnectionExists(
						"cosmic_vpn_connection.foo-bar", &vpnConnection),
					testAccCheckCosmicVPNConnectionExists(
						"cosmic_vpn_connection.bar-foo", &vpnConnection),
				),
			},
		},
	})
}

func testAccCheckCosmicVPNConnectionExists(
	n string, vpnConnection *cosmic.VpnConnection) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPN Connection ID is set")
		}

		cs := testAccProvider.Meta().(*cosmic.CosmicClient)
		v, _, err := cs.VPN.GetVpnConnectionByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if v.Id != rs.Primary.ID {
			return fmt.Errorf("VPN Connection not found")
		}

		*vpnConnection = *v

		return nil
	}
}

func testAccCheckCosmicVPNConnectionDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cosmic.CosmicClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cosmic_vpn_connection" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPN Connection ID is set")
		}

		_, _, err := cs.VPN.GetVpnConnectionByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("VPN Connection %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

var testAccCosmicVPNConnection_basic = fmt.Sprintf(`
resource "cosmic_vpc" "foo" {
  name         = "terraform-vpc-foo"
  cidr         = "10.0.10.0/22"
  vpc_offering = "%s"
  zone         = "%s"
}

resource "cosmic_vpc" "bar" {
  name         = "terraform-vpc-bar"
  cidr         = "10.0.20.0/22"
  vpc_offering = "%s"
  zone         = "%s"
}

resource "cosmic_vpn_gateway" "foo" {
  vpc_id = "${cosmic_vpc.foo.id}"
}

resource "cosmic_vpn_gateway" "bar" {
  vpc_id = "${cosmic_vpc.bar.id}"
}

resource "cosmic_vpn_customer_gateway" "foo" {
  name       = "terraform-foo"
  cidr_list  = ["${cosmic_vpc.foo.cidr}"]
  esp_policy = "aes256-sha1"
  gateway    = "${cosmic_vpn_gateway.foo.public_ip}"
  ike_policy = "aes256-sha1;modp1024"
  ipsec_psk  = "terraform"
}

resource "cosmic_vpn_customer_gateway" "bar" {
  name       = "terraform-bar"
  cidr_list  = ["${cosmic_vpc.bar.cidr}"]
  esp_policy = "aes256-sha1"
  gateway    = "${cosmic_vpn_gateway.bar.public_ip}"
  ike_policy = "aes256-sha1;modp1024"
  ipsec_psk  = "terraform"
}

resource "cosmic_vpn_connection" "foo-bar" {
  customer_gateway_id = "${cosmic_vpn_customer_gateway.foo.id}"
  vpn_gateway_id      = "${cosmic_vpn_gateway.bar.id}"
}

resource "cosmic_vpn_connection" "bar-foo" {
  customer_gateway_id = "${cosmic_vpn_customer_gateway.bar.id}"
  vpn_gateway_id      = "${cosmic_vpn_gateway.foo.id}"
}`,
	COSMIC_VPC_OFFERING,
	COSMIC_ZONE,
	COSMIC_VPC_OFFERING,
	COSMIC_ZONE)
