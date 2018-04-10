package cosmic

import (
	"fmt"
	"testing"

	"github.com/MissionCriticalCloud/go-cosmic/cosmic"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCosmicSecurityGroup_basic(t *testing.T) {
	var sg cosmic.SecurityGroup
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicSecurityGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCosmicSecurityGroup_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicSecurityGroupExists(
						"cosmic_security_group.foo", &sg),
					testAccCheckCosmicSecurityGroupBasicAttributes(&sg),
				),
			},
		},
	})
}

func testAccCheckCosmicSecurityGroupExists(
	n string, sg *cosmic.SecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No security group ID is set")
		}

		cs := testAccProvider.Meta().(*cosmic.CosmicClient)
		resp, _, err := cs.SecurityGroup.GetSecurityGroupByID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp.Id != rs.Primary.ID {
			return fmt.Errorf("Network ACL not found")
		}

		*sg = *resp

		return nil
	}
}

func testAccCheckCosmicSecurityGroupBasicAttributes(
	sg *cosmic.SecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if sg.Name != "terraform-security-group" {
			return fmt.Errorf("Bad name: %s", sg.Name)
		}

		if sg.Description != "terraform-security-group-text" {
			return fmt.Errorf("Bad description: %s", sg.Description)
		}

		return nil
	}
}

func testAccCheckCosmicSecurityGroupDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cosmic.CosmicClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cosmic_security_group" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No security group ID is set")
		}

		_, _, err := cs.SecurityGroup.GetSecurityGroupByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Security group list %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

var testAccCosmicSecurityGroup_basic = fmt.Sprintf(`
resource "cosmic_security_group" "foo" {
  name = "terraform-security-group"
	description = "terraform-security-group-text"
}`)
