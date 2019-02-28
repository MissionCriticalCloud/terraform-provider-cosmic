package cosmic

import (
	"fmt"
	"testing"

	"github.com/MissionCriticalCloud/go-cosmic/cosmic"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCosmicDisk_basic(t *testing.T) {
	if COSMIC_DISK_OFFERING_1 == "" {
		t.Skip("This test requires an existing disk offering (set it by exporting COSMIC_DISK_OFFERING_1)")
	}

	var disk cosmic.Volume

	resource.Test(t, resource.TestCase{
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
	if COSMIC_DISK_OFFERING_1 == "" {
		t.Skip("This test requires an existing disk offering (set it by exporting COSMIC_DISK_OFFERING_1)")
	}

	var disk cosmic.Volume

	resource.Test(t, resource.TestCase{
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
	if COSMIC_DISK_OFFERING_1 == "" {
		t.Skip("This test requires an existing disk offering (set it by exporting COSMIC_DISK_OFFERING_1)")
	}

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

	var disk cosmic.Volume

	resource.Test(t, resource.TestCase{
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

	var disk cosmic.Volume

	resource.Test(t, resource.TestCase{
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

	var disk cosmic.Volume

	resource.Test(t, resource.TestCase{
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

func testAccCheckCosmicDiskExists(
	n string, disk *cosmic.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No disk ID is set")
		}

		cs := testAccProvider.Meta().(*cosmic.CosmicClient)
		volume, _, err := cs.Volume.GetVolumeByID(rs.Primary.ID)

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

func testAccCheckCosmicDiskAttributes(
	disk *cosmic.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if disk.Name != "terraform-disk" {
			return fmt.Errorf("Bad name: %s", disk.Name)
		}

		if disk.Diskofferingname != COSMIC_DISK_OFFERING_1 {
			return fmt.Errorf("Bad disk offering: %s", disk.Diskofferingname)
		}

		return nil
	}
}

func testAccCheckCosmicDiskDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cosmic.CosmicClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cosmic_disk" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No disk ID is set")
		}

		_, _, err := cs.Volume.GetVolumeByID(rs.Primary.ID)
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
  zone          = "%s"
}`,
	COSMIC_DISK_OFFERING_1,
	COSMIC_ZONE)

var testAccCosmicDisk_update = fmt.Sprintf(`
resource "cosmic_disk" "foo" {
  name          = "terraform-disk"
  attach        = false
  size          = "20"
  disk_offering = "%s"
  zone          = "%s"
}`,
	COSMIC_DISK_OFFERING_1,
	COSMIC_ZONE)

var testAccCosmicDisk_attachBasic = fmt.Sprintf(`
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
  display_name     = "terraform"
  service_offering = "%s"
  network_id       = "${cosmic_network.foo.id}"
  template         = "%s"
  zone             = "${cosmic_network.foo.zone}"
  expunge          = true
}

resource "cosmic_disk" "foo" {
  name               = "terraform-disk"
  attach             = true
  size               = "10"
  disk_offering      = "%s"
  virtual_machine_id = "${cosmic_instance.foo.id}"
  zone               = "${cosmic_instance.foo.zone}"
}`,
	COSMIC_VPC_NETWORK_OFFERING,
	COSMIC_VPC_ID,
	COSMIC_ZONE,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_TEMPLATE,
	COSMIC_DISK_OFFERING_1)

var testAccCosmicDisk_attachDeviceID = fmt.Sprintf(`
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
  display_name     = "terraform"
  service_offering = "%s"
  network_id       = "${cosmic_network.foo.id}"
  template         = "%s"
  zone             = "${cosmic_network.foo.zone}"
  expunge          = true
}

resource "cosmic_disk" "foo" {
  name               = "terraform-disk"
  attach             = true
  device_id          = 4
  size               = "10"
  disk_offering      = "%s"
  virtual_machine_id = "${cosmic_instance.foo.id}"
  zone               = "${cosmic_instance.foo.zone}"
}`,
	COSMIC_VPC_NETWORK_OFFERING,
	COSMIC_VPC_ID,
	COSMIC_ZONE,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_TEMPLATE,
	COSMIC_DISK_OFFERING_1)

var testAccCosmicDisk_attachUpdate = fmt.Sprintf(`
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
  display_name     = "terraform"
  service_offering = "%s"
  network_id       = "${cosmic_network.foo.id}"
  template         = "%s"
  zone             = "${cosmic_network.foo.zone}"
  expunge          = true
}

resource "cosmic_disk" "foo" {
  name               = "terraform-disk"
  attach             = true
  size               = "20"
  disk_offering      = "%s"
  virtual_machine_id = "${cosmic_instance.foo.id}"
  zone               = "${cosmic_instance.foo.zone}"
}`,
	COSMIC_VPC_NETWORK_OFFERING,
	COSMIC_VPC_ID,
	COSMIC_ZONE,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_TEMPLATE,
	COSMIC_DISK_OFFERING_1)
