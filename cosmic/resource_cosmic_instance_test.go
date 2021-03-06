package cosmic

import (
	"fmt"
	"strings"
	"testing"

	"github.com/MissionCriticalCloud/go-cosmic/v6/cosmic"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccCosmicInstance_basic(t *testing.T) {
	if COSMIC_SERVICE_OFFERING_1 == "" {
		t.Skip("This test requires an existing service offering (set it by exporting COSMIC_SERVICE_OFFERING_1)")
	}

	if COSMIC_TEMPLATE == "" {
		t.Skip("This test requires an existing instance template (set it by exporting COSMIC_TEMPLATE)")
	}

	if COSMIC_VPC_NETWORK_OFFERING == "" {
		t.Skip("This test requires an existing VPC network offering (set it by exporting COSMIC_VPC_NETWORK_OFFERING)")
	}

	if COSMIC_VPC_OFFERING == "" {
		t.Skip("This test requires an existing VPC offering (set it by exporting COSMIC_VPC_OFFERING)")
	}

	var instance cosmic.VirtualMachine

	resource.ParallelTest(t, resource.TestCase{
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
						"cosmic_instance.foo", "optimise_for", "Generic"),
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
		t.Skip("This test requires an existing second service offering (set it by exporting COSMIC_SERVICE_OFFERING_2)")
	}

	if COSMIC_TEMPLATE == "" {
		t.Skip("This test requires an existing instance template (set it by exporting COSMIC_TEMPLATE)")
	}

	if COSMIC_VPC_NETWORK_OFFERING == "" {
		t.Skip("This test requires an existing VPC network offering (set it by exporting COSMIC_VPC_NETWORK_OFFERING)")
	}

	if COSMIC_VPC_OFFERING == "" {
		t.Skip("This test requires an existing VPC offering (set it by exporting COSMIC_VPC_OFFERING)")
	}

	var instance cosmic.VirtualMachine

	resource.ParallelTest(t, resource.TestCase{
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

func TestAccCosmicInstance_diskController(t *testing.T) {
	if COSMIC_SERVICE_OFFERING_1 == "" {
		t.Skip("This test requires an existing service offering (set it by exporting COSMIC_SERVICE_OFFERING_1)")
	}

	if COSMIC_TEMPLATE == "" {
		t.Skip("This test requires an existing instance template (set it by exporting COSMIC_TEMPLATE)")
	}

	if COSMIC_VPC_NETWORK_OFFERING == "" {
		t.Skip("This test requires an existing VPC network offering (set it by exporting COSMIC_VPC_NETWORK_OFFERING)")
	}

	if COSMIC_VPC_OFFERING == "" {
		t.Skip("This test requires an existing VPC offering (set it by exporting COSMIC_VPC_OFFERING)")
	}

	var instance cosmic.VirtualMachine

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicInstance_diskController,
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

func TestAccCosmicInstance_fixedIP(t *testing.T) {
	if COSMIC_SERVICE_OFFERING_1 == "" {
		t.Skip("This test requires an existing service offering (set it by exporting COSMIC_SERVICE_OFFERING_1)")
	}

	if COSMIC_TEMPLATE == "" {
		t.Skip("This test requires an existing instance template (set it by exporting COSMIC_TEMPLATE)")
	}

	if COSMIC_VPC_NETWORK_OFFERING == "" {
		t.Skip("This test requires an existing VPC network offering (set it by exporting COSMIC_VPC_NETWORK_OFFERING)")
	}

	if COSMIC_VPC_OFFERING == "" {
		t.Skip("This test requires an existing VPC offering (set it by exporting COSMIC_VPC_OFFERING)")
	}

	var instance cosmic.VirtualMachine

	resource.ParallelTest(t, resource.TestCase{
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

	if COSMIC_TEMPLATE == "" {
		t.Skip("This test requires an existing instance template (set it by exporting COSMIC_TEMPLATE)")
	}

	if COSMIC_VPC_NETWORK_OFFERING == "" {
		t.Skip("This test requires an existing VPC network offering (set it by exporting COSMIC_VPC_NETWORK_OFFERING)")
	}

	if COSMIC_VPC_OFFERING == "" {
		t.Skip("This test requires an existing VPC offering (set it by exporting COSMIC_VPC_OFFERING)")
	}

	var instance cosmic.VirtualMachine

	keyPairName := fmt.Sprintf("terraform-test-keypair-%v", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicInstance_keyPair(keyPairName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicInstanceExists(
						"cosmic_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"cosmic_instance.foo", "keypair", keyPairName),
				),
			},
		},
	})
}

func TestAccCosmicInstance_import(t *testing.T) {
	if COSMIC_SERVICE_OFFERING_1 == "" {
		t.Skip("This test requires an existing service offering (set it by exporting COSMIC_SERVICE_OFFERING_1)")
	}

	if COSMIC_TEMPLATE == "" {
		t.Skip("This test requires an existing instance template (set it by exporting COSMIC_TEMPLATE)")
	}

	if COSMIC_VPC_NETWORK_OFFERING == "" {
		t.Skip("This test requires an existing VPC network offering (set it by exporting COSMIC_VPC_NETWORK_OFFERING)")
	}

	if COSMIC_VPC_OFFERING == "" {
		t.Skip("This test requires an existing VPC offering (set it by exporting COSMIC_VPC_OFFERING)")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicInstance_basic,
			},

			{
				ResourceName:            "cosmic_instance.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"expunge", "user_data"},
			},
		},
	})
}

func testAccCheckCosmicInstanceExists(n string, instance *cosmic.VirtualMachine) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No instance ID is set")
		}

		client := testAccProvider.Meta().(*CosmicClient)
		vm, _, err := client.VirtualMachine.GetVirtualMachineByID(rs.Primary.ID)

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

func testAccCheckCosmicInstanceAttributes(instance *cosmic.VirtualMachine) resource.TestCheckFunc {
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

func testAccCheckCosmicInstanceRenamedAndResized(instance *cosmic.VirtualMachine) resource.TestCheckFunc {
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
	client := testAccProvider.Meta().(*CosmicClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cosmic_instance" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No instance ID is set")
		}

		_, _, err := client.VirtualMachine.GetVirtualMachineByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Virtual Machine %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

var testAccCosmicInstance_basic = fmt.Sprintf(`
resource "cosmic_vpc" "foo" {
  name           = "terraform-vpc"
  display_text   = "terraform-vpc"
  cidr           = "10.0.10.0/22"
  network_domain = "terraform-domain"
  vpc_offering   = "%s"
}

resource "cosmic_network" "foo" {
  name             = "terraform-network"
  cidr             = "10.0.10.0/24"
  gateway          = "10.0.10.1"
  network_offering = "%s"
  vpc_id           = "${cosmic_vpc.foo.id}"
}

resource "cosmic_instance" "foo" {
  name             = "terraform-test"
  display_name     = "terraform-test"
  service_offering = "%s"
  network_id       = "${cosmic_network.foo.id}"
  template         = "%s"
  user_data        = "foobar\nfoo\nbar"
  expunge          = true
}`,
	COSMIC_VPC_OFFERING,
	COSMIC_VPC_NETWORK_OFFERING,
	strings.ToLower(COSMIC_SERVICE_OFFERING_1),
	COSMIC_TEMPLATE,
)

var testAccCosmicInstance_renameAndResize = fmt.Sprintf(`
resource "cosmic_vpc" "foo" {
  name           = "terraform-vpc"
  display_text   = "terraform-vpc"
  cidr           = "10.0.10.0/22"
  network_domain = "terraform-domain"
  vpc_offering   = "%s"
}

resource "cosmic_network" "foo" {
  name             = "terraform-network"
  cidr             = "10.0.10.0/24"
  gateway          = "10.0.10.1"
  network_offering = "%s"
  vpc_id           = "${cosmic_vpc.foo.id}"
}

resource "cosmic_instance" "foo" {
  name             = "terraform-updated"
  display_name     = "terraform-updated"
  service_offering = "%s"
  network_id       = "${cosmic_network.foo.id}"
  template         = "%s"
  user_data        = "foobar\nfoo\nbar"
  expunge          = true
}`,
	COSMIC_VPC_OFFERING,
	COSMIC_VPC_NETWORK_OFFERING,
	COSMIC_SERVICE_OFFERING_2,
	COSMIC_TEMPLATE,
)

var testAccCosmicInstance_diskController = fmt.Sprintf(`
resource "cosmic_vpc" "foo" {
  name           = "terraform-vpc"
  display_text   = "terraform-vpc"
  cidr           = "10.0.10.0/22"
  network_domain = "terraform-domain"
  vpc_offering   = "%s"
}

resource "cosmic_network" "foo" {
  name             = "terraform-network"
  cidr             = "10.0.10.0/24"
  gateway          = "10.0.10.1"
  network_offering = "%s"
  vpc_id           = "${cosmic_vpc.foo.id}"
}

resource "cosmic_instance" "foo" {
  name             = "terraform-test"
  display_name     = "terraform-test"
  service_offering = "%s"
  network_id       = "${cosmic_network.foo.id}"
  template         = "%s"
  disk_controller  = "SCSI"
  user_data        = "foobar\nfoo\nbar"
  expunge          = true
}`,
	COSMIC_VPC_OFFERING,
	COSMIC_VPC_NETWORK_OFFERING,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_TEMPLATE,
)

var testAccCosmicInstance_fixedIP = fmt.Sprintf(`
resource "cosmic_vpc" "foo" {
  name           = "terraform-vpc"
  display_text   = "terraform-vpc"
  cidr           = "10.0.10.0/22"
  network_domain = "terraform-domain"
  vpc_offering   = "%s"
}

resource "cosmic_network" "foo" {
  name             = "terraform-network"
  cidr             = "10.0.10.0/24"
  gateway          = "10.0.10.1"
  network_offering = "%s"
  vpc_id           = "${cosmic_vpc.foo.id}"
}

resource "cosmic_instance" "foo" {
  name             = "terraform-test"
  display_name     = "terraform-test"
  service_offering = "%s"
  network_id       = "${cosmic_network.foo.id}"
  ip_address       = "10.0.10.10"
  template         = "%s"
  expunge          = true
}`,
	COSMIC_VPC_OFFERING,
	COSMIC_VPC_NETWORK_OFFERING,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_TEMPLATE,
)

func testAccCosmicInstance_keyPair(keyPairName string) string {
	return fmt.Sprintf(`
resource "cosmic_vpc" "foo" {
  name           = "terraform-vpc"
  display_text   = "terraform-vpc"
  cidr           = "10.0.10.0/22"
  network_domain = "terraform-domain"
  vpc_offering   = "%s"
}

resource "cosmic_network" "foo" {
  name             = "terraform-network"
  cidr             = "10.0.10.0/24"
  gateway          = "10.0.10.1"
  network_offering = "%s"
  vpc_id           = "${cosmic_vpc.foo.id}"
}

resource "cosmic_ssh_keypair" "foo" {
  name = "%s"
}

resource "cosmic_instance" "foo" {
  name             = "terraform-test"
  display_name     = "terraform-test"
  service_offering = "%s"
  network_id       = "${cosmic_network.foo.id}"
  ip_address       = "10.0.10.10"
  template         = "%s"
  keypair          = "${cosmic_ssh_keypair.foo.name}"
  expunge          = true
}`,
		COSMIC_VPC_OFFERING,
		COSMIC_VPC_NETWORK_OFFERING,
		keyPairName,
		COSMIC_SERVICE_OFFERING_1,
		COSMIC_TEMPLATE,
	)
}
