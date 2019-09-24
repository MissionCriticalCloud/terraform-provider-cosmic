package cosmic

import (
	"fmt"
	"strings"
	"testing"

	"github.com/MissionCriticalCloud/go-cosmic/v6/cosmic"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCosmicPortForward_basic(t *testing.T) {
	if COSMIC_SERVICE_OFFERING_1 == "" {
		t.Skip("This test requires an existing service offering (set it by exporting COSMIC_SERVICE_OFFERING_1)")
	}

	if COSMIC_TEMPLATE == "" {
		t.Skip("This test requires an existing instance template (set it by exporting COSMIC_TEMPLATE)")
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
	if COSMIC_SERVICE_OFFERING_1 == "" {
		t.Skip("This test requires an existing service offering (set it by exporting COSMIC_SERVICE_OFFERING_1)")
	}

	if COSMIC_TEMPLATE == "" {
		t.Skip("This test requires an existing instance template (set it by exporting COSMIC_TEMPLATE)")
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

func TestAccCosmicPortForward_endPort(t *testing.T) {
	if COSMIC_SERVICE_OFFERING_1 == "" {
		t.Skip("This test requires an existing service offering (set it by exporting COSMIC_SERVICE_OFFERING_1)")
	}

	if COSMIC_TEMPLATE == "" {
		t.Skip("This test requires an existing instance template (set it by exporting COSMIC_TEMPLATE)")
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
		CheckDestroy: testAccCheckCosmicPortForwardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicPortForward_endPort,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosmicPortForwardsExist("cosmic_port_forward.foo"),
					resource.TestCheckResourceAttr(
						"cosmic_port_forward.foo", "forward.#", "1"),
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
data "cosmic_network_acl" "default_allow" {
  filter {
    name  = "name"
    value = "default_allow"
  }
}

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
  acl_id = "${data.cosmic_network_acl.default_allow.id}"
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
)

var testAccCosmicPortForward_update = fmt.Sprintf(`
data "cosmic_network_acl" "default_allow" {
  filter {
    name  = "name"
    value = "default_allow"
  }
}

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
  acl_id = "${data.cosmic_network_acl.default_allow.id}"
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
)

var testAccCosmicPortForward_endPort = fmt.Sprintf(`
data "cosmic_network_acl" "default_allow" {
  filter {
    name  = "name"
    value = "default_allow"
  }
}

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
  acl_id = "${data.cosmic_network_acl.default_allow.id}"
  vpc_id = "${cosmic_network.foo.vpc_id}"
}

resource "cosmic_port_forward" "foo" {
  ip_address_id = "${cosmic_ipaddress.foo.id}"

  forward {
    protocol           = "tcp"
	private_port       = 443
	private_end_port   = 444
	public_port        = 443
	public_end_port    = 444
    virtual_machine_id = "${cosmic_instance.foo.id}"
  }
}`,
	COSMIC_VPC_NETWORK_OFFERING,
	COSMIC_VPC_ID,
	COSMIC_ZONE,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_TEMPLATE,
)
