package cosmic

import (
	"fmt"
	"testing"

	"github.com/MissionCriticalCloud/go-cosmic/v6/cosmic"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCosmicAffinityGroup_basic(t *testing.T) {
	var affinityGroup cosmic.AffinityGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicAffinityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicAffinityGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicAffinityGroupExists("cosmic_affinity_group.foo", &affinityGroup),
					testAccCheckCosmicAffinityGroupAttributes(&affinityGroup),
				),
			},
		},
	})
}

func testAccCheckCosmicAffinityGroupExists(
	n string, affinityGroup *cosmic.AffinityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No affinity group ID is set")
		}

		cs := testAccProvider.Meta().(*cosmic.CosmicClient)
		ag, _, err := cs.AffinityGroup.GetAffinityGroupByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if ag.Id != rs.Primary.ID {
			return fmt.Errorf("Affinity group not found")
		}

		*affinityGroup = *ag

		return nil
	}
}

func testAccCheckCosmicAffinityGroupAttributes(
	affinityGroup *cosmic.AffinityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if affinityGroup.Name != "terraform-affinity-group" {
			return fmt.Errorf("Bad name: %s", affinityGroup.Name)
		}

		if affinityGroup.Description != "terraform-affinity-group" {
			return fmt.Errorf("Bad description: %s", affinityGroup.Description)
		}

		if affinityGroup.Type != "host anti-affinity" {
			return fmt.Errorf("Bad type: %s", affinityGroup.Type)
		}

		return nil
	}
}

func testAccCheckCosmicAffinityGroupDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cosmic.CosmicClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cosmic_affinity_group" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No affinity group ID is set")
		}

		_, _, err := cs.AffinityGroup.GetAffinityGroupByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Affinity group %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

const testAccCosmicAffinityGroup = `
resource "cosmic_affinity_group" "foo" {
  name = "terraform-affinity-group"
  type = "host anti-affinity"
}`
