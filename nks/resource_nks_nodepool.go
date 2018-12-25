package nks

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/StackPointCloud/nks-sdk-go/nks"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNKSNodePool() *schema.Resource {
	return &schema.Resource{
		Create: resourceNKSNodePoolCreate,
		Read:   resourceNKSNodePoolRead,
		Update: resourceNKSNodePoolUpdate,
		Delete: resourceNKSNodePoolDelete,
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
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
			"provider_subnet_id_requested": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"provider_subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
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

func resourceNKSNodePoolCreate(d *schema.ResourceData, meta interface{}) error {
	// Get client for API
	config := meta.(*Config)
	clusterID := d.Get("cluster_id").(int)
	orgID := d.Get("org_id").(int)

	newNodepool := nks.NodePool{
		Name:      "TerraForm NodePool",
		NodeCount: d.Get("worker_count").(int),
		Size:      d.Get("worker_size").(string),
		Platform:  d.Get("platform").(string),
	}
	if d.Get("provider_code").(string) == "aws" {
		if _, ok := d.GetOk("zone"); !ok {
			return fmt.Errorf("NKS needs zone for AWS clusters.")
		}
		newNodepool.Zone = d.Get("zone").(string)
	}
	if d.Get("provider_code").(string) == "aws" || d.Get("provider_code").(string) == "azure" {
		// Allow user to submit values for provider_subnet_id_requested, and put real value in computed provider_subnet_id
		if _, ok := d.GetOk("provider_subnet_id_requested"); !ok {
			newNodepool.ProviderSubnetID = "__new__"
		} else {
			newNodepool.ProviderSubnetID = d.Get("provider_subnet_id_requested").(string)
		}
		if _, ok := d.GetOk("provider_subnet_cidr"); !ok {
			newNodepool.ProviderSubnetCidr = "10.0.1.0/24"
		} else {
			newNodepool.ProviderSubnetCidr = d.Get("provider_subnet_cidr").(string)
		}
	}
	log.Println("[DEBUG] Nodepool creation running")
	pool, err := config.Client.CreateNodePool(orgID, clusterID, newNodepool)
	if err != nil {
		log.Printf("[DEBUG] Nodepool creation failed: %s\n", err)
		return err
	}
	timeout := int(d.Timeout("Create").Seconds())
	if v, ok := d.GetOk("timeout"); ok {
		timeout = v.(int)
	}
	config.Client.WaitNodePoolProvisioned(orgID, clusterID, pool.ID, timeout)

	// Set ID in TF
	d.SetId(strconv.Itoa(pool.ID))

	// Nodepools now change to active status before nodes inside them are provisioned,
	// so wait for new nodes to provision
	nodes, err := config.Client.GetNodesInPool(orgID, clusterID, pool.ID)
	if err != nil {
		log.Printf("[DEBUG] Nodepool GetNodesInPool failed: %s\n", err)
		return err
	}
	for i := 0; i < len(nodes); i++ {
		if err = config.Client.WaitNodeProvisioned(orgID, clusterID, nodes[i].ID, timeout); err != nil {
			log.Printf("[DEBUG] Nodepool WaitNodeProvisioned failed: %s\n", err)
			return err
		}
	}

	// Set ID in TF
	d.SetId(strconv.Itoa(pool.ID))

	return resourceNKSNodePoolRead(d, meta)
}

func resourceNKSNodePoolRead(d *schema.ResourceData, meta interface{}) error {
	nodepoolID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	config := meta.(*Config)
	clusterID := d.Get("cluster_id").(int)
	orgID := d.Get("org_id").(int)

	nodepool, err := config.Client.GetNodePool(orgID, clusterID, nodepoolID)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			log.Println("[DEBUG] Nodepool read got a 404, delete")
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] Nodepool GetNodePool failed: %s\n", err)
		return err
	}
	d.Set("state", nodepool.State)
	d.Set("platform", nodepool.Platform)
	d.Set("provider_subnet_id", nodepool.ProviderSubnetID)
	d.Set("provider_subnet_cidr", nodepool.ProviderSubnetCidr)
	d.Set("worker_size", nodepool.Size)
	d.Set("worker_count", nodepool.NodeCount)
	d.Set("cluster_id", nodepool.ClusterID)
	d.Set("name", nodepool.Name)
	d.Set("autoscaled", nodepool.Autoscaled)
	d.Set("autoscale_min_count", nodepool.MinCount)
	d.Set("autoscale_max_count", nodepool.MaxCount)
	d.Set("instance_id", nodepool.InstanceID)

	return nil
}

func resourceNKSNodePoolUpdate(d *schema.ResourceData, meta interface{}) error {
	nodepoolID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	config := meta.(*Config)
	clusterID := d.Get("cluster_id").(int)
	orgID := d.Get("org_id").(int)

	if d.HasChange("worker_count") {
		oldV, newV := d.GetChange("worker_count")
		oldVi, newVi := oldV.(int), newV.(int)

		if oldVi > newVi {
			// Decrease worker count, try to cull the herd to match wanted value
			nodes, err := config.Client.GetNodesInPool(orgID, clusterID, nodepoolID)
			if err != nil {
				log.Printf("[DEBUG] Nodepool GetNodesInPool failed: %s\n", err)
				return err
			}
			workerCount := len(nodes)

			// Delete only as many workers (from workerCount) as it takes to get down to wanted value (newVi)
			for i := 0; i < (workerCount - newVi); i++ {
				if err = config.Client.DeleteNode(orgID, clusterID, nodes[i].ID); err != nil {
					log.Printf("[DEBUG] Nodepool DeleteNode failed: %s\n", err)
					return err
				}
				timeout := int(d.Timeout("Delete").Seconds())
				if v, ok := d.GetOk("timeout"); ok {
					timeout = v.(int)
				}
				if err = config.Client.WaitNodeDeleted(orgID, clusterID, nodes[i].ID, timeout); err != nil {
					log.Printf("[DEBUG] Nodepool WaitNodeDeleted failed: %s\n", err)
					return err
				}
			}
		} else {
			// Increase worker count
			newNode := nks.NodeAddToPool{
				Count:      (newVi - oldVi),
				Role:       "worker",
				NodePoolID: nodepoolID,
			}
			nodes, err := config.Client.AddNodesToNodePool(orgID, clusterID, nodepoolID, newNode)
			if err != nil {
				log.Printf("[DEBUG] Nodepool AddNodesToNodePool failed: %s\n", err)
				return err
			}
			for i := 0; i < len(nodes); i++ {
				timeout := int(d.Timeout("Create").Seconds())
				if v, ok := d.GetOk("timeout"); ok {
					timeout = v.(int)
				}
				if err = config.Client.WaitNodeProvisioned(orgID, clusterID, nodes[i].ID, timeout); err != nil {
					log.Printf("[DEBUG] Nodepool WaitNodeProvisioned failed: %s\n", err)
					return err
				}
			}
		}
	}
	return resourceNKSNodePoolRead(d, meta)
}

func resourceNKSNodePoolDelete(d *schema.ResourceData, meta interface{}) error {
	nodepoolID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	config := meta.(*Config)
	clusterID := d.Get("cluster_id").(int)
	orgID := d.Get("org_id").(int)

	// Gather list of nodes, delete them all before calling for nodepool deletion
	nodes, err := config.Client.GetNodesInPool(orgID, clusterID, nodepoolID)
	if err != nil {
		log.Printf("[DEBUG] Nodepool GetNodesInPool failed in deletion: %s\n", err)
		return err
	}
	for i := 0; i < len(nodes); i++ {
		// Delete node if active
		if nodes[i].State == nks.NodeRunningStateString {
			if err = config.Client.DeleteNode(orgID, clusterID, nodes[i].ID); err != nil {
				log.Printf("[DEBUG] Nodepool DeleteNode failed: %s\n", err)
				return err
			}
			timeout := int(d.Timeout("Delete").Seconds())
			if v, ok := d.GetOk("timeout"); ok {
				timeout = v.(int)
			}
			if err = config.Client.WaitNodeDeleted(orgID, clusterID, nodes[i].ID, timeout); err != nil {
				log.Printf("[DEBUG] Nodepool WaitNodeDeleted failed: %s\n", err)
				return err
			}
		}
	}
	// Delete the actual nodepool
	if err = config.Client.DeleteNodePool(orgID, clusterID, nodepoolID); err != nil {
		log.Printf("[DEBUG] Nodepool DeleteNodePool failed: %s\n", err)
		return err
	}
	log.Printf("[DEBUG] NodePool deletion completed")
	d.SetId("")
	return nil
}
