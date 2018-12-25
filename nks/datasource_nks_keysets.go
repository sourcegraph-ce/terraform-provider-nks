package nks

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/StackPointCloud/nks-sdk-go/nks"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNKSKeyset() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNKSKeysetsRead,
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"category": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					filter := v.(string)
					if filter != "provider" && filter != "user_ssh" {
						errors = append(errors, fmt.Errorf("category can be either 'provider' or 'user_ssh'"))
					}
					return
				},
			},
			"entity": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceNKSKeysetsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	var name string
	var entity string
	var orgID int
	if temp, ok := d.GetOk("org_id"); ok {
		orgID = temp.(int)
	}
	category := d.Get("category").(string)
	if category == "provider" {
		if temp, ok := d.GetOk("entity"); ok {
			entity = temp.(string)
		} else {
			return fmt.Errorf("if category is set to 'provider' entity must be set")
		}
	}

	if tmp, ok := d.GetOk("name"); ok {
		name = tmp.(string)
	}

	// Fetch userprofile based on API token
	userProfile, err := config.Client.GetUserProfile()
	if orgID == 0 && len(userProfile) > 0 {
		for _, org := range userProfile[0].OrgMems {
			if org.IsDefault {
				orgID = org.Org.ID
				break
			}
		}
	} else if len(userProfile) == 0 {
		return fmt.Errorf("userprofile not found please check your credentials and the API endpoint")
	}

	keysets, err := config.Client.GetKeysets(orgID)
	if err != nil {
		return err
	}

	var userKeys []nks.Keyset
	var providerKeys []nks.Keyset

	for _, c := range keysets {
		if category == "provider" {
			providerKeys = append(providerKeys, c)
		} else if category == "user_ssh" {
			userKeys = append(userKeys, c)
		}
	}

	if len(providerKeys) > 0 {
		var subKeys []nks.Keyset
		for _, p := range providerKeys {
			if entity == p.Entity {
				subKeys = append(subKeys, p)
			}
		}
		if name != "" {
			var newKeys []nks.Keyset
			for _, p := range subKeys {
				if strings.Contains(strings.ToLower(p.Name), strings.ToLower(name)) {
					newKeys = append(newKeys, p)
				}
			}
			subKeys = newKeys
		}
		if len(subKeys) > 1 {
			return fmt.Errorf("there is more than one keyset in category '%s' and entity '%s' refine the search with 'name' parameter ", category, entity)
		} else if len(subKeys) == 0 {
			return fmt.Errorf("there are no keysets that match search criteria")
		}
		d.SetId(strconv.Itoa(subKeys[0].ID))
		return nil
	}
	if len(userKeys) > 0 {
		var subKeys []nks.Keyset
		if name != "" {
			for _, u := range userKeys {
				if strings.Contains(strings.ToLower(u.Name), strings.ToLower(name)) {
					subKeys = append(subKeys, u)
				}
			}
		} else {
			for _, u := range userKeys {
				if u.IsDefault {
					subKeys = append(subKeys, u)
					break
				}
			}
		}

		if len(subKeys) > 1 {
			return fmt.Errorf("there is more than one keyset in category '%s' refine the search with 'name' parameter ", category)
		} else if len(subKeys) == 0 {
			return fmt.Errorf("there are no keysets that match search criteria")
		}
		d.SetId(strconv.Itoa(subKeys[0].ID))
		return nil
	}
	return nil
}
