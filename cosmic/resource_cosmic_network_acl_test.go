package cosmic

import (
	"fmt"
	"testing"

	"github.com/MissionCriticalCloud/go-cosmic/v6/cosmic"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccCosmicNetworkACL_basic(t *testing.T) {
	if COSMIC_VPC_OFFERING == "" {
		t.Skip("This test requires an existing VPC offering (set it by exporting COSMIC_VPC_OFFERING)")
	}

	var id string
	var acl cosmic.NetworkACLList

	createAttributes := &testAccCheckCosmicNetworkACLExpectedAttributes{
		Name:        "terraform-acl",
		Description: "terraform-acl-text",
	}

	updateAttributes := &testAccCheckCosmicNetworkACLExpectedAttributes{
		Name:        "terraform-acl-updated",
		Description: "terraform-acl-text-updated",
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicNetworkACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicNetworkACL_basic(createAttributes),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosmicNetworkACLExists("cosmic_network_acl.foo", &id, &acl),
					testAccCheckCosmicNetworkACLBasicAttributes(&acl, createAttributes),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl.foo", "name", createAttributes.Name),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl.foo", "description", createAttributes.Description),
				),
			},

			{
				Config: testAccCosmicNetworkACL_basic(updateAttributes),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCosmicNetworkACLExists("cosmic_network_acl.foo", &id, &acl),
					testAccCheckCosmicNetworkACLBasicAttributes(&acl, updateAttributes),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl.foo", "name", updateAttributes.Name),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl.foo", "description", updateAttributes.Description),
				),
			},
		},
	})
}

func testAccCheckCosmicNetworkACLExists(n string, id *string, acl *cosmic.NetworkACLList) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No network ACL ID is set")
		}

		if id != nil {
			if *id != "" && *id != rs.Primary.ID {
				return fmt.Errorf("Resource ID has changed")
			}

			*id = rs.Primary.ID
		}

		client := testAccProvider.Meta().(*CosmicClient)
		acllist, count, err := client.NetworkACL.GetNetworkACLListByID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if count == 0 {
			return fmt.Errorf("Network ACL not found")
		}

		*acl = *acllist

		return nil
	}
}

type testAccCheckCosmicNetworkACLExpectedAttributes struct {
	Description string
	Name        string
}

func testAccCheckCosmicNetworkACLBasicAttributes(acl *cosmic.NetworkACLList, want *testAccCheckCosmicNetworkACLExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if acl.Name != want.Name {
			return fmt.Errorf("Bad name: got %s; want %s", acl.Name, want.Name)
		}

		if acl.Description != want.Description {
			return fmt.Errorf("Bad name: got %s; want %s", acl.Description, want.Description)
		}

		return nil
	}
}

func testAccCheckCosmicNetworkACLDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CosmicClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cosmic_network_acl" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No network ACL ID is set")
		}

		_, _, err := client.NetworkACL.GetNetworkACLListByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Network ACl list %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCosmicNetworkACL_basic(attr *testAccCheckCosmicNetworkACLExpectedAttributes) string {
	return fmt.Sprintf(`
resource "cosmic_vpc" "foo" {
  name           = "terraform-vpc"
  display_text   = "terraform-vpc"
  cidr           = "10.0.10.0/22"
  vpc_offering   = "%s"
  network_domain = "terraform-domain"
}

resource "cosmic_network_acl" "foo" {
  name        = "%s"
  description = "%s"
  vpc_id      = "${cosmic_vpc.foo.id}"
}`,
		COSMIC_VPC_OFFERING,
		attr.Name,
		attr.Description,
	)
}
