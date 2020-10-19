package cosmic

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCosmicSSHKeyPair() *schema.Resource {
	return &schema.Resource{
		Create: resourceCosmicSSHKeyPairCreate,
		Read:   resourceCosmicSSHKeyPairRead,
		Delete: resourceCosmicSSHKeyPairDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"public_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"private_key": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"fingerprint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCosmicSSHKeyPairCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	name := d.Get("name").(string)
	publicKey := d.Get("public_key").(string)

	if publicKey != "" {
		// Register supplied key
		p := client.SSH.NewRegisterSSHKeyPairParams(name, publicKey)

		_, err := client.SSH.RegisterSSHKeyPair(p)
		if err != nil {
			return err
		}
	} else {
		// No key supplied, must create one and return the private key
		p := client.SSH.NewCreateSSHKeyPairParams(name)

		r, err := client.SSH.CreateSSHKeyPair(p)
		if err != nil {
			return err
		}
		d.Set("private_key", r.Privatekey)
	}

	log.Printf("[DEBUG] Key pair successfully generated at Cosmic")
	d.SetId(name)

	return resourceCosmicSSHKeyPairRead(d, meta)
}

func resourceCosmicSSHKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	log.Printf("[DEBUG] looking for key pair with name %s", d.Id())

	p := client.SSH.NewListSSHKeyPairsParams()
	p.SetName(d.Id())

	r, err := client.SSH.ListSSHKeyPairs(p)
	if err != nil {
		return err
	}
	if r.Count == 0 {
		log.Printf("[DEBUG] Key pair %s does not exist", d.Id())
		d.SetId("")
		return nil
	}

	//SSHKeyPair name is unique in a cosmic account so dont need to check for multiple
	d.Set("name", r.SSHKeyPairs[0].Name)
	d.Set("fingerprint", r.SSHKeyPairs[0].Fingerprint)

	return nil
}

func resourceCosmicSSHKeyPairDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CosmicClient)

	// Create a new parameter struct
	p := client.SSH.NewDeleteSSHKeyPairParams(d.Id())

	// Remove the SSH Keypair
	_, err := client.SSH.DeleteSSHKeyPair(p)
	if err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if strings.Contains(err.Error(), fmt.Sprintf(
			"A key pair with name '%s' does not exist for account", d.Id())) {
			return nil
		}

		return fmt.Errorf("Error deleting key pair: %s", err)
	}

	return nil
}
