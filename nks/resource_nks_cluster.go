package stackpoint

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/StackPointCloud/stackpoint-sdk-go/stackpointio"
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
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region": {
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
			"kubeconfig_path": {
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
			"aws": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"do", "packet"},
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:     schema.TypeString,
							Required: true,
						},
						"zone": {
							Type:     schema.TypeString,
							Required: true,
						},
						"provider_network_id_requested": {
							Type:     schema.TypeString,
							Default:  "__new__",
							Optional: true,
						},
						"provider_network_cidr": {
							Type:     schema.TypeString,
							Default:  "10.0.0.0/16",
							Optional: true,
						},
						"provider_subnet_cidr": {
							Type:     schema.TypeString,
							Default:  "10.0.0.0/16",
							Optional: true,
						},
						"provider_subnet_id_requested": {
							Type:     schema.TypeString,
							Default:  "__new__",
							Optional: true,
						},

						"provider_resource_group": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"provider_network_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"provider_subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"provider_resource_group_requested": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"do": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"aws", "packet"},
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"packet": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"aws", "do"},
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:     schema.TypeString,
							Required: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Required: true,
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

	// Set up cluster structure based on input from user
	newCluster := stackpointio.Cluster{
		Name: d.Get("cluster_name").(string),
		// Provider:          d.Get("provider_code").(string),
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

	if _, ok := d.GetOk("aws"); ok {
		newCluster.Region = d.Get("aws.0.region").(string)

		newCluster.Zone = d.Get("aws.0.zone").(string)

		// Allow user to submit values for provider_network_id_requested, and put real value in computed provider_network_id
		if temp, ok := d.GetOk("aws.0.provider_network_id_requested"); !ok {
			newCluster.ProviderNetworkID = "__new__"
		} else {
			newCluster.ProviderNetworkID = temp.(string)
		}
		if temp, ok := d.GetOk("aws.0.provider_network_cidr"); !ok {
			newCluster.ProviderNetworkCdr = "10.0.0.0/16"
		} else {
			newCluster.ProviderNetworkCdr = temp.(string)
		}

		// Allow user to submit values for provider_subnet_id_requested, and put real value in computed provider_subnet_id
		if temp, ok := d.GetOk("aws.0.provider_subnet_id_requested"); !ok {
			newCluster.ProviderSubnetID = "__new__"
		} else {
			newCluster.ProviderSubnetID = temp.(string)
		}
		if temp, ok := d.GetOk("aws.0.provider_subnet_cidr"); !ok {
			newCluster.ProviderSubnetCidr = "10.0.0.0/24"
		} else {
			newCluster.ProviderSubnetCidr = temp.(string)
		}
	}
	if _, ok := d.GetOk("do"); ok {
		newCluster.Region = d.Get("do.0.region").(string)
		newCluster.Provider = "do"
	}
	if _, ok := d.GetOk("gce"); ok {
		newCluster.Region = d.Get("gce.0.region").(string)
		newCluster.Provider = "gce"
	}
	if _, ok := d.GetOk("gke"); ok {
		newCluster.Region = d.Get("gke.0.region").(string)
		newCluster.Provider = "gke"
	}
	if _, ok := d.GetOk("oneandone"); ok {
		newCluster.Region = d.Get("oneandone.0.region").(string)
		newCluster.Provider = "oneandone"
	}
	if _, ok := d.GetOk("azure"); ok {
		newCluster.Provider = "azure"
		if temp, ok := d.GetOk("azure.0.provider_resource_group_requested"); !ok {
			newCluster.ProviderResourceGp = "__new__"
		} else {
			newCluster.ProviderResourceGp = temp.(string)
		}
		newCluster.Region = d.Get("azure.0.region").(string)

		// Allow user to submit values for provider_network_id_requested, and put real value in computed provider_network_id
		if temp, ok := d.GetOk("azure.0.provider_network_id_requested"); !ok {
			newCluster.ProviderNetworkID = "__new__"
		} else {
			newCluster.ProviderNetworkID = temp.(string)
		}
		if temp, ok := d.GetOk("azure.0.provider_network_cidr"); !ok {
			newCluster.ProviderNetworkCdr = "10.0.0.0/16"
		} else {
			newCluster.ProviderNetworkCdr = temp.(string)
		}
		// Allow user to submit values for provider_subnet_id_requested, and put real value in computed provider_subnet_id
		if temp, ok := d.GetOk("azure.0.provider_subnet_id_requested"); !ok {
			newCluster.ProviderSubnetID = "__new__"
		} else {
			newCluster.ProviderSubnetID = temp.(string)
		}
		if temp, ok := d.GetOk("azure.0.provider_subnet_cidr"); !ok {
			newCluster.ProviderSubnetCidr = "10.0.0.0/24"
		} else {
			newCluster.ProviderSubnetCidr = temp.(string)
		}
	}

	if _, ok := d.GetOk("packet"); ok {
		newCluster.Provider = "packet"

		newCluster.Region = d.Get("packet.0.region").(string)
		newCluster.ProjectID = d.Get("packet.0.project_id").(string)
	}

	// Do cluster creation call
	cluster, err := config.Client.CreateCluster(orgID, newCluster)

	// Don't bail until request and response are logged above
	if err != nil {
		reqJSON, _ := json.Marshal(newCluster)
		log.Printf("[DEBUG] Cluster error at CreateCluster: Error: %s", err)
		log.Printf("[DEBUG] Request: %s", string(reqJSON))

		return err
	}

	// Wait until provisioned
	timeout := int(d.Timeout("Create").Seconds())
	if v, ok := d.GetOk("timeout"); ok {
		timeout = v.(int)
	}
	if err = config.Client.WaitClusterRunning(orgID, cluster.ID, false, timeout); err != nil {
		log.Printf("[DEBUG] Cluster error at WaitClusterProvisioned: %s", err)
		return err
	}
	// Set ID in TF
	d.SetId(strconv.Itoa(cluster.ID))

	return nil //resourceNKSClusterRead(d, meta)
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
	d.Set("k8s_version_upgrades", cluster.KubernetesMigrationVersions)

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
		if err = config.Client.WaitClusterRunning(orgID, clusterID, false, timeout); err != nil {
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
