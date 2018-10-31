package nks

import (
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNKSOrganization() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNKSOrganizationRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceNKSOrganizationRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	var name string
	if tmp, ok := d.GetOk("name"); ok {
		name = tmp.(string)
	}

	if name == "" {
		userProfile, err := config.Client.GetUserProfile()
		if err != nil {
			return err
		}
		if len(userProfile) > 0 {
			for _, org := range userProfile[0].OrgMems {
				if org.IsDefault {
					d.SetId(strconv.Itoa(org.Org.ID))
					return nil
				}
			}
		}
	}

	organizations, err := config.Client.GetOrganizations()
	if err != nil {
		return err
	}

	for _, org := range organizations {
		if strings.Contains(strings.ToLower(org.Name), strings.ToLower(name)) {
			d.SetId(strconv.Itoa(org.ID))
			return nil
		}
	}

	return nil
}
