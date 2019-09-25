package cosmic

import (
	"fmt"
	"testing"

	"github.com/MissionCriticalCloud/go-cosmic/v6/cosmic"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccCosmicTemplate_basic(t *testing.T) {
	var template cosmic.Template

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicTemplate_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicTemplateExists("cosmic_template.foo", &template),
					testAccCheckCosmicTemplateBasicAttributes(&template),
					resource.TestCheckResourceAttr(
						"cosmic_template.foo", "display_text", "terraform-test"),
				),
			},
		},
	})
}

func TestAccCosmicTemplate_update(t *testing.T) {
	var template cosmic.Template

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosmicTemplate_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicTemplateExists("cosmic_template.foo", &template),
					testAccCheckCosmicTemplateBasicAttributes(&template),
				),
			},

			{
				Config: testAccCosmicTemplate_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicTemplateExists(
						"cosmic_template.foo", &template),
					testAccCheckCosmicTemplateUpdatedAttributes(&template),
					resource.TestCheckResourceAttr(
						"cosmic_template.foo", "display_text", "terraform-updated"),
				),
			},
		},
	})
}

func testAccCheckCosmicTemplateExists(n string, template *cosmic.Template) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No template ID is set")
		}

		client := testAccProvider.Meta().(*CosmicClient)
		tmpl, _, err := client.Template.GetTemplateByID(rs.Primary.ID, "executable")

		if err != nil {
			return err
		}

		if tmpl.Id != rs.Primary.ID {
			return fmt.Errorf("Template not found")
		}

		*template = *tmpl

		return nil
	}
}

func testAccCheckCosmicTemplateBasicAttributes(template *cosmic.Template) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if template.Name != "terraform-test" {
			return fmt.Errorf("Bad name: %s", template.Name)
		}

		if template.Format != "QCOW2" {
			return fmt.Errorf("Bad format: %s", template.Format)
		}

		if template.Hypervisor != "KVM" {
			return fmt.Errorf("Bad hypervisor: %s", template.Hypervisor)
		}

		if template.Ostypename != "Other PV (64-bit)" {
			return fmt.Errorf("Bad os type: %s", template.Ostypename)
		}

		if template.Zonename != COSMIC_ZONE {
			return fmt.Errorf("Bad zone: %s", template.Zonename)
		}

		return nil
	}
}

func testAccCheckCosmicTemplateUpdatedAttributes(template *cosmic.Template) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if template.Displaytext != "terraform-updated" {
			return fmt.Errorf("Bad name: %s", template.Displaytext)
		}

		if !template.Isdynamicallyscalable {
			return fmt.Errorf("Bad is_dynamically_scalable: %t", template.Isdynamicallyscalable)
		}

		if !template.Passwordenabled {
			return fmt.Errorf("Bad password_enabled: %t", template.Passwordenabled)
		}

		return nil
	}
}

func testAccCheckCosmicTemplateDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CosmicClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cosmic_template" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No template ID is set")
		}

		_, _, err := client.Template.GetTemplateByID(rs.Primary.ID, "executable")
		if err == nil {
			return fmt.Errorf("Template %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

var testAccCosmicTemplate_basic = fmt.Sprintf(`
resource "cosmic_template" "foo" {
  name       = "terraform-test"
  format     = "QCOW2"
  hypervisor = "KVM"
  os_type    = "Other PV (64-bit)"
  url        = "http://cloud.centos.org/centos/7/images/CentOS-7-x86_64-GenericCloud.qcow2"
  zone       = "%s"
}`,
	COSMIC_ZONE,
)

var testAccCosmicTemplate_update = fmt.Sprintf(`
resource "cosmic_template" "foo" {
  name                    = "terraform-test"
  display_text            = "terraform-updated"
  format                  = "QCOW2"
  hypervisor              = "KVM"
  os_type                 = "Other PV (64-bit)"
  url                     = "http://cloud.centos.org/centos/7/images/CentOS-7-x86_64-GenericCloud.qcow2"
  zone                    = "%s"
  is_dynamically_scalable = true
  password_enabled        = true
}`,
	COSMIC_ZONE,
)
