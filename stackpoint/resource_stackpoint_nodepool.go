package stackpoint

import (
	"encoding/json"
	"github.com/StackPointCloud/stackpoint-sdk-go/stackpointio"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strconv"
	"time"
)

func resourceStackPointNodePool() *schema.Resource {
	return &schema.Resource{
		Create: resourceStackPointNodePoolCreate,
		Read:   resourceStackPointNodePoolRead,
		Update: resourceStackPointNodePoolUpdate,
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
                        "cluster_id": {                         
                                Type:     schema.TypeInt,   
                                Required: true,             
                        },
			"pool_name": {
				Type:     schema.TypeString,
				Required: true,
			},
                        "number_nodes": {                     
                                Type:     schema.TypeInt,   
                                Required: true,                                   
                        },
			"worker_size": {
				Type:     schema.TypeString,
				Required: true,
			},
			"platform": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceStackPointClusterCreate(d *schema.ResourceData, meta interface{}) error {
        newNodePool := stackpointio.NodePool {
		Name: d.Get("pool_name").(string),
                NodeCount: d.Get("number_nodes").(int),
                Size:      d.Get("worker_size").(string),
                Platform:  d.Get("platform").(string),
	}
        // Create new nodepool
	client := meta.(*stackpointio.APIClient)
        pool, err := client.CreateNodePool(d.Get("org_id").(int), d.Get("cluster_id").(int), newNodePool)

	reqJSON, _ := json.Marshal(newNodePool)
	resJSON, _ := json.Marshal(pool)

	log.Println("[DEBUG] NodePool create request", string(reqJSON))
	log.Println("[DEBUG] NodePool create response", string(resJSON))

	// Don't bail until request and response are logged above
	if err != nil {
		return err
	}
	// Wait until provisioned (until "state" is "running")
	for i := 1; ; i++ {
		isActive, err := client.IsNodePoolActive(d.Get("org_id").(int), d.Get("cluster_id").(int), pool.ID)
		if err != nil {
			return err
		}
		if isActive {
			d.SetId(strconv.Itoa(pool.ID))
			break
		}
		time.Sleep(time.Second)
	}
	return resourceStackPointNodePoolRead(d, meta)
}

func resourceStackPointNodePoolRead(d *schema.ResourceData, meta interface{}) error {
	nodePoolID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	client := meta.(*stackpointio.APIClient)
	pool, err := client.GetNodePool(d.Get("org_id").(int), nodePoolID)
	if err != nil {
		return err
	}
	d.Set("state", cluster.State)
	d.Set("instanceID", cluster.InstanceID)
	return nil
}

func resourceStackPointNodePoolUpdate(d *schema.ResourceData, meta interface{}) error {
//////
        newNode := spio.NodeAddToPool{Count: nodeCount,
                Role:       "worker",
                NodePoolID: nodepoolID}

        nodes, err := client.AddNodesToNodePool(orgID, clusterID, nodepoolID, newNode)

        return nil                                             
}

func resourceStackPointNodePoolDelete(d *schema.ResourceData, meta interface{}) error {
	// Can't delete nodepools, but unset and forget, I guess?? (seems like a bad idea)
	d.SetId("")
	return nil
}
