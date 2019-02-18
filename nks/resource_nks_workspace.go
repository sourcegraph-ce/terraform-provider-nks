package nks

import (
	"strconv"

	"github.com/NetApp/nks-sdk-go/nks"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNKSWorkspace() *schema.Resource {
	return &schema.Resource{
		Create: resourceNKSWorkspaceCreate,
		Read:   resourceNKSWorkspaceRead,
		Delete: resourceNKSWorkspaceDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"org_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNKSWorkspaceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	newWorkspace := nks.Workspace{
		Org:            d.Get("org_id").(int),
		Name:           d.Get("name").(string),
		TeamWorkspaces: []nks.TeamWorkspace{},
	}

	if temp, ok := d.GetOk("default"); ok {
		newWorkspace.IsDefault = temp.(bool)
	}

	workspace, err := config.Client.CreateWorkspace(newWorkspace.Org, newWorkspace)
	if err != nil {
		return err
	}
	d.SetId(strconv.Itoa(workspace.ID))

	return resourceNKSWorkspaceRead(d, meta)
}

func resourceNKSWorkspaceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	orgID := d.Get("org_id").(int)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	workspaces, err := config.Client.GetWorkspaces(orgID)
	if err != nil {
		return err
	}

	var workspace nks.Workspace
	for _, w := range workspaces {
		if w.ID == id {
			workspace = w
			break
		}
	}

	d.Set("name", workspace.Name)
	d.Set("default", workspace.IsDefault)

	return nil
}

func resourceNKSWorkspaceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	orgID := d.Get("org_id").(int)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	return config.Client.DeleteWorkspace(orgID, id)
}
