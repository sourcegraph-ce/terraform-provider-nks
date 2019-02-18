package nks

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/NetApp/nks-sdk-go/nks"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNKSCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceNKSClusterCreate,
		Read:   resourceNKSClusterRead,
		Update: resourceNKSClusterUpdate,
		Delete: resourceNKSClusterDelete,
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
				Optional: true,
			},
			"startup_worker_count": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"startup_worker_size": {
				Type:     schema.TypeString,
				Required: true,
			},
			"startup_worker_min_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"startup_worker_max_count": {
				Type:     schema.TypeInt,
				Optional: true,
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
			"provider_resource_group_requested": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"provider_resource_group": {
				Type:     schema.TypeString,
				Computed: true,
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
			"kubeconfig": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dashboard_installed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"k8s_version_upgrades": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"public_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"network_component": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"component_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"zone": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"provider_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceNKSClusterCreate(d *schema.ResourceData, meta interface{}) error {
	// Get client for API
	config := meta.(*Config)
	orgID := d.Get("org_id").(int)
	sshKeyID := d.Get("ssh_keyset").(int)
	providerCode := d.Get("provider_code").(string)

	// Set up cluster structure based on input from user
	newCluster := nks.Cluster{
		Name:              d.Get("cluster_name").(string),
		Provider:          d.Get("provider_code").(string),
		ProviderKey:       d.Get("provider_keyset").(int),
		WorkerCount:       d.Get("startup_worker_count").(int),
		WorkerSize:        d.Get("startup_worker_size").(string),
		KubernetesVersion: d.Get("k8s_version").(string),
		RbacEnabled:       d.Get("rbac_enabled").(bool),
		DashboardEnabled:  d.Get("dashboard_enabled").(bool),
		EtcdType:          d.Get("etcd_type").(string),
		Platform:          d.Get("platform").(string),
		Channel:           d.Get("channel").(string),
		SSHKeySet:         sshKeyID,
		Solutions:         []nks.Solution{}, // helm_tiller will get automatically installed
		NetworkComponents: []nks.NetworkComponent{},
	}

	if providerCode == "aws" || providerCode == "azure" || providerCode == "gce" || providerCode == "gke" {
		newCluster.MasterCount = 1
		newCluster.MasterSize = d.Get("startup_master_size").(string)
	}

	if providerCode == "eks" {
		if _, ok := d.GetOk("startup_worker_min_count"); !ok {
			return fmt.Errorf("NKS needs min node number")
		}
		newCluster.MinNodeCount = d.Get("startup_worker_min_count").(int)

		if _, ok := d.GetOk("startup_worker_max_count"); !ok {
			return fmt.Errorf("NKS needs max node number")
		}
		newCluster.MaxNodeCount = d.Get("startup_worker_max_count").(int)

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
	}

	// Grab provider-specific fields
	if d.Get("provider_code").(string) == "aws" {
		if _, ok := d.GetOk("region"); !ok {
			return fmt.Errorf("NKS needs region for AWS clusters.")
		}
		if _, ok := d.GetOk("zone"); !ok {
			return fmt.Errorf("NKS needs zone for AWS clusters.")
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
			return fmt.Errorf("NKS needs region for DigitalOcean/GCE/GKE clusters.")
		}
		newCluster.Region = d.Get("region").(string)
	} else if d.Get("provider_code").(string) == "azure" {
		// Allow user to submit values for provider_resource_group_requested, and put real value in computed provider_resource_group
		if _, ok := d.GetOk("provider_resource_group_requested"); !ok {
			newCluster.ProviderResourceGp = "__new__"
		} else {
			newCluster.ProviderResourceGp = d.Get("provider_resource_group_requested").(string)
		}
		if _, ok := d.GetOk("region"); !ok {
			return fmt.Errorf("NKS needs region for Azure clusters.")
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
	} else if d.Get("provider_code").(string) == "packet" {
		if _, ok := d.GetOk("region"); !ok {
			return fmt.Errorf("NKS needs region for Packet clusters.")
		}
		if _, ok := d.GetOk("project_id"); !ok {
			return fmt.Errorf("NKS needs project_id for Packet clusters.")
		}
		newCluster.Region = d.Get("region").(string)
		newCluster.ProjectID = d.Get("project_id").(string)
	}

	if _, ok := d.GetOk("region"); !ok {
		return fmt.Errorf("NKS needs region for clusters.")
	}
	newCluster.Region = d.Get("region").(string)

	//Network Components
	if vRaw, ok := d.GetOk("network_component"); ok {
		componentRaw := vRaw.(*schema.Set).List()
		for _, raw := range componentRaw {
			rawMap := raw.(map[string]interface{})
			if rawMap["id"] == nil || rawMap["cidr"] == nil || rawMap["component_type"] == nil || rawMap["provider_id"] == nil || rawMap["vpc_id"] == nil || rawMap["zone"] == nil {
				return fmt.Errorf("Required fields for network component are id, cidr, component_type, provider_id, vpc_id, zone")
			}
			netComponent := nks.NetworkComponent{
				ID:            rawMap["id"].(string),
				Cidr:          rawMap["cidr"].(string),
				ComponentType: rawMap["component_type"].(string),
				ProviderID:    rawMap["provider_id"].(string),
				VpcID:         rawMap["vpc_id"].(string),
				Zone:          rawMap["zone"].(string),
			}
			if rawMap["name"] != nil {
				netComponent.Name = rawMap["name"].(string)
			}
			newCluster.NetworkComponents = append(newCluster.NetworkComponents, netComponent)
		}
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
	if err = config.Client.WaitClusterRunning(orgID, cluster.ID, true, timeout); err != nil {
		log.Printf("[DEBUG] Cluster error at WaitClusterProvisioned: %s", err)
		return err
	}
	// Set ID in TF
	d.SetId(strconv.Itoa(cluster.ID))

	return resourceNKSClusterRead(d, meta)
}

func resourceNKSClusterRead(d *schema.ResourceData, meta interface{}) error {
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
	nodes, err := config.Client.GetNodes(orgID, clusterID)
	if err != nil {
		return err
	}

	rawNodes := make([]map[string]interface{}, len(nodes))

	for i, n := range nodes {
		rawNode := map[string]interface{}{
			"instance_id": n.InstanceID,
			"public_ip":   n.PublicIP,
			"private_ip":  n.PrivateIP,
		}

		rawNodes[i] = rawNode
	}

	if err := d.Set("nodes", rawNodes); err != nil {
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
	d.Set("rbac_enabled", cluster.RbacEnabled)
	d.Set("master_count", cluster.MasterCount)
	d.Set("master_size", cluster.MasterSize)
	d.Set("etcd_type", cluster.EtcdType)
	d.Set("platform", cluster.Platform)
	d.Set("image", cluster.Image)
	d.Set("channel", cluster.Channel)
	d.Set("k8s_version_upgrades", cluster.KubernetesMigrationVersions)

	kubeconfig, err := config.Client.GetKubeConfig(orgID, clusterID)
	if err != nil {
		return err
	}
	d.Set("kubeconfig", kubeconfig)

	return nil
}

func resourceNKSClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	clusterID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	config := meta.(*Config)
	orgID := d.Get("org_id").(int)

	if d.HasChange("k8s_version") {
		oldV, newV := d.GetChange("k8s_version")
		log.Printf("[DEBUG] Cluster has a change in k8s_version, old value %s, new value %s\n", oldV, newV)
		cluster, err := config.Client.GetCluster(orgID, clusterID)
		if err != nil {
			log.Printf("[DEBUG] Cluster in change in k8s_version, could not fetch cluster info: %s\n", err)
			return err
		}
		err = config.Client.UpgradeClusterToVersion(*cluster, newV.(string))
		if err != nil {
			log.Printf("[DEBUG] Cluster in change in k8s_version, failed to upgrade cluster: %s\n", err)
			return err
		}
		timeout := int(d.Timeout("Update").Seconds())
		if v, ok := d.GetOk("timeout"); ok {
			timeout = v.(int)
		}
		if err = config.Client.WaitClusterRunning(orgID, clusterID, true, timeout); err != nil {
			log.Printf("[DEBUG] Cluster error at WaitClusterDeleted: %s", err)
			return err
		}
		log.Println("[DEBUG] Cluster successfully upgraded k8s_version")
	}

	return resourceNKSClusterRead(d, meta)
}

func resourceNKSClusterDelete(d *schema.ResourceData, meta interface{}) error {
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
