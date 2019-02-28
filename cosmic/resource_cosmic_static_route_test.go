package cosmic

import (
	"fmt"
	"testing"

	"github.com/MissionCriticalCloud/go-cosmic/cosmic"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCosmicStaticRoute_basic(t *testing.T) {
	if COSMIC_VPC_ID == "" {
		t.Skip("This test requires an existing VPC ID (set it by exporting COSMIC_VPC_ID)")
	}

	var route cosmic.StaticRoute

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicStaticRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicStaticRoute_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicStaticRouteExists(
						"cosmic_static_route.foo", &route),
					testAccCheckCosmicStaticRouteAttributes(&route),
				),
			},
		},
	})
}

func testAccCheckCosmicStaticRouteExists(
	n string, route *cosmic.StaticRoute) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Static Route ID is set")
		}

		cs := testAccProvider.Meta().(*cosmic.CosmicClient)
		r, _, err := cs.VPC.GetStaticRouteByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if r.Id != rs.Primary.ID {
			return fmt.Errorf("Static Route not found")
		}

		*route = *r

		return nil
	}
}

func testAccCheckCosmicStaticRouteAttributes(
	route *cosmic.StaticRoute) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if route.Cidr != "172.16.0.0/16" {
			return fmt.Errorf("Bad CIDR: %s", route.Cidr)
		}

		return nil
	}
}

func testAccCheckCosmicStaticRouteDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cosmic.CosmicClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cosmic_static_route" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No static route ID is set")
		}

		route, _, err := cs.VPC.GetStaticRouteByID(rs.Primary.ID)
		if err == nil && route.Id != "" {
			return fmt.Errorf("Static route %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

var testAccCosmicStaticRoute_basic = fmt.Sprintf(`
resource "cosmic_static_route" "foo" {
  cidr    = "172.16.0.0/16"
  nexthop = "10.0.252.1"
  vpc_id  = "%s"
}`, COSMIC_VPC_ID)
