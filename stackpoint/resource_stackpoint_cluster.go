package stackpoint

import (
	"encoding/json"
	"fmt"
	"github.com/StackPointCloud/stackpoint-sdk-go/stackpointio"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strconv"
	"strings"
)

func resourceStackPointCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceStackPointClusterCreate,
		Read:   resourceStackPointClusterRead,
		Update: resourceStackPointClusterUpdate,
		Delete: resourceStackPointClusterDelete,
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"provider_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"provider_keyset": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"k8s_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"startup_master_size": {
				Type:     schema.TypeString,
				Required: true,
			},
			"startup_worker_count": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"startup_worker_size": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rbac_enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"dashboard_enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"etcd_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"platform": {
				Type:     schema.TypeString,
				Required: true,
			},
			"channel": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ssh_keyset": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"provider_resource_group": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"provider_network_id_requested": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"provider_network_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provider_network_cidr": {
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
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provider_keyset_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"notified": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"image": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kubeconfig_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dashboard_installed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceStackPointClusterCreate(d *schema.ResourceData, meta interface{}) error {
	// Get client for API
	config := meta.(*Config)
	orgID := d.Get("org_id").(int)
	sshKeyID := d.Get("ssh_keyset").(int)

	// Set up cluster structure based on input from user
	newCluster := stackpointio.Cluster{
		Name:              d.Get("cluster_name").(string),
		Provider:          d.Get("provider_code").(string),
		ProviderKey:       d.Get("provider_keyset").(int),
		MasterCount:       1,
		MasterSize:        d.Get("startup_master_size").(string),
		WorkerCount:       d.Get("startup_worker_count").(int),
		WorkerSize:        d.Get("startup_worker_size").(string),
		KubernetesVersion: d.Get("k8s_version").(string),
		RbacEnabled:       d.Get("rbac_enabled").(bool),
		DashboardEnabled:  d.Get("dashboard_enabled").(bool),
		EtcdType:          d.Get("etcd_type").(string),
		Platform:          d.Get("platform").(string),
		Channel:           d.Get("channel").(string),
		SSHKeySet:         sshKeyID,
		Solutions:         []stackpointio.Solution{}, // helm_tiller will get automatically installed
	}
	// Grab provider-specific fields
	if d.Get("provider_code").(string) == "aws" {
		if _, ok := d.GetOk("region"); !ok {
			return fmt.Errorf("StackPoint needs region for AWS clusters.")
		}
		if _, ok := d.GetOk("zone"); !ok {
			return fmt.Errorf("StackPoint needs zone for AWS clusters.")
		}
		// Allow user to submit values for provider_network_id_requested, and put real value in computed provider_network_id
		if _, ok := d.GetOk("provider_network_id_requested"); !ok {
			newCluster.ProviderNetworkID = "__new__"
		} else {
			newCluster.ProviderNetworkID = d.Get("provider_network_id_requested").(string)
		}
		if _, ok := d.GetOk("provider_network_cidr"); !ok {
			newCluster.ProviderNetworkCdr = "10.0.0.0/16"
		} else {
			newCluster.ProviderNetworkCdr = d.Get("provider_network_cidr").(string)
		}
		// Allow user to submit values for provider_subnet_id_requested, and put real value in computed provider_subnet_id
		if _, ok := d.GetOk("provider_subnet_id_requested"); !ok {
			newCluster.ProviderSubnetID = "__new__"
		} else {
			newCluster.ProviderSubnetID = d.Get("provider_subnet_id_requested").(string)
		}
		if _, ok := d.GetOk("provider_subnet_cidr"); !ok {
			newCluster.ProviderSubnetCidr = "10.0.0.0/24"
		} else {
			newCluster.ProviderSubnetCidr = d.Get("provider_subnet_cidr").(string)
		}
		newCluster.Region = d.Get("region").(string)
		newCluster.Zone = d.Get("zone").(string)
	} else if d.Get("provider_code").(string) == "do" || d.Get("provider_code").(string) == "gce" ||
		d.Get("provider_code").(string) == "gke" || d.Get("provider_code").(string) == "oneandone" {
		if _, ok := d.GetOk("region"); !ok {
			return fmt.Errorf("StackPoint needs region for DigitalOcean/GCE/GKE clusters.")
		}
		newCluster.Region = d.Get("region").(string)
	} else if d.Get("provider_code").(string) == "azure" {
		if _, ok := d.GetOk("provider_resource_group"); !ok {
			return fmt.Errorf("StackPoint needs provider_resource_group for Azure clusters.")
		}
		if _, ok := d.GetOk("region"); !ok {
			return fmt.Errorf("StackPoint needs region for Azure clusters.")
		}
		// Allow user to submit values for provider_network_id_requested, and put real value in computed provider_network_id
		if _, ok := d.GetOk("provider_network_id_requested"); !ok {
			newCluster.ProviderNetworkID = "__new__"
		} else {
			newCluster.ProviderNetworkID = d.Get("provider_network_id_requested").(string)
		}
		if _, ok := d.GetOk("provider_network_cidr"); !ok {
			newCluster.ProviderNetworkCdr = "10.0.0.0/16"
		} else {
			newCluster.ProviderNetworkCdr = d.Get("provider_network_cidr").(string)
		}
		// Allow user to submit values for provider_subnet_id_requested, and put real value in computed provider_subnet_id
		if _, ok := d.GetOk("provider_subnet_id_requested"); !ok {
			newCluster.ProviderSubnetID = "__new__"
		} else {
			newCluster.ProviderSubnetID = d.Get("provider_subnet_id_requested").(string)
		}
		if _, ok := d.GetOk("provider_subnet_cidr"); !ok {
			newCluster.ProviderSubnetCidr = "10.0.0.0/24"
		} else {
			newCluster.ProviderSubnetCidr = d.Get("provider_subnet_cidr").(string)
		}
		newCluster.ProviderResourceGp = d.Get("provider_resource_group").(string)
		newCluster.Region = d.Get("region").(string)
	} else if d.Get("provider_code").(string) == "packet" {
		if _, ok := d.GetOk("region"); !ok {
			return fmt.Errorf("StackPoint needs region for Packet clusters.")
		}
		if _, ok := d.GetOk("project_id"); !ok {
			return fmt.Errorf("StackPoint needs project_id for Packet clusters.")
		}
		newCluster.Region = d.Get("region").(string)
		newCluster.ProjectID = d.Get("project_id").(string)
	}
	// Do cluster creation call
	cluster, err := config.Client.CreateCluster(orgID, newCluster)

	reqJSON, _ := json.Marshal(newCluster)
	resJSON, _ := json.Marshal(cluster)

	log.Println("[DEBUG] Cluster create request", string(reqJSON))
	log.Println("[DEBUG] Cluster create response", string(resJSON))

	// Don't bail until request and response are logged above
	if err != nil {
		log.Printf("[DEBUG] Cluster error at CreateCluster: %s", err)
		return err
	}

	// Wait until provisioned
	timeout := int(d.Timeout("Create").Seconds())
	if v, ok := d.GetOk("timeout"); ok {
		timeout = v.(int)
	}
	err = config.Client.WaitClusterProvisioned(orgID, cluster.ID, timeout)
	if err != nil {
		log.Printf("[DEBUG] Cluster error at WaitClusterProvisioned: %s", err)
		return err
	}
	// Set ID in TF
	d.SetId(strconv.Itoa(cluster.ID))

	return resourceStackPointClusterRead(d, meta)
}

func resourceStackPointClusterRead(d *schema.ResourceData, meta interface{}) error {
	clusterID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	config := meta.(*Config)
	orgID := d.Get("org_id").(int)

	cluster, err := config.Client.GetCluster(orgID, clusterID)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			log.Println("[DEBUG] Cluster read got a 404")
			d.SetId("")
			return nil
		}
		return err
	}
	d.Set("state", cluster.State)
	d.Set("instanceID", cluster.InstanceID)
	d.Set("cluster_name", cluster.Name)
	d.Set("provider_name", cluster.Provider)
	d.Set("provider_keyset", cluster.ProviderKey)
	d.Set("provider_keyset_name", cluster.ProviderKeyName)
	d.Set("region", cluster.Region)
	d.Set("zone", cluster.Zone)
	d.Set("project_id", cluster.ProjectID)
	d.Set("provider_resource_group", cluster.ProviderResourceGp)
	d.Set("provider_network_id", cluster.ProviderNetworkID)
	d.Set("provider_network_cidr", cluster.ProviderNetworkCdr)
	d.Set("provider_subnet_id", cluster.ProviderSubnetID)
	d.Set("provider_subnet_cidr", cluster.ProviderSubnetCidr)
	d.Set("owner", cluster.Owner)
	d.Set("notified", cluster.Notified)
	d.Set("k8s_version", cluster.KubernetesVersion)
	d.Set("dashboard_enabled", cluster.DashboardEnabled)
	d.Set("dashboard_installed", cluster.DashboardInstalled)
	d.Set("kubeconfig_path", cluster.KubeconfigPath)
	d.Set("rbac_enabled", cluster.RbacEnabled)
	d.Set("master_count", cluster.MasterCount)
	d.Set("master_size", cluster.MasterSize)
	d.Set("etcd_type", cluster.EtcdType)
	d.Set("platform", cluster.Platform)
	d.Set("image", cluster.Image)
	d.Set("channel", cluster.Channel)

	return nil
}

func resourceStackPointClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	//clusterID, err := strconv.Atoi(d.Id())
	//if err != nil {
	//	return err
	//}
	//config := meta.(*Config)
	//orgID := d.Get("org_id").(int)

	return resourceStackPointClusterRead(d, meta)
}

func resourceStackPointClusterDelete(d *schema.ResourceData, meta interface{}) error {
	clusterID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	config := meta.(*Config)
	orgID := d.Get("org_id").(int)

	if err = config.Client.DeleteCluster(orgID, clusterID); err != nil {
		log.Printf("[DEBUG] Cluster error while calling DeleteCluster: %s", err)
		return err
	}
	timeout := int(d.Timeout("Delete").Seconds())
	if v, ok := d.GetOk("timeout"); ok {
		timeout = v.(int)
	}
	if err = config.Client.WaitClusterDeleted(orgID, clusterID, timeout); err != nil {
		log.Printf("[DEBUG] Cluster error at WaitClusterDeleted: %s", err)
		return err
	}
	log.Println("[DEBUG] Cluster deletion complete")
	d.SetId("")
	return nil
}
