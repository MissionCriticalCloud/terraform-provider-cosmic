package cosmic

import (
	"fmt"
	"testing"

	"github.com/MissionCriticalCloud/go-cosmic/v6/cosmic"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCosmicVPNCustomerGateway_basic(t *testing.T) {
	if COSMIC_VPC_OFFERING == "" {
		t.Skip("This test requires an existing VPC offering (set it by exporting COSMIC_VPC_OFFERING)")
	}

	var vpnCustomerGateway cosmic.VpnCustomerGateway

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicVPNCustomerGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicVPNCustomerGateway_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicVPNCustomerGatewayExists(
						"cosmic_vpn_customer_gateway.foo", &vpnCustomerGateway),
					testAccCheckCosmicVPNCustomerGatewayAttributes(&vpnCustomerGateway),
					resource.TestCheckResourceAttr(
						"cosmic_vpn_customer_gateway.foo", "name", "terraform-foo"),
					resource.TestCheckResourceAttr(
						"cosmic_vpn_customer_gateway.bar", "name", "terraform-bar"),
					resource.TestCheckResourceAttr(
						"cosmic_vpn_customer_gateway.bar", "esp_policy", "aes256-sha1"),
					resource.TestCheckResourceAttr(
						"cosmic_vpn_customer_gateway.foo", "ike_policy", "aes256-sha1;modp1024"),
				),
			},
		},
	})
}

func testAccCheckCosmicVPNCustomerGatewayExists(n string, vpnCustomerGateway *cosmic.VpnCustomerGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPN CustomerGateway ID is set")
		}

		cs := testAccProvider.Meta().(*cosmic.CosmicClient)
		v, _, err := cs.VPN.GetVpnCustomerGatewayByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if v.Id != rs.Primary.ID {
			return fmt.Errorf("VPN CustomerGateway not found")
		}

		*vpnCustomerGateway = *v

		return nil
	}
}

func testAccCheckCosmicVPNCustomerGatewayAttributes(vpnCustomerGateway *cosmic.VpnCustomerGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if vpnCustomerGateway.Esppolicy != "aes256-sha1" {
			return fmt.Errorf("Bad ESP policy: %s", vpnCustomerGateway.Esppolicy)
		}

		if vpnCustomerGateway.Ikepolicy != "aes256-sha1;modp1024" {
			return fmt.Errorf("Bad IKE policy: %s", vpnCustomerGateway.Ikepolicy)
		}

		if vpnCustomerGateway.Ipsecpsk != "terraform" {
			return fmt.Errorf("Bad IPSEC pre-shared key: %s", vpnCustomerGateway.Ipsecpsk)
		}

		return nil
	}
}

func testAccCheckCosmicVPNCustomerGatewayDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cosmic.CosmicClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cosmic_vpn_customer_gateway" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPN Customer Gateway ID is set")
		}

		_, _, err := cs.VPN.GetVpnCustomerGatewayByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("VPN Customer Gateway %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

var testAccCosmicVPNCustomerGateway_basic = fmt.Sprintf(`
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
  gateway    = "${cosmic_vpn_gateway.foo.public_ip}"
  esp_policy = "aes256-sha1"
  ike_policy = "aes256-sha1;modp1024"
  ipsec_psk  = "terraform"
}

resource "cosmic_vpn_customer_gateway" "bar" {
  name       = "terraform-bar"
  cidr_list  = ["${cosmic_vpc.bar.cidr}"]
  gateway    = "${cosmic_vpn_gateway.bar.public_ip}"
  esp_policy = "aes256-sha1"
  ike_policy = "aes256-sha1;modp1024"
  ipsec_psk  = "terraform"
}`,
	COSMIC_VPC_OFFERING,
	COSMIC_ZONE,
	COSMIC_VPC_OFFERING,
	COSMIC_ZONE,
)
