package cosmic

import (
	"fmt"
	"testing"

	"github.com/MissionCriticalCloud/go-cosmic/v6/cosmic"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccCosmicSecondaryIPAddress_basic(t *testing.T) {
	var ip cosmic.AddIpToNicResponse

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicSecondaryIPAddressDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicSecondaryIPAddress_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicSecondaryIPAddressExists(
						"cosmic_secondary_ipaddress.foo", &ip),
				),
			},
		},
	})
}

func TestAccCosmicSecondaryIPAddress_fixedIP(t *testing.T) {
	var ip cosmic.AddIpToNicResponse

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicSecondaryIPAddressDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicSecondaryIPAddress_fixedIP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicSecondaryIPAddressExists(
						"cosmic_secondary_ipaddress.foo", &ip),
					testAccCheckCosmicSecondaryIPAddressAttributes(&ip),
					resource.TestCheckResourceAttr(
						"cosmic_secondary_ipaddress.foo", "ip_address", "10.0.10.10"),
				),
			},
		},
	})
}

func testAccCheckCosmicSecondaryIPAddressExists(n string, ip *cosmic.AddIpToNicResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No IP address ID is set")
		}

		client := testAccProvider.Meta().(*CosmicClient)

		virtualmachine, ok := rs.Primary.Attributes["virtual_machine_id"]
		if !ok {
			virtualmachine, ok = rs.Primary.Attributes["virtual_machine"]
		}

		// Retrieve the virtual_machine ID
		virtualmachineid, e := retrieveID(client, "virtual_machine", virtualmachine)
		if e != nil {
			return e.Error()
		}

		// Get the virtual machine details
		vm, count, err := client.VirtualMachine.GetVirtualMachineByID(virtualmachineid)
		if err != nil {
			if count == 0 {
				return fmt.Errorf("Instance not found")
			}
			return err
		}

		nicid, ok := rs.Primary.Attributes["nic_id"]
		if !ok {
			nicid, ok = rs.Primary.Attributes["nicid"]
		}
		if !ok {
			nicid = vm.Nic[0].Id
		}

		p := client.Nic.NewListNicsParams(virtualmachineid)
		p.SetNicid(nicid)

		l, err := client.Nic.ListNics(p)
		if err != nil {
			return err
		}

		if l.Count == 0 {
			return fmt.Errorf("NIC not found")
		}

		if l.Count > 1 {
			return fmt.Errorf("Found more then one possible result: %v", l.Nics)
		}

		for _, sip := range l.Nics[0].Secondaryip {
			if sip.Id == rs.Primary.ID {
				ip.Ipaddress = sip.Ipaddress
				ip.Nicid = l.Nics[0].Id
				return nil
			}
		}

		return fmt.Errorf("IP address not found")
	}
}

func testAccCheckCosmicSecondaryIPAddressAttributes(ip *cosmic.AddIpToNicResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if ip.Ipaddress != "10.0.10.10" {
			return fmt.Errorf("Bad IP address: %s", ip.Ipaddress)
		}
		return nil
	}
}

func testAccCheckCosmicSecondaryIPAddressDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CosmicClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cosmic_secondary_ipaddress" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No IP address ID is set")
		}

		virtualmachine, ok := rs.Primary.Attributes["virtual_machine_id"]
		if !ok {
			virtualmachine, ok = rs.Primary.Attributes["virtual_machine"]
		}

		// Retrieve the virtual_machine ID
		virtualmachineid, e := retrieveID(client, "virtual_machine", virtualmachine)
		if e != nil {
			return e.Error()
		}

		// Get the virtual machine details
		vm, count, err := client.VirtualMachine.GetVirtualMachineByID(virtualmachineid)
		if err != nil {
			if count == 0 {
				return nil
			}
			return err
		}

		nicid, ok := rs.Primary.Attributes["nic_id"]
		if !ok {
			nicid, ok = rs.Primary.Attributes["nicid"]
		}
		if !ok {
			nicid = vm.Nic[0].Id
		}

		p := client.Nic.NewListNicsParams(virtualmachineid)
		p.SetNicid(nicid)

		l, err := client.Nic.ListNics(p)
		if err != nil {
			return err
		}

		if l.Count == 0 {
			return fmt.Errorf("NIC not found")
		}

		if l.Count > 1 {
			return fmt.Errorf("Found more then one possible result: %v", l.Nics)
		}

		for _, sip := range l.Nics[0].Secondaryip {
			if sip.Id == rs.Primary.ID {
				return fmt.Errorf("IP address %s still exists", rs.Primary.ID)
			}
		}

		return nil
	}

	return nil
}

var testAccCosmicSecondaryIPAddress_basic = fmt.Sprintf(`
resource "cosmic_vpc" "foo" {
  name           = "terraform-vpc"
  display_text   = "terraform-vpc"
  cidr           = "10.0.10.0/22"
  network_domain = "terraform-domain"
  vpc_offering   = "%s"
  zone           = "%s"
}

resource "cosmic_network" "foo" {
  name             = "terraform-network"
  cidr             = "10.0.10.0/24"
  gateway          = "10.0.10.1"
  network_offering = "%s"
  vpc_id           = "${cosmic_vpc.foo.id}"
  zone             = "${cosmic_vpc.foo.zone}"
}

resource "cosmic_instance" "foo" {
  name             = "terraform-test"
  display_name     = "terraform-test"
  service_offering = "%s"
  network_id       = "${cosmic_network.foo.id}"
  template         = "%s"
  zone             = "${cosmic_vpc.foo.zone}"
  user_data        = "foobar\nfoo\nbar"
  expunge          = true
}

resource "cosmic_secondary_ipaddress" "foo" {
  virtual_machine_id = "${cosmic_instance.foo.id}"
}`,
	COSMIC_VPC_OFFERING,
	COSMIC_ZONE,
	COSMIC_VPC_NETWORK_OFFERING,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_TEMPLATE,
)

var testAccCosmicSecondaryIPAddress_fixedIP = fmt.Sprintf(`
resource "cosmic_vpc" "foo" {
  name           = "terraform-vpc"
  display_text   = "terraform-vpc"
  cidr           = "10.0.10.0/22"
  network_domain = "terraform-domain"
  vpc_offering   = "%s"
  zone           = "%s"
}

resource "cosmic_network" "foo" {
  name             = "terraform-network"
  cidr             = "10.0.10.0/24"
  gateway          = "10.0.10.1"
  network_offering = "%s"
  vpc_id           = "${cosmic_vpc.foo.id}"
  zone             = "${cosmic_vpc.foo.zone}"
}

resource "cosmic_instance" "foo" {
  name             = "terraform-test"
  display_name     = "terraform-test"
  service_offering = "%s"
  network_id       = "${cosmic_network.foo.id}"
  template         = "%s"
  zone             = "${cosmic_vpc.foo.zone}"
  user_data        = "foobar\nfoo\nbar"
  expunge          = true
}

resource "cosmic_secondary_ipaddress" "foo" {
  ip_address         = "10.0.10.10"
  virtual_machine_id = "${cosmic_instance.foo.id}"
}`,
	COSMIC_VPC_OFFERING,
	COSMIC_ZONE,
	COSMIC_VPC_NETWORK_OFFERING,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_TEMPLATE,
)
