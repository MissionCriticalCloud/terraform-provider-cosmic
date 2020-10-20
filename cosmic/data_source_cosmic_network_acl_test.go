package cosmic

import (
	"fmt"
	"testing"

	"github.com/MissionCriticalCloud/go-cosmic/v6/cosmic"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceCosmicNetworkACL_basic(t *testing.T) {
	var aclList cosmic.NetworkACLList

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCosmicNetworkACL_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicNetworkACLDataSourceExists("data.cosmic_network_acl.default_allow", &aclList),
					testAccCheckCosmicNetworkACLDataSourceAttributes(&aclList),
				),
			},
		},
	})
}

func testAccCheckCosmicNetworkACLDataSourceExists(n string, aclList *cosmic.NetworkACLList) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Network ACL List data source ID not set")
		}

		client := testAccProvider.Meta().(*CosmicClient)

		list, _, err := client.NetworkACL.GetNetworkACLListByID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if list.Id != rs.Primary.ID {
			return fmt.Errorf("Network ACL List not found")
		}

		*aclList = *list

		return nil
	}
}

func testAccCheckCosmicNetworkACLDataSourceAttributes(aclList *cosmic.NetworkACLList) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if aclList.Name != "default_allow" {
			return fmt.Errorf("Bad name: %s", aclList.Name)
		}

		if aclList.Description != "Default Network ACL Allow All" {
			return fmt.Errorf("Bad description: %s", aclList.Description)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "cosmic_network_acl" {
				continue
			}

			if aclList.Id != rs.Primary.ID {
				return fmt.Errorf("Bad Network ACL List ID: %s", aclList.Id)
			}
		}

		return nil
	}
}

const testAccDataSourceCosmicNetworkACL_basic = `
data "cosmic_network_acl" "default_allow" {
  filter {
    name  = "name"
    value = "default_allow"
  }
}
`
