package stackpoint

import (
	"fmt"
	"github.com/StackPointCloud/stackpoint-sdk-go/stackpointio"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strconv"
	"strings"
)

func resourceStackPointNodePool() *schema.Resource {
	return &schema.Resource{
		Create: resourceStackPointNodePoolCreate,
		Read:   resourceStackPointNodePoolRead,
		Update: resourceStackPointNodePoolUpdate,
		Delete: resourceStackPointNodePoolDelete,
		Schema: map[string]*schema.Schema{
			"platform": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"provider_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"worker_size": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"worker_count": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"cluster_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"provider_subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"provider_subnet_cidr": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"autoscaled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"autoscale_min_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"autoscale_max_count": {
				Type:     schema.TypeInt,
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
		},
	}
}

func resourceStackPointNodePoolCreate(d *schema.ResourceData, meta interface{}) error {
	// Get client for API
	config := meta.(*Config)
	clusterID := d.Get("cluster_id").(int)

	outDebug(fmt.Sprintf("In nodepool creation\n"))
	newNodepool := stackpointio.NodePool{
		Name:      "TerraForm NodePool",
		NodeCount: d.Get("worker_count").(int),
		Size:      d.Get("worker_size").(string),
		Platform:  d.Get("platform").(string),
	}
	if d.Get("provider_code").(string) == "aws" {
		if _, ok := d.GetOk("zone"); !ok {
			return fmt.Errorf("StackPoint needs zone for AWS clusters.")
		}
		newNodepool.Zone = d.Get("zone").(string)
	}
	if d.Get("provider_code").(string) == "aws" || d.Get("provider_code").(string) == "azure" {
		if _, ok := d.GetOk("provider_subnet_id"); !ok {
			return fmt.Errorf("StackPoint needs provider_subnet_id for AWS and Azure clusters.")
		}
		if _, ok := d.GetOk("provider_subnet_cidr"); !ok {
			return fmt.Errorf("StackPoint needs provider_subnet_cidr for AWS and Azure clusters.")
		}
		newNodepool.ProviderSubnetID = d.Get("provider_subnet_id").(string)
		newNodepool.ProviderSubnetCidr = d.Get("provider_subnet_cidr").(string)
	}
	log.Println("[DEBUG] Nodepool creation running\n")
	pool, err := config.Client.CreateNodePool(config.OrgID, clusterID, newNodepool)
	if err != nil {
		return err
	}
	timeout := int(d.Timeout("Create").Seconds())
	if v, ok := d.GetOk("timeout"); ok {
		timeout = v.(int)
	}
	config.Client.WaitNodePoolProvisioned(config.OrgID, clusterID, pool.ID, timeout)

	// Set ID in TF
	d.SetId(strconv.Itoa(pool.ID))

	// Nodepools now change to active status before nodes inside them are provisioned,
	// so wait for new nodes to provision
	nodes, err := config.Client.GetNodesInPool(config.OrgID, clusterID, pool.ID)
	if err != nil {
		return err
	}
	for i := 0; i < len(nodes); i++ {
		outDebug(fmt.Sprintf("In nodepoolcreate, waiting for node.ID to provision: %d\n", nodes[i].ID))
		if err = config.Client.WaitNodeProvisioned(config.OrgID, clusterID, nodes[i].ID, timeout); err != nil {
			return err
		}
	}

	// Set ID in TF
	d.SetId(strconv.Itoa(pool.ID))

	return resourceStackPointNodePoolRead(d, meta)
}

func resourceStackPointNodePoolRead(d *schema.ResourceData, meta interface{}) error {
	nodepoolID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	config := meta.(*Config)
	clusterID := d.Get("cluster_id").(int)
	nodepool, err := config.Client.GetNodePool(config.OrgID, clusterID, nodepoolID)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			log.Println("[DEBUG] Nodepool read got a 404, delete")
			d.SetId("")
			outDebug(fmt.Sprintf("In nodepool read, nodepool is gone, got 404\n"))
			return nil
		}
		outDebug(fmt.Sprintf("In nodepool creation\n"))
		return err
	}
	d.Set("state", nodepool.State)
	d.Set("platform", nodepool.Platform)
	d.Set("worker_size", nodepool.Size)
	d.Set("worker_count", nodepool.NodeCount)
	d.Set("cluster_id", nodepool.ClusterID)
	d.Set("name", nodepool.Name)
	d.Set("autoscaled", nodepool.Autoscaled)
	d.Set("autoscale_min_count", nodepool.MinCount)
	d.Set("autoscale_max_count", nodepool.MaxCount)
	d.Set("instance_id", nodepool.InstanceID)
	outDebug(fmt.Sprintf("In NodePoolRead, instance_id: %s\n", nodepool.InstanceID))

	return nil
}

func resourceStackPointNodePoolUpdate(d *schema.ResourceData, meta interface{}) error {
	nodepoolID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	config := meta.(*Config)
	clusterID := d.Get("cluster_id").(int)

	if d.HasChange("worker_count") {
		outDebug(fmt.Sprintf("In nodepoolupdate, worker_count has change, worker_count: %d, nodepooidID: %d\n", d.Get("worker_count").(int), nodepoolID))
		oldV, newV := d.GetChange("worker_count")
		oldVi, newVi := oldV.(int), newV.(int)

		if oldVi > newVi {
			// Decrease worker count, try to cull the herd to match wanted value
			nodes, err := config.Client.GetNodesInPool(config.OrgID, clusterID, nodepoolID)
			if err != nil {
				return err
			}
			workerCount := len(nodes)

			// Delete only as many workers (from workerCount) as it takes to get down to wanted value (newVi)
			for i := 0; i < (workerCount - newVi); i++ {
				if err = config.Client.DeleteNode(config.OrgID, clusterID, nodes[i].ID); err != nil {
					return err
				}
				timeout := int(d.Timeout("Delete").Seconds())
				if v, ok := d.GetOk("timeout"); ok {
					timeout = v.(int)
				}
				if err = config.Client.WaitNodeDeleted(config.OrgID, clusterID, nodes[i].ID, timeout); err != nil {
					return err
				}
			}
		} else {
			// Increase worker count
			newNode := stackpointio.NodeAddToPool{
				Count:      (newVi - oldVi),
				Role:       "worker",
				NodePoolID: nodepoolID,
			}
			nodes, err := config.Client.AddNodesToNodePool(config.OrgID, clusterID, nodepoolID, newNode)
			if err != nil {
				return err
			}
			for i := 0; i < len(nodes); i++ {
				timeout := int(d.Timeout("Create").Seconds())
				if v, ok := d.GetOk("timeout"); ok {
					timeout = v.(int)
				}
				if err = config.Client.WaitNodeProvisioned(config.OrgID, clusterID, nodes[i].ID, timeout); err != nil {
					return err
				}
			}
		}
	}
	return resourceStackPointNodePoolRead(d, meta)
}

func resourceStackPointNodePoolDelete(d *schema.ResourceData, meta interface{}) error {
	nodepoolID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	config := meta.(*Config)
	clusterID := d.Get("cluster_id").(int)
	outDebug(fmt.Sprintf("In nodepooldelete, ID is: %d\n", nodepoolID))

	// Nodepools can't actually be deleted, but we can delete all the worker nodes within it
	nodes, err := config.Client.GetNodesInPool(config.OrgID, clusterID, nodepoolID)
	if err != nil {
		return err
	}
	for i := 0; i < len(nodes); i++ {
		// Delete node if active
		if nodes[i].State == stackpointio.NodeRunningStateString {
			outDebug(fmt.Sprintf("In nodepooldelete, deleting node.ID: %d\n", nodes[i].ID))
			if err = config.Client.DeleteNode(config.OrgID, clusterID, nodes[i].ID); err != nil {
				return err
			}
			timeout := int(d.Timeout("Delete").Seconds())
			if v, ok := d.GetOk("timeout"); ok {
				timeout = v.(int)
			}
			if err = config.Client.WaitNodeDeleted(config.OrgID, clusterID, nodes[i].ID, timeout); err != nil {
				return err
			}
		}
	}
	log.Println("[DEBUG] NodePool deletion completed")
	d.SetId("")
	return nil
}
