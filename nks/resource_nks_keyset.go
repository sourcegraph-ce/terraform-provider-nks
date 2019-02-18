package nks

import (
	"fmt"
	"strconv"

	"github.com/NetApp/nks-sdk-go/nks"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNKSKeyset() *schema.Resource {
	return &schema.Resource{
		Create: resourceNKSKeysetCreate,
		Read:   resourceNKSKeysetRead,
		Delete: resourceNKSKeysetDelete,
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"category": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					if v.(string) != "provider" && v.(string) != "user_ssh" && v.(string) != "solution" {
						errors = append(errors, fmt.Errorf("Category can be one of following 'provider', 'user_ssh' or 'solution'"))
					}
					return
				},
			},
			"entity": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"keys": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"workspaces": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNKSKeysetCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	orgID := d.Get("org_id").(int)
	name := d.Get("name").(string)
	category := d.Get("category").(string)
	req := nks.Keyset{
		Org:        orgID,
		Name:       name,
		Category:   category,
		Workspaces: []int{},
	}

	if temp, ok := d.GetOk("entity"); ok {
		if category == "user_ssh" {
			return fmt.Errorf("when 'category' is set to '%s', 'entity' cannot be set", category)
		}
		req.Entity = temp.(string)
	}

	rawKeys := d.Get("keys").([]interface{})
	req.Keys = make([]nks.Key, len(rawKeys))
	for i, v := range rawKeys {
		value := v.(map[string]interface{})
		req.Keys[i] = nks.Key{
			Type:  value["key_type"].(string),
			Value: value["key"].(string),
		}
	}

	rawWorkspaces := d.Get("workspaces").([]interface{})
	req.Workspaces = make([]int, len(rawWorkspaces))
	for i, v := range rawWorkspaces {
		req.Workspaces[i] = v.(int)
	}

	keyset, err := config.Client.CreateKeyset(orgID, req)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(keyset.ID))

	return resourceNKSKeysetRead(d, meta)
}

func resourceNKSKeysetRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	orgID := d.Get("org_id").(int)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	keyset, err := config.Client.GetKeyset(orgID, id)
	if err != nil {
		return err
	}

	d.Set("name", keyset.Name)
	d.Set("category", keyset.Category)
	d.Set("entity", keyset.Entity)

	workspaces := make([]interface{}, len(keyset.Workspaces))
	for i, w := range keyset.Workspaces {
		workspaces[i] = w
	}
	d.Set("workspaces", workspaces)

	return nil
}

func resourceNKSKeysetDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	orgID := d.Get("org_id").(int)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	return config.Client.DeleteKeyset(orgID, id)
}
