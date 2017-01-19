package cosmic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/xanzy/go-cosmic/cosmic"
)

func TestAccCosmicStaticRoute_basic(t *testing.T) {
	var route cosmic.StaticRoute

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicStaticRouteDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCosmicStaticRoute_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicStaticRouteExists(
						"cosmic_static_route.bar", &route),
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

		if route.Cidr != CLOUDSTACK_STATIC_ROUTE_CIDR {
			return fmt.Errorf("Bad Cidr: %s", route.Cidr)
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
resource "cosmic_vpc" "foobar" {
  name = "terraform-vpc"
  cidr = "%s"
  vpc_offering = "%s"
  zone = "%s"
}

resource "cosmic_network" "foo" {
	name = "terraform-network"
	cidr = "%s"
	network_offering = "%s"
	vpc_id = "${cosmic_vpc.foobar.id}"
	zone = "${cosmic_vpc.foobar.zone}"
}

resource "cosmic_static_route" "bar" {
  cidr = "%s"
	nexthop = "%s"
  vpc_id = "${cosmic_vpc.foobar.id}"
}`,
	CLOUDSTACK_VPC_CIDR_1,
	CLOUDSTACK_VPC_OFFERING,
	CLOUDSTACK_ZONE,
	CLOUDSTACK_VPC_NETWORK_CIDR,
	CLOUDSTACK_VPC_NETWORK_OFFERING,
	CLOUDSTACK_STATIC_ROUTE_CIDR,
	CLOUDSTACK_VPC_NETWORK_IPADDRESS)
