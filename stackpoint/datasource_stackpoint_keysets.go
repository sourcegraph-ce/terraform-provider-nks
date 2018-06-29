package stackpoint

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strconv"
)

func dataSourceStackPointKeysets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceStackPointKeysetsRead,
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"aws_keyset": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"azure_keyset": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"do_keyset": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"gce_keyset": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"gke_keyset": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"oneandone_keyset": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"packet_keyset": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"user_ssh_keyset": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceStackPointKeysetsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// Fetch userprofile based on API token
	up, err := config.Client.GetUserProfile()
	if err != nil {
		log.Println("[DEBUG] Keysets GetUserProfile failed: %s\n", err)
		return err
	}
	if up == nil {
		return fmt.Errorf("Could not fetch user profile, cannot proceed with keyset import from StackPointCloud\n")
	}

	// Use supplied org ID or fetch org ID from userprofile
	var orgID int
	if _, ok := d.GetOk("org_id"); !ok {
		orgID, err = config.Client.GetUserProfileDefaultOrg(&up[0])
		if err != nil {
			log.Println("[DEBUG] Keysets GetUserProfileDefaultOrg failed: %s\n", err)
			return err
		}
		d.Set("org_id", orgID)
	} else {
		orgID = d.Get("org_id").(int)
	}

	// Loop through keysets loaded into userprofile, store them into variables
	for _, ks := range []string{"aws", "azure", "do", "gce", "gke", "oneandone", "packet", "user_ssh"} {
		ksid, _ := config.Client.GetUserProfileKeysetID(&up[0], ks)
		d.Set(ks+"_keyset", ksid)
	}

	d.SetId(strconv.Itoa(orgID))

	return nil
}
