package cosmic

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccCosmicNetworkACLRule_basic(t *testing.T) {
	if COSMIC_VPC_OFFERING == "" {
		t.Skip("This test requires an existing VPC offering (set it by exporting COSMIC_VPC_OFFERING)")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicNetworkACLRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicNetworkACLRule_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicNetworkACLRulesExist("cosmic_network_acl.foo"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.1480917538.action", "allow"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.1480917538.cidr_list.#", "1"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.1480917538.cidr_list.3056857544", "172.18.100.0/24"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.1480917538.icmp_code", "-1"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.1480917538.icmp_type", "-1"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.1480917538.ports.#", "0"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.1480917538.protocol", "icmp"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.1480917538.traffic_type", "ingress"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2898748868.action", "allow"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2898748868.cidr_list.#", "1"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2898748868.cidr_list.2835005819", "172.16.100.0/24"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2898748868.ports.#", "2"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2898748868.ports.1889509032", "80"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2898748868.ports.3638101695", "443"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2898748868.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2898748868.traffic_type", "ingress"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.3943933455.action", "allow"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.3943933455.cidr_list.#", "1"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.3943933455.cidr_list.3056857544", "172.18.100.0/24"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.3943933455.ports.#", "0"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.3943933455.protocol", "all"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.3943933455.traffic_type", "ingress"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.865820176.action", "allow"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.865820176.cidr_list.#", "1"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.865820176.cidr_list.3056857544", "172.18.100.0/24"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.865820176.ports.#", "0"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.865820176.protocol", "47"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.865820176.traffic_type", "ingress"),
				),
			},
		},
	})
}

func TestAccCosmicNetworkACLRule_update(t *testing.T) {
	if COSMIC_VPC_OFFERING == "" {
		t.Skip("This test requires an existing VPC offering (set it by exporting COSMIC_VPC_OFFERING)")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicNetworkACLRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicNetworkACLRule_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicNetworkACLRulesExist("cosmic_network_acl.foo"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.1480917538.action", "allow"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.1480917538.cidr_list.#", "1"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.1480917538.cidr_list.3056857544", "172.18.100.0/24"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.1480917538.icmp_code", "-1"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.1480917538.icmp_type", "-1"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.1480917538.ports.#", "0"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.1480917538.protocol", "icmp"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.1480917538.traffic_type", "ingress"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2898748868.action", "allow"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2898748868.cidr_list.#", "1"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2898748868.cidr_list.2835005819", "172.16.100.0/24"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2898748868.ports.#", "2"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2898748868.ports.1889509032", "80"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2898748868.ports.3638101695", "443"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2898748868.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2898748868.traffic_type", "ingress"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.3943933455.action", "allow"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.3943933455.cidr_list.#", "1"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.3943933455.cidr_list.3056857544", "172.18.100.0/24"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.3943933455.ports.#", "0"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.3943933455.protocol", "all"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.3943933455.traffic_type", "ingress"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.865820176.action", "allow"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.865820176.cidr_list.#", "1"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.865820176.cidr_list.3056857544", "172.18.100.0/24"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.865820176.ports.#", "0"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.865820176.protocol", "47"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.865820176.traffic_type", "ingress"),
				),
			},

			{
				Config: testAccCosmicNetworkACLRule_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicNetworkACLRulesExist("cosmic_network_acl.foo"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.#", "5"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.1724235854.action", "deny"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.1724235854.cidr_list.#", "1"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.1724235854.cidr_list.3482919157", "10.0.0.0/24"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.1724235854.icmp_code", "0"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.1724235854.icmp_type", "0"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.1724235854.ports.#", "2"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.1724235854.ports.1209010669", "1000-2000"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.1724235854.ports.1889509032", "80"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.1724235854.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.1724235854.traffic_type", "egress"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2090315355.action", "deny"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2090315355.cidr_list.#", "2"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2090315355.cidr_list.2104435309", "172.18.101.0/24"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2090315355.cidr_list.3056857544", "172.18.100.0/24"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2090315355.icmp_code", "-1"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2090315355.icmp_type", "-1"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2090315355.ports.#", "0"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2090315355.protocol", "icmp"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2090315355.traffic_type", "ingress"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2548582181.action", "deny"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2548582181.cidr_list.#", "1"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2548582181.cidr_list.3056857544", "172.18.100.0/24"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2548582181.icmp_code", "0"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2548582181.icmp_type", "0"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2548582181.ports.#", "0"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2548582181.protocol", "47"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2548582181.traffic_type", "ingress"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2576683033.action", "allow"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2576683033.cidr_list.#", "1"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2576683033.cidr_list.3056857544", "172.18.100.0/24"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2576683033.icmp_code", "0"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2576683033.icmp_type", "0"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2576683033.ports.#", "2"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2576683033.ports.1889509032", "80"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2576683033.ports.3638101695", "443"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2576683033.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.2576683033.traffic_type", "ingress"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.3171160373.action", "deny"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.3171160373.cidr_list.#", "1"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.3171160373.cidr_list.3056857544", "172.18.100.0/24"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.3171160373.icmp_code", "0"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.3171160373.icmp_type", "0"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.3171160373.ports.#", "0"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.3171160373.protocol", "all"),
					resource.TestCheckResourceAttr(
						"cosmic_network_acl_rule.foo", "rule.3171160373.traffic_type", "ingress"),
				),
			},
		},
	})
}

func testAccCheckCosmicNetworkACLRulesExist(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No network ACL rule ID is set")
		}

		for k, id := range rs.Primary.Attributes {
			if !strings.Contains(k, ".uuids.") || strings.HasSuffix(k, ".uuids.%") {
				continue
			}

			client := testAccProvider.Meta().(*CosmicClient)
			_, count, err := client.NetworkACL.GetNetworkACLByID(id)

			if err != nil {
				return err
			}

			if count == 0 {
				return fmt.Errorf("Network ACL rule %s not found", k)
			}
		}

		return nil
	}
}

func testAccCheckCosmicNetworkACLRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CosmicClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cosmic_network_acl_rule" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No network ACL rule ID is set")
		}

		for k, id := range rs.Primary.Attributes {
			if !strings.Contains(k, ".uuids.") || strings.HasSuffix(k, ".uuids.%") {
				continue
			}

			_, _, err := client.NetworkACL.GetNetworkACLByID(id)
			if err == nil {
				return fmt.Errorf("Network ACL rule %s still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

var testAccCosmicNetworkACLRule_basic = fmt.Sprintf(`
resource "cosmic_vpc" "foo" {
  name           = "terraform-vpc"
  display_text   = "terraform-vpc"
  cidr           = "10.0.10.0/22"
  vpc_offering   = "%s"
  network_domain = "terraform-domain"
  zone           = "%s"
}

resource "cosmic_network_acl" "foo" {
  name        = "terraform-acl"
  description = "terraform-acl-text"
  vpc_id      = "${cosmic_vpc.foo.id}"
}

resource "cosmic_network_acl_rule" "foo" {
  acl_id = "${cosmic_network_acl.foo.id}"

  rule {
    action       = "allow"
    cidr_list    = ["172.18.100.0/24"]
    protocol     = "all"
    traffic_type = "ingress"
  }

  rule {
    action       = "allow"
    cidr_list    = ["172.18.100.0/24"]
    protocol     = "icmp"
    icmp_type    = "-1"
    icmp_code    = "-1"
    traffic_type = "ingress"
  }

  rule {
    action       = "allow"
    cidr_list    = ["172.18.100.0/24"]
    protocol     = "47"
    traffic_type = "ingress"
  }

  rule {
    cidr_list    = ["172.16.100.0/24"]
    protocol     = "tcp"
    ports        = ["80", "443"]
    traffic_type = "ingress"
  }
}`,
	COSMIC_VPC_OFFERING,
	COSMIC_ZONE,
)

var testAccCosmicNetworkACLRule_update = fmt.Sprintf(`
resource "cosmic_vpc" "foo" {
  name           = "terraform-vpc"
  display_text   = "terraform-vpc"
  cidr           = "10.0.10.0/22"
  vpc_offering   = "%s"
  network_domain = "terraform-domain"
  zone           = "%s"
}

resource "cosmic_network_acl" "foo" {
  name        = "terraform-acl"
  description = "terraform-acl-text"
  vpc_id      = "${cosmic_vpc.foo.id}"
}

resource "cosmic_network_acl_rule" "foo" {
  acl_id = "${cosmic_network_acl.foo.id}"

  rule {
    action       = "deny"
    cidr_list    = ["172.18.100.0/24"]
    protocol     = "all"
    traffic_type = "ingress"
  }

  rule {
    action       = "deny"
    cidr_list    = ["172.18.100.0/24", "172.18.101.0/24"]
    protocol     = "icmp"
    icmp_type    = "-1"
    icmp_code    = "-1"
    traffic_type = "ingress"
  }

  rule {
    action       = "deny"
    cidr_list    = ["172.18.100.0/24"]
    protocol     = "47"
    traffic_type = "ingress"
  }

  rule {
    action       = "allow"
    cidr_list    = ["172.18.100.0/24"]
    protocol     = "tcp"
    ports        = ["80", "443"]
    traffic_type = "ingress"
  }

  rule {
    action       = "deny"
    cidr_list    = ["10.0.0.0/24"]
    protocol     = "tcp"
    ports        = ["80", "1000-2000"]
    traffic_type = "egress"
  }
}`,
	COSMIC_VPC_OFFERING,
	COSMIC_ZONE,
)
