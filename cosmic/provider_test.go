package cosmic

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"cosmic": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testSetValueOnResourceData(t *testing.T) {
	d := schema.ResourceData{}
	d.Set("id", "name")

	setValueOrID(&d, "id", "name", "54711781-274e-41b2-83c0-17194d0108f7")

	if d.Get("id").(string) != "name" {
		t.Fatal("err: 'id' does not match 'name'")
	}
}

func testSetIDOnResourceData(t *testing.T) {
	d := schema.ResourceData{}
	d.Set("id", "54711781-274e-41b2-83c0-17194d0108f7")

	setValueOrID(&d, "id", "name", "54711781-274e-41b2-83c0-17194d0108f7")

	if d.Get("id").(string) != "54711781-274e-41b2-83c0-17194d0108f7" {
		t.Fatal("err: 'id' doest not match '54711781-274e-41b2-83c0-17194d0108f7'")
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("COSMIC_API_URL"); v == "" {
		t.Fatal("COSMIC_API_URL must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_API_KEY"); v == "" {
		t.Fatal("COSMIC_API_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_SECRET_KEY"); v == "" {
		t.Fatal("COSMIC_SECRET_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_ZONE"); v == "" {
		t.Fatal("COSMIC_ZONE must be set for acceptance tests")
	}
}

// Name of a valid disk offering
var COSMIC_DISK_OFFERING_1 = os.Getenv("COSMIC_DISK_OFFERING_1")

// Name of a disk offering that COSMIC_DISK_OFFERING_1 can resize to
var COSMIC_DISK_OFFERING_2 = os.Getenv("COSMIC_DISK_OFFERING_2")

// Name of a valid service offering
var COSMIC_SERVICE_OFFERING_1 = os.Getenv("COSMIC_SERVICE_OFFERING_1")

// Name of a service offering that COSMIC_SERVICE_OFFERING_1 can resize to
var COSMIC_SERVICE_OFFERING_2 = os.Getenv("COSMIC_SERVICE_OFFERING_2")

// ID of an existing VPC
var COSMIC_VPC_ID = os.Getenv("COSMIC_VPC_ID")

// An available VPC offering
var COSMIC_VPC_OFFERING = os.Getenv("COSMIC_VPC_OFFERING")

// Name of an available network offering with forvpc=true
var COSMIC_VPC_NETWORK_OFFERING = os.Getenv("COSMIC_VPC_NETWORK_OFFERING")

// Name of a template that exists already for building VMs
var COSMIC_TEMPLATE = os.Getenv("COSMIC_TEMPLATE")

// Name of a project that exists already
var COSMIC_PROJECT_NAME = os.Getenv("COSMIC_PROJECT_NAME")

// Name of a zone that exists already
var COSMIC_ZONE = os.Getenv("COSMIC_ZONE")

// ID of the "default_acl" built-in ACL
var COSMIC_DEFAULT_ALLOW_ACL_ID = os.Getenv("COSMIC_DEFAULT_ALLOW_ACL_ID")
