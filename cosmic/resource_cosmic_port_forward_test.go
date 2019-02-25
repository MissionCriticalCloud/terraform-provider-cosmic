package cosmic

import (
	"fmt"
	"strings"
	"testing"

	"github.com/MissionCriticalCloud/go-cosmic/cosmic"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCosmicPortForward_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicPortForwardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicPortForward_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicPortForwardsExist("cosmic_port_forward.foo"),
					resource.TestCheckResourceAttr(
						"cosmic_port_forward.foo", "forward.#", "1"),
				),
			},
		},
	})
}

func TestAccCosmicPortForward_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicPortForwardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicPortForward_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicPortForwardsExist("cosmic_port_forward.foo"),
					resource.TestCheckResourceAttr(
						"cosmic_port_forward.foo", "forward.#", "1"),
				),
			},

			{
				Config: testAccCosmicPortForward_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicPortForwardsExist("cosmic_port_forward.foo"),
					resource.TestCheckResourceAttr(
						"cosmic_port_forward.foo", "forward.#", "2"),
				),
			},
		},
	})
}

func testAccCheckCosmicPortForwardsExist(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No port forward ID is set")
		}

		for k, id := range rs.Primary.Attributes {
			if !strings.Contains(k, "uuid") {
				continue
			}

			cs := testAccProvider.Meta().(*cosmic.CosmicClient)
			_, count, err := cs.Firewall.GetPortForwardingRuleByID(id)

			if err != nil {
				return err
			}

			if count == 0 {
				return fmt.Errorf("Port forward for %s not found", k)
			}
		}

		return nil
	}
}

func testAccCheckCosmicPortForwardDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cosmic.CosmicClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cosmic_port_forward" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No port forward ID is set")
		}

		for k, id := range rs.Primary.Attributes {
			if !strings.Contains(k, "uuid") {
				continue
			}

			_, _, err := cs.Firewall.GetPortForwardingRuleByID(id)
			if err == nil {
				return fmt.Errorf("Port forward %s still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

var testAccCosmicPortForward_basic = fmt.Sprintf(`
resource "cosmic_network" "foo" {
  name             = "terraform-network"
  cidr             = "10.0.10.0/24"
  gateway          = "10.0.10.1"
  network_offering = "%s"
  vpc_id           = "%s"
  zone             = "%s"
}

resource "cosmic_instance" "foo" {
  name             = "terraform-test"
  service_offering = "%s"
  network_id       = "${cosmic_network.foo.id}"
  template         = "%s"
  zone             = "${cosmic_network.foo.zone}"
  expunge          = true
}

resource "cosmic_ipaddress" "foo" {
  acl_id = "%s"
  vpc_id = "${cosmic_network.foo.vpc_id}"
}

resource "cosmic_port_forward" "foo" {
  ip_address_id = "${cosmic_ipaddress.foo.id}"

  forward {
    protocol           = "tcp"
    private_port       = 443
    public_port        = 8443
    virtual_machine_id = "${cosmic_instance.foo.id}"
  }
}`,
	COSMIC_VPC_NETWORK_OFFERING,
	COSMIC_VPC_ID,
	COSMIC_ZONE,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_TEMPLATE,
	COSMIC_DEFAULT_ALLOW_ACL_ID)

var testAccCosmicPortForward_update = fmt.Sprintf(`
resource "cosmic_network" "foo" {
  name             = "terraform-network"
  cidr             = "10.0.10.0/24"
  gateway          = "10.0.10.1"
  network_offering = "%s"
  vpc_id           = "%s"
  zone             = "%s"
}

resource "cosmic_instance" "foo" {
  name             = "terraform-test"
  service_offering = "%s"
  network_id       = "${cosmic_network.foo.id}"
  template         = "%s"
  zone             = "${cosmic_network.foo.zone}"
  expunge          = true
}

resource "cosmic_ipaddress" "foo" {
  acl_id = "%s"
  vpc_id = "${cosmic_network.foo.vpc_id}"
}

resource "cosmic_port_forward" "foo" {
  ip_address_id = "${cosmic_ipaddress.foo.id}"

  forward {
    protocol           = "tcp"
    private_port       = 443
    public_port        = 8443
    virtual_machine_id = "${cosmic_instance.foo.id}"
  }

  forward {
    protocol           = "tcp"
    private_port       = 80
    public_port        = 8080
    virtual_machine_id = "${cosmic_instance.foo.id}"
  }
}`,
	COSMIC_VPC_NETWORK_OFFERING,
	COSMIC_VPC_ID,
	COSMIC_ZONE,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_TEMPLATE,
	COSMIC_DEFAULT_ALLOW_ACL_ID)
