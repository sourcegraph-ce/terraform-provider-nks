package nks

import (
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNKSWorkspace() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNKSWorkspaceRead,
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceNKSWorkspaceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	var name string
	if tmp, ok := d.GetOk("name"); ok {
		name = tmp.(string)
	}
	if len(name) == 0 {
		name = "Default"
	}

	var orgID int
	if tmp, ok := d.GetOk("org_id"); ok {
		orgID = tmp.(int)
	}

	workspaces, err := config.Client.GetWorkspaces(orgID)
	if err != nil {
		return err
	}

	for _, workspace := range workspaces {
		if strings.Contains(strings.ToLower(workspace.Name), strings.ToLower(name)) {
			d.SetId(strconv.Itoa(workspace.ID))
			return nil
		}
	}

	return nil
}
