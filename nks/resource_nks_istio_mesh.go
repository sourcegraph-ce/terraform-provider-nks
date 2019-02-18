package nks

import (
	"fmt"
	"strconv"

	"github.com/NetApp/nks-sdk-go/nks"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNKSIstioMesh() *schema.Resource {
	return &schema.Resource{
		Create: resourceNKSIstioMeshCreate,
		Read:   resourceNKSIstioMeshRead,
		Delete: resourceNKSIstioMeshDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"mesh_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					filter := v.(string)
					if filter != "cross_cluster" {
						errors = append(errors, fmt.Errorf("'mesh_type' can only be 'cross_cluster'"))
					}
					return
				},
			},
			"workspace": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"org_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"members": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MinItems: 2,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"role": {
							Type:     schema.TypeString,
							Required: true,
						},
						"istio_solution_id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceNKSIstioMeshCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	orgID := d.Get("org_id").(int)
	workspace := d.Get("workspace").(int)

	im := nks.IstioMeshRequest{
		Name:      d.Get("name").(string),
		MeshType:  d.Get("mesh_type").(string),
		Workspace: workspace,
		Members:   []nks.MemberRequest{},
	}

	membersRaw := d.Get("members").([]interface{})
	im.Members = make([]nks.MemberRequest, len(membersRaw))
	solutionIds := make([]int, len(membersRaw))
	for i, v := range membersRaw {
		value := v.(map[string]interface{})
		solutionIds[i] = value["istio_solution_id"].(int)
		im.Members[i] = nks.MemberRequest{
			Cluster: value["cluster"].(int),
			Role:    value["role"].(string),
		}
	}

	timeout := int(d.Timeout("Create").Seconds())
	if v, ok := d.GetOk("timeout"); ok {
		timeout = v.(int)
	}

	for i, solutionID := range solutionIds {
		if err := config.Client.WaitSolutionInstalled(orgID, im.Members[i].Cluster, solutionID, timeout); err != nil {
			return fmt.Errorf("Solution %d create failed while waiting for 'installed' state: %s", solutionID, err)
		}
	}

	istioMesh, err := config.Client.CreateIstioMesh(orgID, workspace, im)
	if err != nil {
		return err
	}

	if err := config.Client.WaitIstioMeshCreated(orgID, workspace, istioMesh.ID, timeout); err != nil {
		return fmt.Errorf("Istio mesh %s create failed while waiting: %s", d.Get("name").(string), err)
	}

	d.SetId(strconv.Itoa(istioMesh.ID))

	return resourceNKSIstioMeshRead(d, meta)
}

func resourceNKSIstioMeshRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	orgID := d.Get("org_id").(int)
	workspace := d.Get("workspace").(int)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	var mesh nks.IstioMesh

	meshes, err := config.Client.GetIstioMeshes(orgID, workspace)
	for _, m := range meshes {
		if m.ID == id {
			mesh = m
		}
	}

	d.Set("name", mesh.Name)
	d.Set("mesh_type", mesh.MeshType)

	return nil
}

func resourceNKSIstioMeshDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	orgID := d.Get("org_id").(int)
	workspace := d.Get("workspace").(int)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	err = config.Client.DeleteIstioMesh(orgID, workspace, id)
	if err != nil {
		return err
	}
	return nil
}
