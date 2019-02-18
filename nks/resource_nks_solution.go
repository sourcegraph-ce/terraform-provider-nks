package nks

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/NetApp/nks-sdk-go/nks"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNKSSolution() *schema.Resource {
	return &schema.Resource{
		Create: resourceNKSSolutionCreate,
		Read:   resourceNKSSolutionRead,
		Update: resourceNKSSolutionUpdate,
		Delete: resourceNKSSolutionDelete,
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"cluster_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"solution": {
				Type:     schema.TypeString,
				Required: true,
			},
			"config": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deleteable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"keyset": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceNKSSolutionCreate(d *schema.ResourceData, meta interface{}) error {
	// Get client for API
	config := meta.(*Config)
	clusterID := d.Get("cluster_id").(int)
	orgID := d.Get("org_id").(int)

	// If config file is sent, use that to try to install solution
	log.Printf("[DEBUG] Solution create attempting to add solution: %s\n", d.Get("name").(string))
	var (
		solution *nks.Solution
		err      error
	)
	if c, ok := d.GetOk("config"); ok {
		solutionJSON := c.(string)
		if solutionJSON == "" {
			return fmt.Errorf("no config data sent for solution %s", d.Get("name").(string))
		}
		solution, err = config.Client.AddSolutionFromJSON(orgID, clusterID, solutionJSON)
		if err != nil {
			log.Printf("[DEBUG] Solution %s create failed to add solution from JSON: %s\n", d.Get("name").(string), err)
			return err
		}
	} else {
		// No config file sent, try to install solution simply by name
		newSolution := nks.Solution{
			Solution: d.Get("solution").(string),
			State:    "draft",
		}
		if temp, ok := d.GetOk("keyset"); ok {
			newSolution.Keyset = temp.(int)
		}

		solution, err = config.Client.AddSolution(orgID, clusterID, newSolution)
		if err != nil {
			log.Printf("[DEBUG] Solution %s create failed to add solution from name: %s\n", d.Get("name").(string), err)
			return err
		}
	}
	timeout := int(d.Timeout("Create").Seconds())
	if v, ok := d.GetOk("timeout"); ok {
		timeout = v.(int)
	}
	if err := config.Client.WaitSolutionInstalled(orgID, clusterID, solution.ID, timeout); err != nil {
		log.Printf("[DEBUG] Solution %s create failed while waiting for installed: %s\n", d.Get("name").(string), err)
		return err
	}

	// Set ID in TF
	d.SetId(strconv.Itoa(solution.ID))

	return resourceNKSSolutionRead(d, meta)
}

func resourceNKSSolutionRead(d *schema.ResourceData, meta interface{}) error {
	solutionID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	config := meta.(*Config)
	clusterID := d.Get("cluster_id").(int)
	orgID := d.Get("org_id").(int)

	solution, err := config.Client.GetSolution(orgID, clusterID, solutionID)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			log.Println("[DEBUG] Solution read got a 404, delete")
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] Solution %s read failed for solution: %s\n", d.Get("name").(string), err)
		return err
	}
	d.Set("state", solution.State)
	d.Set("name", solution.Name)
	d.Set("deleteable", solution.Deleteable)

	return nil
}

func resourceNKSSolutionUpdate(d *schema.ResourceData, meta interface{}) error {
	// No updates possible
	return resourceNKSSolutionRead(d, meta)
}

func resourceNKSSolutionDelete(d *schema.ResourceData, meta interface{}) error {
	solutionID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	// Check if solution is deleteable
	if !d.Get("deleteable").(bool) {
		//return fmt.Errorf("Solution %s is not deleteable\n", d.Get("name").(string))
		log.Printf("[DEBUG] Solution %s is not deleteable, but skipping and deleting ID\n", d.Get("name").(string))
		d.SetId("")
		return nil
	}

	config := meta.(*Config)
	clusterID := d.Get("cluster_id").(int)
	orgID := d.Get("org_id").(int)

	if err = config.Client.DeleteSolution(orgID, clusterID, solutionID); err != nil {
		log.Printf("[DEBUG] Solution %s delete failed for solution: %s\n", d.Get("name").(string), err)
		return err
	}
	timeout := int(d.Timeout("Delete").Seconds())
	if v, ok := d.GetOk("timeout"); ok {
		timeout = v.(int)
	}
	if err = config.Client.WaitSolutionDeleted(orgID, clusterID, solutionID, timeout); err != nil {
		log.Printf("[DEBUG] Solution %s delete failed while waiting for solution to delete: %s\n", d.Get("name").(string), err)
		return err
	}
	log.Printf("[DEBUG] Solution %s deletion complete\n", d.Get("name").(string))
	d.SetId("")
	return nil
}
