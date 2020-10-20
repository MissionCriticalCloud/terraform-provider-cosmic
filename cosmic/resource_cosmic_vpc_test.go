package cosmic

import (
	"fmt"
	"strings"
	"testing"

	"github.com/MissionCriticalCloud/go-cosmic/v6/cosmic"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCosmicVPC_basic(t *testing.T) {
	if COSMIC_VPC_OFFERING == "" {
		t.Skip("This test requires an existing VPC offering (set it by exporting COSMIC_VPC_OFFERING)")
	}

	var vpc cosmic.VPC

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicVPCDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicVPC_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicVPCExists(
						"cosmic_vpc.foo", &vpc),
					testAccCheckCosmicVPCAttributes(&vpc),
					resource.TestCheckResourceAttr(
						"cosmic_vpc.foo", "vpc_offering", COSMIC_VPC_OFFERING),
				),
			},
		},
	})
}

func testAccCheckCosmicVPCExists(n string, vpc *cosmic.VPC) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPC ID is set")
		}

		client := testAccProvider.Meta().(*CosmicClient)
		v, _, err := client.VPC.GetVPCByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if v.Id != rs.Primary.ID {
			return fmt.Errorf("VPC not found")
		}

		*vpc = *v

		return nil
	}
}

func testAccCheckCosmicVPCAttributes(vpc *cosmic.VPC) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if vpc.Name != "terraform-vpc" {
			return fmt.Errorf("Bad name: %s", vpc.Name)
		}

		if vpc.Displaytext != "terraform-vpc-text" {
			return fmt.Errorf("Bad display text: %s", vpc.Displaytext)
		}

		if vpc.Cidr != "10.0.10.0/22" {
			return fmt.Errorf("Bad VPC CIDR: %s", vpc.Cidr)
		}

		if vpc.Networkdomain != "terraform-domain" {
			return fmt.Errorf("Bad network domain: %s", vpc.Networkdomain)
		}

		return nil
	}
}

func testAccCheckCosmicVPCDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CosmicClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cosmic_vpc" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPC ID is set")
		}

		_, _, err := client.VPC.GetVPCByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("VPC %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

var testAccCosmicVPC_basic = fmt.Sprintf(`
resource "cosmic_vpc" "foo" {
  name           = "terraform-vpc"
  display_text   = "terraform-vpc-text"
  cidr           = "10.0.10.0/22"
  vpc_offering   = "%s"
  network_domain = "terraform-domain"
}`,
	strings.ToLower(COSMIC_VPC_OFFERING),
)
