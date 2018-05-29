package stackpoint

import (
	"fmt"
	"github.com/StackPointCloud/stackpoint-sdk-go/stackpointio"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strconv"
	"strings"
)

func resourceStackPointMasterNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceStackPointMasterNodeCreate,
		Read:   resourceStackPointMasterNodeRead,
		Update: resourceStackPointMasterNodeUpdate,
		Delete: resourceStackPointMasterNodeDelete,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"node_size": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
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

func resourceStackPointMasterNodeCreate(d *schema.ResourceData, meta interface{}) error {
	// Get client for API
	config := meta.(*Config)
	clusterID := d.Get("cluster_id").(int)

	// Set up new master node
	newNode := stackpointio.NodeAdd{
		Count: 1,
		Role:  "master",
		Size:  d.Get("node_size").(string),
	}
	if d.Get("provider_code").(string) == "aws" {
		if _, ok := d.GetOk("zone"); !ok {
			return fmt.Errorf("StackPoint needs zone for AWS clusters.")
		}
		newNode.Zone = d.Get("zone").(string)
	}
	if d.Get("provider_code").(string) == "aws" || d.Get("provider_code").(string) == "azure" {
		if _, ok := d.GetOk("provider_subnet_id"); !ok {
			return fmt.Errorf("StackPoint needs provider_subnet_id for AWS and Azure clusters.")
		}
		if _, ok := d.GetOk("provider_subnet_cidr"); !ok {
			return fmt.Errorf("StackPoint needs provider_subnet_cidr for AWS and Azure clusters.")
		}
		newNode.ProviderSubnetID = d.Get("provider_subnet_id").(string)
		newNode.ProviderSubnetCidr = d.Get("provider_subnet_cidr").(string)
	}
	log.Println("[DEBUG] Cluster update attempting to add master node\n")
	nodes, err := config.Client.AddNode(config.OrgID, clusterID, newNode)
	if err != nil {
		return err
	}
	timeout := int(d.Timeout("Create").Seconds())
	if v, ok := d.GetOk("timeout"); ok {
		timeout = v.(int)
	}
	if err := config.Client.WaitNodeProvisioned(config.OrgID, clusterID, nodes[0].ID, timeout); err != nil {
		return err
	}

	// Set ID in TF
	d.SetId(strconv.Itoa(nodes[0].ID))

	return resourceStackPointMasterNodeRead(d, meta)
}

func resourceStackPointMasterNodeRead(d *schema.ResourceData, meta interface{}) error {
	nodeID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	config := meta.(*Config)
	clusterID := d.Get("cluster_id").(int)
	node, err := config.Client.GetNode(config.OrgID, clusterID, nodeID)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			log.Println("[DEBUG] Master node read got a 404, delete")
			d.SetId("")
			return nil
		}
		return err
	}
	d.Set("state", node.State)
	d.Set("node_size", node.Size)
	d.Set("platform", node.Platform)
	d.Set("provider_subnet_id", node.ProviderSubnetID)
	d.Set("provider_subnet_cidr", node.ProviderSubnetCidr)
	d.Set("location", node.Location)
	d.Set("private_ip", node.PrivateIP)
	d.Set("public_ip", node.PublicIP)
	d.Set("cluster_id", node.ClusterID)
	d.Set("instance_id", node.InstanceID)
	outDebug(fmt.Sprintf("In MasterNodeRead, instance_id: %s\n", node.InstanceID))

	return nil
}

func resourceStackPointMasterNodeUpdate(d *schema.ResourceData, meta interface{}) error {
	// No updates possible, everything requires rebuild and is set ForceNew
	return resourceStackPointMasterNodeRead(d, meta)
}

func resourceStackPointMasterNodeDelete(d *schema.ResourceData, meta interface{}) error {
	nodeID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	config := meta.(*Config)
	clusterID := d.Get("cluster_id").(int)
	if err = config.Client.DeleteNode(config.OrgID, clusterID, nodeID); err != nil {
		return err
	}
	timeout := int(d.Timeout("Delete").Seconds())
	if v, ok := d.GetOk("timeout"); ok {
		timeout = v.(int)
	}
	if err = config.Client.WaitNodeDeleted(config.OrgID, clusterID, nodeID, timeout); err != nil {
		return err
	}
	log.Println("[DEBUG] Master node deletion complete")
	d.SetId("")
	return nil
}
