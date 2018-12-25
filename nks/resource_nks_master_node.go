package nks

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/StackPointCloud/nks-sdk-go/nks"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNKSMasterNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceNKSMasterNodeCreate,
		Read:   resourceNKSMasterNodeRead,
		Update: resourceNKSMasterNodeUpdate,
		Delete: resourceNKSMasterNodeDelete,
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
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
			},
			"provider_subnet_id_requested": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"provider_subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provider_subnet_cidr": {
				Type:     schema.TypeString,
				Optional: true,
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

func resourceNKSMasterNodeCreate(d *schema.ResourceData, meta interface{}) error {
	// Get client for API
	config := meta.(*Config)
	clusterID := d.Get("cluster_id").(int)
	orgID := d.Get("org_id").(int)

	// Set up new master node
	newNode := nks.NodeAdd{
		Count: 1,
		Role:  "master",
		Size:  d.Get("node_size").(string),
	}
	if d.Get("provider_code").(string) == "aws" {
		if _, ok := d.GetOk("zone"); !ok {
			return fmt.Errorf("NKS needs zone for AWS clusters.")
		}
		newNode.Zone = d.Get("zone").(string)
	}
	if d.Get("provider_code").(string) == "aws" || d.Get("provider_code").(string) == "azure" {
		// Allow user to submit values for provider_subnet_id_requested, and put real value in computed provider_subnet_id
		if _, ok := d.GetOk("provider_subnet_id_requested"); !ok {
			newNode.ProviderSubnetID = "__new__"
		} else {
			newNode.ProviderSubnetID = d.Get("provider_subnet_id_requested").(string)
		}
		if _, ok := d.GetOk("provider_subnet_cidr"); !ok {
			newNode.ProviderSubnetCidr = "10.0.1.0/24"
		} else {
			newNode.ProviderSubnetCidr = d.Get("provider_subnet_cidr").(string)
		}
	}
	log.Println("[DEBUG] MasterNode update attempting to add master node")
	nodes, err := config.Client.AddNode(orgID, clusterID, newNode)
	if err != nil {
		log.Printf("[DEBUG] MasterNode failed when creating new master node: %s\n", err)
		return err
	}
	timeout := int(d.Timeout("Create").Seconds())
	if v, ok := d.GetOk("timeout"); ok {
		timeout = v.(int)
	}
	if err := config.Client.WaitNodeProvisioned(orgID, clusterID, nodes[0].ID, timeout); err != nil {
		log.Printf("[DEBUG] MasterNode failed when waiting for new master node: %s\n", err)
		return err
	}

	// Set ID in TF
	d.SetId(strconv.Itoa(nodes[0].ID))

	return resourceNKSMasterNodeRead(d, meta)
}

func resourceNKSMasterNodeRead(d *schema.ResourceData, meta interface{}) error {
	nodeID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	config := meta.(*Config)
	clusterID := d.Get("cluster_id").(int)
	orgID := d.Get("org_id").(int)

	node, err := config.Client.GetNode(orgID, clusterID, nodeID)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			log.Println("[DEBUG] MasterNode read got a 404, delete")
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] MasterNode GetNode failed in read: %s\n", err)
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

	return nil
}

func resourceNKSMasterNodeUpdate(d *schema.ResourceData, meta interface{}) error {
	// No updates possible, everything requires rebuild and is set ForceNew
	return resourceNKSMasterNodeRead(d, meta)
}

func resourceNKSMasterNodeDelete(d *schema.ResourceData, meta interface{}) error {
	nodeID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	config := meta.(*Config)
	clusterID := d.Get("cluster_id").(int)
	orgID := d.Get("org_id").(int)

	if err = config.Client.DeleteNode(orgID, clusterID, nodeID); err != nil {
		log.Printf("[DEBUG] MasterNode DeleteNode failed: %s\n", err)
		return err
	}
	timeout := int(d.Timeout("Delete").Seconds())
	if v, ok := d.GetOk("timeout"); ok {
		timeout = v.(int)
	}
	if err = config.Client.WaitNodeDeleted(orgID, clusterID, nodeID, timeout); err != nil {
		log.Printf("[DEBUG] MasterNode WaitNodeDeleted failed when deleting node: %s\n", err)
		return err
	}
	log.Println("[DEBUG] Master node deletion complete")
	d.SetId("")
	return nil
}
