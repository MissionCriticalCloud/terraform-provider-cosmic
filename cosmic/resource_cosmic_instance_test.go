package cosmic

import (
	"fmt"
	"testing"

	"github.com/MissionCriticalCloud/go-cosmic/v6/cosmic"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCosmicInstance_basic(t *testing.T) {
	if COSMIC_SERVICE_OFFERING_1 == "" {
		t.Skip("This test requires an existing service offering (set it by exporting COSMIC_SERVICE_OFFERING_1)")
	}

	if COSMIC_SERVICE_OFFERING_2 == "" {
		t.Skip("This test requires a second existing service offering (set it by exporting COSMIC_SERVICE_OFFERING_2)")
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

	var instance cosmic.VirtualMachine

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicInstance_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicInstanceExists(
						"cosmic_instance.foo", &instance),
					testAccCheckCosmicInstanceAttributes(&instance),
					resource.TestCheckResourceAttr(
						"cosmic_instance.foo", "user_data", "0cf3dcdc356ec8369494cb3991985ecd5296cdd5"),
				),
			},
		},
	})
}

func TestAccCosmicInstance_update(t *testing.T) {
	if COSMIC_SERVICE_OFFERING_1 == "" {
		t.Skip("This test requires an existing service offering (set it by exporting COSMIC_SERVICE_OFFERING_1)")
	}

	if COSMIC_SERVICE_OFFERING_2 == "" {
		t.Skip("This test requires a second existing service offering (set it by exporting COSMIC_SERVICE_OFFERING_2)")
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

	var instance cosmic.VirtualMachine

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicInstance_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicInstanceExists(
						"cosmic_instance.foo", &instance),
					testAccCheckCosmicInstanceAttributes(&instance),
					resource.TestCheckResourceAttr(
						"cosmic_instance.foo", "user_data", "0cf3dcdc356ec8369494cb3991985ecd5296cdd5"),
				),
			},

			{
				Config: testAccCosmicInstance_renameAndResize,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicInstanceExists(
						"cosmic_instance.foo", &instance),
					testAccCheckCosmicInstanceRenamedAndResized(&instance),
					resource.TestCheckResourceAttr(
						"cosmic_instance.foo", "name", "terraform-updated"),
					resource.TestCheckResourceAttr(
						"cosmic_instance.foo", "display_name", "terraform-updated"),
					resource.TestCheckResourceAttr(
						"cosmic_instance.foo", "service_offering", COSMIC_SERVICE_OFFERING_2),
				),
			},
		},
	})
}

func TestAccCosmicInstance_fixedIP(t *testing.T) {
	if COSMIC_SERVICE_OFFERING_1 == "" {
		t.Skip("This test requires an existing service offering (set it by exporting COSMIC_SERVICE_OFFERING_1)")
	}

	if COSMIC_SERVICE_OFFERING_2 == "" {
		t.Skip("This test requires a second existing service offering (set it by exporting COSMIC_SERVICE_OFFERING_2)")
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

	var instance cosmic.VirtualMachine

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicInstance_fixedIP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicInstanceExists(
						"cosmic_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"cosmic_instance.foo", "ip_address", "10.0.10.10"),
				),
			},
		},
	})
}

func TestAccCosmicInstance_keyPair(t *testing.T) {
	if COSMIC_SERVICE_OFFERING_1 == "" {
		t.Skip("This test requires an existing service offering (set it by exporting COSMIC_SERVICE_OFFERING_1)")
	}

	if COSMIC_SERVICE_OFFERING_2 == "" {
		t.Skip("This test requires a second existing service offering (set it by exporting COSMIC_SERVICE_OFFERING_2)")
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

	var instance cosmic.VirtualMachine

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicInstance_keyPair,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicInstanceExists(
						"cosmic_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"cosmic_instance.foo", "keypair", "terraform-test-keypair"),
				),
			},
		},
	})
}

func TestAccCosmicInstance_project(t *testing.T) {
	if COSMIC_PROJECT_NAME == "" {
		t.Skip("This test requires an existing project (set it by exporting COSMIC_PROJECT_NAME)")
	}

	if COSMIC_SERVICE_OFFERING_1 == "" {
		t.Skip("This test requires an existing service offering (set it by exporting COSMIC_SERVICE_OFFERING_1)")
	}

	if COSMIC_SERVICE_OFFERING_2 == "" {
		t.Skip("This test requires a second existing service offering (set it by exporting COSMIC_SERVICE_OFFERING_2)")
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

	var instance cosmic.VirtualMachine

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicInstance_project,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicInstanceExists(
						"cosmic_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"cosmic_instance.foo", "project", COSMIC_PROJECT_NAME),
				),
			},
		},
	})
}

func testAccCheckCosmicInstanceExists(
	n string, instance *cosmic.VirtualMachine) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No instance ID is set")
		}

		cs := testAccProvider.Meta().(*cosmic.CosmicClient)
		vm, _, err := cs.VirtualMachine.GetVirtualMachineByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if vm.Id != rs.Primary.ID {
			return fmt.Errorf("Instance not found")
		}

		*instance = *vm

		return nil
	}
}

func testAccCheckCosmicInstanceAttributes(
	instance *cosmic.VirtualMachine) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if instance.Name != "terraform-test" {
			return fmt.Errorf("Bad name: %s", instance.Name)
		}

		if instance.Displayname != "terraform-test" {
			return fmt.Errorf("Bad display name: %s", instance.Displayname)
		}

		if instance.Serviceofferingname != COSMIC_SERVICE_OFFERING_1 {
			return fmt.Errorf("Bad service offering: %s", instance.Serviceofferingname)
		}

		if instance.Templatename != COSMIC_TEMPLATE {
			return fmt.Errorf("Bad template: %s", instance.Templatename)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "cosmic_network" {
				continue
			}

			if instance.Nic[0].Networkid != rs.Primary.ID {
				return fmt.Errorf("Bad network ID: %s", instance.Nic[0].Networkid)
			}
		}

		return nil
	}
}

func testAccCheckCosmicInstanceRenamedAndResized(
	instance *cosmic.VirtualMachine) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if instance.Name != "terraform-updated" {
			return fmt.Errorf("Bad name: %s", instance.Name)
		}

		if instance.Displayname != "terraform-updated" {
			return fmt.Errorf("Bad display name: %s", instance.Displayname)
		}

		if instance.Serviceofferingname != COSMIC_SERVICE_OFFERING_2 {
			return fmt.Errorf("Bad service offering: %s", instance.Serviceofferingname)
		}

		return nil
	}
}

func testAccCheckCosmicInstanceDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cosmic.CosmicClient)

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

var testAccCosmicInstance_basic = fmt.Sprintf(`
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
  display_name     = "terraform-test"
  service_offering = "%s"
  network_id       = "${cosmic_network.foo.id}"
  template         = "%s"
  zone             = "${cosmic_network.foo.zone}"
  user_data        = "foobar\nfoo\nbar"
  expunge          = true
}`,
	COSMIC_VPC_NETWORK_OFFERING,
	COSMIC_VPC_ID,
	COSMIC_ZONE,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_TEMPLATE)

var testAccCosmicInstance_renameAndResize = fmt.Sprintf(`
resource "cosmic_network" "foo" {
  name             = "terraform-network"
  cidr             = "10.0.10.0/24"
  gateway          = "10.0.10.1"
  network_offering = "%s"
  vpc_id           = "%s"
  zone             = "%s"
}

resource "cosmic_instance" "foo" {
  name             = "terraform-updated"
  display_name     = "terraform-updated"
  service_offering = "%s"
  network_id       = "${cosmic_network.foo.id}"
  template         = "%s"
  zone             = "${cosmic_network.foo.zone}"
  user_data        = "foobar\nfoo\nbar"
  expunge          = true
}`,
	COSMIC_VPC_NETWORK_OFFERING,
	COSMIC_VPC_ID,
	COSMIC_ZONE,
	COSMIC_SERVICE_OFFERING_2,
	COSMIC_TEMPLATE)

var testAccCosmicInstance_fixedIP = fmt.Sprintf(`
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
  display_name     = "terraform-test"
  service_offering = "%s"
  network_id       = "${cosmic_network.foo.id}"
  ip_address       = "10.0.10.10"
  template         = "%s"
  zone             = "${cosmic_network.foo.zone}"
  expunge          = true
}`,
	COSMIC_VPC_NETWORK_OFFERING,
	COSMIC_VPC_ID,
	COSMIC_ZONE,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_TEMPLATE)

var testAccCosmicInstance_keyPair = fmt.Sprintf(`
resource "cosmic_ssh_keypair" "foo" {
  name = "terraform-test-keypair"
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
  display_name     = "terraform-test"
  service_offering = "%s"
  network_id       = "${cosmic_network.foo.id}"
  ip_address       = "10.0.10.10"
  template         = "%s"
  zone             = "${cosmic_network.foo.zone}"
  keypair          = "${cosmic_ssh_keypair.foo.name}"
  expunge          = true
}`,
	COSMIC_VPC_NETWORK_OFFERING,
	COSMIC_VPC_ID,
	COSMIC_ZONE,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_TEMPLATE)

var testAccCosmicInstance_project = fmt.Sprintf(`
resource "cosmic_vpc" "foo" {
  name           = "terraform-vpc"
  display_text   = "terraform-vpc-text"
  cidr           = "10.0.10.0/22"
  vpc_offering   = "%s"
  network_domain = "terraform-domain"
  project        = "%s"
  zone           = "%s"
}

resource "cosmic_network" "foo" {
  name             = "terraform-network"
  cidr             = "10.0.10.0/24"
  gateway          = "10.0.10.1"
  network_offering = "%s"
  project          = "${cosmic_vpc.foo.project}"
  vpc_id           = "${cosmic_vpc.foo.id}"
  zone             = "${cosmic_vpc.foo.zone}"
}

resource "cosmic_instance" "foo" {
  name             = "terraform-test"
  display_name     = "terraform-test"
  service_offering = "%s"
  network_id       = "${cosmic_network.foo.id}"
  template         = "%s"
  project          = "${cosmic_vpc.foo.project}"
  zone             = "${cosmic_vpc.foo.zone}"
  expunge          = true
}`,
	COSMIC_VPC_OFFERING,
	COSMIC_PROJECT_NAME,
	COSMIC_ZONE,
	COSMIC_VPC_NETWORK_OFFERING,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_TEMPLATE)
