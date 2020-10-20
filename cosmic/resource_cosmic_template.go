package cosmic

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCosmicTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceCosmicTemplateCreate,
		Read:   resourceCosmicTemplateRead,
		Update: resourceCosmicTemplateUpdate,
		Delete: resourceCosmicTemplateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"display_text": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"format": {
				Type:     schema.TypeString,
				Required: true,
			},

			"hypervisor": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"os_type": {
				Type:     schema.TypeString,
				Required: true,
			},

			"url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"is_dynamically_scalable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"is_extractable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"is_featured": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"is_public": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"password_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"is_ready": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"is_ready_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300,
			},

			"zone": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				ForceNew:   true,
				Deprecated: deprecatedZoneMsg(),
			},
		},
	}
}

func resourceCosmicTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	if err := verifyTemplateParams(d); err != nil {
		return err
	}

	name := d.Get("name").(string)

	// Compute/set the display text
	displaytext := d.Get("display_text").(string)
	if displaytext == "" {
		displaytext = name
	}

	// Retrieve the os_type ID
	ostypeid, e := retrieveID(client, "os_type", d.Get("os_type").(string))
	if e != nil {
		return e.Error()
	}

	// Retrieve the zone ID
	zoneid, e := retrieveID(client, "zone", client.ZoneName)
	if e != nil {
		return e.Error()
	}

	// Create a new parameter struct
	p := client.Template.NewRegisterTemplateParams(
		displaytext,
		d.Get("format").(string),
		d.Get("hypervisor").(string),
		name,
		ostypeid,
		d.Get("url").(string),
		zoneid)

	// Set optional parameters
	if v, ok := d.GetOk("is_dynamically_scalable"); ok {
		p.SetIsdynamicallyscalable(v.(bool))
	}

	if v, ok := d.GetOk("is_extractable"); ok {
		p.SetIsextractable(v.(bool))
	}

	if v, ok := d.GetOk("is_featured"); ok {
		p.SetIsfeatured(v.(bool))
	}

	if v, ok := d.GetOk("is_public"); ok {
		p.SetIspublic(v.(bool))
	}

	if v, ok := d.GetOk("password_enabled"); ok {
		p.SetPasswordenabled(v.(bool))
	}

	// Create the new template
	r, err := client.Template.RegisterTemplate(p)
	if err != nil {
		return fmt.Errorf("Error creating template %s: %s", name, err)
	}

	d.SetId(r.RegisterTemplate[0].Id)

	// Wait until the template is ready to use, or timeout with an error...
	currentTime := time.Now().Unix()
	timeout := int64(d.Get("is_ready_timeout").(int))
	for {
		// Start with the sleep so the register action has a few seconds
		// to process the registration correctly. Without this wait
		time.Sleep(10 * time.Second)

		err := resourceCosmicTemplateRead(d, meta)
		if err != nil {
			return err
		}

		if d.Get("is_ready").(bool) {
			return nil
		}

		if time.Now().Unix()-currentTime > timeout {
			return fmt.Errorf("Timeout while waiting for template to become ready")
		}
	}
}

func resourceCosmicTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	// Get the template details
	t, count, err := client.Template.GetTemplateByID(d.Id(), "executable")
	if err != nil {
		if count == 0 {
			log.Printf(
				"[DEBUG] Template %s no longer exists", d.Get("name").(string))
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("name", t.Name)
	d.Set("display_text", t.Displaytext)
	d.Set("format", t.Format)
	d.Set("hypervisor", t.Hypervisor)
	d.Set("is_dynamically_scalable", t.Isdynamicallyscalable)
	d.Set("is_extractable", t.Isextractable)
	d.Set("is_featured", t.Isfeatured)
	d.Set("is_public", t.Ispublic)
	d.Set("password_enabled", t.Passwordenabled)
	d.Set("is_ready", t.Isready)

	setValueOrID(d, "os_type", t.Ostypename, t.Ostypeid)
	setValueOrID(d, "zone", t.Zonename, t.Zoneid)

	return nil
}

func resourceCosmicTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)
	name := d.Get("name").(string)

	// Create a new parameter struct
	p := client.Template.NewUpdateTemplateParams(d.Id())

	if d.HasChange("name") {
		p.SetName(name)
	}

	if d.HasChange("display_text") {
		p.SetDisplaytext(d.Get("display_text").(string))
	}

	if d.HasChange("format") {
		p.SetFormat(d.Get("format").(string))
	}

	if d.HasChange("is_dynamically_scalable") {
		p.SetIsdynamicallyscalable(d.Get("is_dynamically_scalable").(bool))
	}

	if d.HasChange("os_type") {
		ostypeid, e := retrieveID(client, "os_type", d.Get("os_type").(string))
		if e != nil {
			return e.Error()
		}
		p.SetOstypeid(ostypeid)
	}

	if d.HasChange("password_enabled") {
		p.SetPasswordenabled(d.Get("password_enabled").(bool))
	}

	_, err := client.Template.UpdateTemplate(p)
	if err != nil {
		return fmt.Errorf("Error updating template %s: %s", name, err)
	}

	return resourceCosmicTemplateRead(d, meta)
}

func resourceCosmicTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	// Create a new parameter struct
	p := client.Template.NewDeleteTemplateParams(d.Id())

	// Delete the template
	log.Printf("[INFO] Deleting template: %s", d.Get("name").(string))
	_, err := client.Template.DeleteTemplate(p)
	if err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", d.Id())) {
			return nil
		}

		return fmt.Errorf("Error deleting template %s: %s", d.Get("name").(string), err)
	}
	return nil
}

func verifyTemplateParams(d *schema.ResourceData) error {
	format := d.Get("format").(string)
	if format != "OVA" && format != "QCOW2" && format != "RAW" && format != "VHD" && format != "VMDK" {
		return fmt.Errorf(
			"%s is not a valid format. Valid options are 'OVA','QCOW2', 'RAW', 'VHD' and 'VMDK'", format)
	}

	return nil
}
