package cosmic

import (
	"fmt"
	"testing"

	"github.com/MissionCriticalCloud/go-cosmic/v6/cosmic"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCosmicNIC_basic(t *testing.T) {
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

	var nic cosmic.Nic

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicNICDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicNIC_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicNICExists(
						"cosmic_instance.foo", "cosmic_nic.bar", &nic),
					testAccCheckCosmicNICAttributes(&nic),
				),
			},
		},
	})
}

func TestAccCosmicNIC_update(t *testing.T) {
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

	var nic cosmic.Nic

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicNICDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicNIC_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicNICExists(
						"cosmic_instance.foo", "cosmic_nic.bar", &nic),
					testAccCheckCosmicNICAttributes(&nic),
				),
			},

			{
				Config: testAccCosmicNIC_ipaddress,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicNICExists(
						"cosmic_instance.foo", "cosmic_nic.bar", &nic),
					testAccCheckCosmicNICAttributes(&nic),
					testAccCheckCosmicNICIPAddress(&nic),
					resource.TestCheckResourceAttr(
						"cosmic_nic.bar", "ip_address", "10.0.11.10"),
				),
			},
		},
	})
}

func testAccCheckCosmicNICExists(v, n string, nic *cosmic.Nic) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rsv, ok := s.RootModule().Resources[v]
		if !ok {
			return fmt.Errorf("Not found: %s", v)
		}

		if rsv.Primary.ID == "" {
			return fmt.Errorf("No instance ID is set")
		}

		rsn, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rsn.Primary.ID == "" {
			return fmt.Errorf("No NIC ID is set")
		}

		cs := testAccProvider.Meta().(*cosmic.CosmicClient)
		vm, _, err := cs.VirtualMachine.GetVirtualMachineByID(rsv.Primary.ID)

		if err != nil {
			return err
		}

		for _, n := range vm.Nic {
			if n.Id == rsn.Primary.ID {
				*nic = n
				return nil
			}
		}

		return fmt.Errorf("NIC not found")
	}
}

func testAccCheckCosmicNICAttributes(nic *cosmic.Nic) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if nic.Networkname != "terraform-network-bar" {
			return fmt.Errorf("Bad network name: %s", nic.Networkname)
		}

		return nil
	}
}

func testAccCheckCosmicNICIPAddress(nic *cosmic.Nic) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if nic.Networkname != "terraform-network-bar" {
			return fmt.Errorf("Bad network name: %s", nic.Networkname)
		}

		if nic.Ipaddress != "10.0.11.10" {
			return fmt.Errorf("Bad IP address: %s", nic.Ipaddress)
		}

		return nil
	}
}

func testAccCheckCosmicNICDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cosmic.CosmicClient)

	// Deleting the instance automatically deletes any additional NICs
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cosmic_instance" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No instance ID is set")
		}

		_, _, err := cs.VirtualMachine.GetVirtualMachineByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Virtual Machine %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

var testAccCosmicNIC_basic = fmt.Sprintf(`
resource "cosmic_network" "foo" {
  name             = "terraform-network-foo"
  cidr             = "10.0.10.0/24"
  gateway          = "10.0.10.1"
  network_offering = "%s"
  vpc_id           = "%s"
  zone             = "%s"
}

resource "cosmic_network" "bar" {
  name             = "terraform-network-bar"
  cidr             = "10.0.11.0/24"
  gateway          = "10.0.11.1"
  network_offering = "${cosmic_network.foo.network_offering}"
  vpc_id           = "${cosmic_network.foo.vpc_id}"
  zone             = "${cosmic_network.foo.zone}"
}

resource "cosmic_instance" "foo" {
  name             = "terraform-test"
  display_name     = "terraform"
  service_offering = "%s"
  network_id       = "${cosmic_network.foo.id}"
  template         = "%s"
  zone             = "${cosmic_network.foo.zone}"
  expunge          = true
}

resource "cosmic_nic" "bar" {
  network_id         = "${cosmic_network.bar.id}"
  virtual_machine_id = "${cosmic_instance.foo.id}"
}`,
	COSMIC_VPC_NETWORK_OFFERING,
	COSMIC_VPC_ID,
	COSMIC_ZONE,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_TEMPLATE,
)

var testAccCosmicNIC_ipaddress = fmt.Sprintf(`
resource "cosmic_network" "foo" {
  name             = "terraform-network-foo"
  cidr             = "10.0.10.0/24"
  gateway          = "10.0.10.1"
  network_offering = "%s"
  vpc_id           = "%s"
  zone             = "%s"
}

resource "cosmic_network" "bar" {
  name             = "terraform-network-bar"
  cidr             = "10.0.11.0/24"
  gateway          = "10.0.11.1"
  network_offering = "${cosmic_network.foo.network_offering}"
  vpc_id           = "${cosmic_network.foo.vpc_id}"
  zone             = "${cosmic_network.foo.zone}"
}

resource "cosmic_instance" "foo" {
  name             = "terraform-test"
  display_name     = "terraform"
  service_offering = "%s"
  network_id       = "${cosmic_network.foo.id}"
  template         = "%s"
  zone             = "${cosmic_network.foo.zone}"
  expunge          = true
}

resource "cosmic_nic" "bar" {
  network_id         = "${cosmic_network.bar.id}"
  ip_address         = "10.0.11.10"
  virtual_machine_id = "${cosmic_instance.foo.id}"
}`,
	COSMIC_VPC_NETWORK_OFFERING,
	COSMIC_VPC_ID,
	COSMIC_ZONE,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_TEMPLATE,
)
