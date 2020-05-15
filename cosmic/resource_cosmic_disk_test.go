package cosmic

import (
	"fmt"
	"strings"
	"testing"

	"github.com/MissionCriticalCloud/go-cosmic/v6/cosmic"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccCosmicDisk_basic(t *testing.T) {
	if COSMIC_DISK_OFFERING == "" {
		t.Skip("This test requires an existing disk offering (set it by exporting COSMIC_DISK_OFFERING)")
	}

	var disk cosmic.Volume

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicDisk_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicDiskExists(
						"cosmic_disk.foo", &disk),
					testAccCheckCosmicDiskAttributes(&disk),
				),
			},
		},
	})
}

func TestAccCosmicDisk_update(t *testing.T) {
	if COSMIC_DISK_OFFERING == "" {
		t.Skip("This test requires an existing disk offering (set it by exporting COSMIC_DISK_OFFERING)")
	}

	var disk cosmic.Volume

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicDisk_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicDiskExists(
						"cosmic_disk.foo", &disk),
					testAccCheckCosmicDiskAttributes(&disk),
					resource.TestCheckResourceAttr(
						"cosmic_disk.foo", "size", "10"),
				),
			},

			{
				Config: testAccCosmicDisk_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicDiskExists(
						"cosmic_disk.foo", &disk),
					resource.TestCheckResourceAttr(
						"cosmic_disk.foo", "size", "20"),
				),
			},
		},
	})
}

func TestAccCosmicDisk_attachBasic(t *testing.T) {
	if COSMIC_DISK_OFFERING == "" {
		t.Skip("This test requires an existing disk offering (set it by exporting COSMIC_DISK_OFFERING)")
	}

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

	var disk cosmic.Volume

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicDisk_attachBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicDiskExists(
						"cosmic_disk.foo", &disk),
					testAccCheckCosmicDiskAttributes(&disk),
				),
			},
		},
	})
}

func TestAccCosmicDisk_attachUpdate(t *testing.T) {
	if COSMIC_DISK_OFFERING == "" {
		t.Skip("This test requires an existing disk offering (set it by exporting COSMIC_DISK_OFFERING)")
	}

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

	var disk cosmic.Volume

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicDisk_attachBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicDiskExists(
						"cosmic_disk.foo", &disk),
					testAccCheckCosmicDiskAttributes(&disk),
					resource.TestCheckResourceAttr(
						"cosmic_disk.foo", "size", "10"),
				),
			},

			{
				Config: testAccCosmicDisk_attachUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicDiskExists(
						"cosmic_disk.foo", &disk),
					resource.TestCheckResourceAttr(
						"cosmic_disk.foo", "size", "20"),
				),
			},
		},
	})
}

func TestAccCosmicDisk_attachDeviceID(t *testing.T) {
	if COSMIC_DISK_OFFERING == "" {
		t.Skip("This test requires an existing disk offering (set it by exporting COSMIC_DISK_OFFERING)")
	}

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

	var disk cosmic.Volume

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicDisk_attachDeviceID,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicDiskExists(
						"cosmic_disk.foo", &disk),
					testAccCheckCosmicDiskAttributes(&disk),
					resource.TestCheckResourceAttr(
						"cosmic_disk.foo", "device_id", "4"),
				),
			},
		},
	})
}

func TestAccCosmicDisk_diskController(t *testing.T) {
	if COSMIC_DISK_OFFERING == "" {
		t.Skip("This test requires an existing disk offering (set it by exporting COSMIC_DISK_OFFERING)")
	}

	var disk cosmic.Volume

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicDisk_diskController,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicDiskExists(
						"cosmic_disk.foo", &disk),
					testAccCheckCosmicDiskAttributes(&disk),
					resource.TestCheckResourceAttr(
						"cosmic_disk.foo", "disk_controller", "SCSI"),
				),
			},
		},
	})
}

func TestAccCosmicDisk_attachDiskController(t *testing.T) {
	if COSMIC_DISK_OFFERING == "" {
		t.Skip("This test requires an existing disk offering (set it by exporting COSMIC_DISK_OFFERING)")
	}

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

	var disk cosmic.Volume

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicDisk_attachDiskController,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicDiskExists(
						"cosmic_disk.foo", &disk),
					testAccCheckCosmicDiskAttributes(&disk),
					resource.TestCheckResourceAttr(
						"cosmic_disk.foo", "disk_controller", "SCSI"),
				),
			},
		},
	})
}

func testAccCheckCosmicDiskExists(n string, disk *cosmic.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No disk ID is set")
		}

		client := testAccProvider.Meta().(*CosmicClient)
		volume, _, err := client.Volume.GetVolumeByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if volume.Id != rs.Primary.ID {
			return fmt.Errorf("Disk not found")
		}

		*disk = *volume

		return nil
	}
}

func testAccCheckCosmicDiskAttributes(disk *cosmic.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if disk.Name != "terraform-disk" {
			return fmt.Errorf("Bad name: %s", disk.Name)
		}

		if disk.Diskofferingname != COSMIC_DISK_OFFERING {
			return fmt.Errorf("Bad disk offering: %s", disk.Diskofferingname)
		}

		return nil
	}
}

func testAccCheckCosmicDiskDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CosmicClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cosmic_disk" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No disk ID is set")
		}

		_, _, err := client.Volume.GetVolumeByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Disk %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

var testAccCosmicDisk_basic = fmt.Sprintf(`
resource "cosmic_disk" "foo" {
  name          = "terraform-disk"
  attach        = false
  size          = "10"
  disk_offering = "%s"
}`,
	strings.ToLower(COSMIC_DISK_OFFERING),
)

var testAccCosmicDisk_update = fmt.Sprintf(`
resource "cosmic_disk" "foo" {
  name          = "terraform-disk"
  attach        = false
  size          = "20"
  disk_offering = "%s"
}`,
	COSMIC_DISK_OFFERING,
)

var testAccCosmicDisk_diskController = fmt.Sprintf(`
resource "cosmic_disk" "foo" {
  name            = "terraform-disk"
  attach          = false
  size            = "10"
  disk_offering   = "%s"
  disk_controller = "SCSI"
}`,
	COSMIC_DISK_OFFERING,
)

var testAccCosmicDisk_attachBasic = fmt.Sprintf(`
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
  expunge          = true
}

resource "cosmic_disk" "foo" {
  name               = "terraform-disk"
  attach             = true
  size               = "10"
  disk_offering      = "%s"
  virtual_machine_id = "${cosmic_instance.foo.id}"
}`,
	COSMIC_VPC_OFFERING,
	COSMIC_VPC_NETWORK_OFFERING,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_TEMPLATE,
	COSMIC_DISK_OFFERING,
)

var testAccCosmicDisk_attachDiskController = fmt.Sprintf(`
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
  expunge          = true
}

resource "cosmic_disk" "foo" {
  name               = "terraform-disk"
  attach             = true
  size               = "10"
  disk_offering      = "%s"
  disk_controller    = "SCSI"
  virtual_machine_id = "${cosmic_instance.foo.id}"
}`,
	COSMIC_VPC_OFFERING,
	COSMIC_VPC_NETWORK_OFFERING,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_TEMPLATE,
	COSMIC_DISK_OFFERING,
)

var testAccCosmicDisk_attachDeviceID = fmt.Sprintf(`
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
  expunge          = true
}

resource "cosmic_disk" "foo" {
  name               = "terraform-disk"
  attach             = true
  device_id          = 4
  size               = "10"
  disk_offering      = "%s"
  virtual_machine_id = "${cosmic_instance.foo.id}"
}`,
	COSMIC_VPC_OFFERING,
	COSMIC_VPC_NETWORK_OFFERING,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_TEMPLATE,
	COSMIC_DISK_OFFERING,
)

var testAccCosmicDisk_attachUpdate = fmt.Sprintf(`
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
  expunge          = true
}

resource "cosmic_disk" "foo" {
  name               = "terraform-disk"
  attach             = true
  size               = "20"
  disk_offering      = "%s"
  virtual_machine_id = "${cosmic_instance.foo.id}"
}`,
	COSMIC_VPC_OFFERING,
	COSMIC_VPC_NETWORK_OFFERING,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_TEMPLATE,
	COSMIC_DISK_OFFERING,
)
