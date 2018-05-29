package stackpoint

import (
	"encoding/json"
	"fmt"
	"github.com/StackPointCloud/stackpoint-sdk-go/stackpointio"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"os"
	"strconv"
	"strings"
)

func outDebug(m string) {
	f, err := os.OpenFile("/tmp/tf_debug.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fmt.Fprintf(f, "%s", m)
}
func resourceStackPointCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceStackPointClusterCreate,
		Read:   resourceStackPointClusterRead,
		Update: resourceStackPointClusterUpdate,
		Delete: resourceStackPointClusterDelete,
		Schema: map[string]*schema.Schema{
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
			"provider_network_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"provider_network_cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"provider_subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
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
			"solutions": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
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
		SSHKeySet:         config.SSHKeyset,
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
		if _, ok := d.GetOk("provider_network_id"); !ok {
			newCluster.ProviderNetworkID = "__new__"
		} else {
			newCluster.ProviderNetworkID = d.Get("provider_network_id").(string)
		}
		if _, ok := d.GetOk("provider_network_cidr"); !ok {
			return fmt.Errorf("StackPoint needs provider_network_cidr for AWS clusters.")
		}
		if _, ok := d.GetOk("provider_subnet_id"); !ok {
			newCluster.ProviderSubnetID = "__new__"
		} else {
			newCluster.ProviderSubnetID = d.Get("provider_subnet_id").(string)
		}
		if _, ok := d.GetOk("provider_subnet_cidr"); !ok {
			return fmt.Errorf("StackPoint needs provider_subnet_cidr for AWS clusters.")
		}
		newCluster.Region = d.Get("region").(string)
		newCluster.Zone = d.Get("zone").(string)
		newCluster.ProviderNetworkCdr = d.Get("provider_network_cidr").(string)
		newCluster.ProviderSubnetCidr = d.Get("provider_subnet_cidr").(string)
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
		if _, ok := d.GetOk("provider_network_id"); !ok {
			newCluster.ProviderNetworkID = "__new__"
		} else {
			newCluster.ProviderNetworkID = d.Get("provider_network_id").(string)
		}
		if _, ok := d.GetOk("provider_network_cidr"); !ok {
			return fmt.Errorf("StackPoint needs provider_network_cidr for Azure clusters.")
		}
		if _, ok := d.GetOk("provider_subnet_id"); !ok {
			newCluster.ProviderSubnetID = "__new__"
		} else {
			newCluster.ProviderSubnetID = d.Get("provider_subnet_id").(string)
		}
		if _, ok := d.GetOk("provider_subnet_cidr"); !ok {
			return fmt.Errorf("StackPoint needs provider_subnet_cidr for Azure clusters.")
		}
		newCluster.ProviderResourceGp = d.Get("provider_resource_group").(string)
		newCluster.Region = d.Get("region").(string)
		newCluster.ProviderNetworkCdr = d.Get("provider_network_cidr").(string)
		newCluster.ProviderSubnetCidr = d.Get("provider_subnet_cidr").(string)
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
	cluster, err := config.Client.CreateCluster(config.OrgID, newCluster)

	reqJSON, _ := json.Marshal(newCluster)
	resJSON, _ := json.Marshal(cluster)

	log.Println("[DEBUG] Cluster create request", string(reqJSON))
	log.Println("[DEBUG] Cluster create response", string(resJSON))

	// Don't bail until request and response are logged above
	if err != nil {
		return err
	}

	// Wait until provisioned
	timeout := int(d.Timeout("Create").Seconds())
	if v, ok := d.GetOk("timeout"); ok {
		timeout = v.(int)
	}
	err = config.Client.WaitClusterProvisioned(config.OrgID, cluster.ID, timeout)
	if err != nil {
		log.Println("[DEBUG] Error while waiting for cluster to be provisioned")
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
	cluster, err := config.Client.GetCluster(config.OrgID, clusterID)
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

	// Collect solutions for TF state
	var solArray []string
	for _, sol := range cluster.Solutions {
		solArray = append(solArray, sol.Solution)
	}
	d.Set("solutions", solArray)

	return nil
}

func resourceStackPointClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	clusterID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	config := meta.(*Config)

	if d.HasChange("solutions") {
		_, newV := d.GetChange("solutions")
		userIntList := newV.([]interface{})
		var userSolutionList []string
		for _, item := range userIntList {
			userSolutionList = append(userSolutionList, item.(string))
		}
		solutions, err := config.Client.GetSolutions(config.OrgID, clusterID)
		if err != nil {
			return err
		}
		// Loop through currently configured solutions, delete any that aren't in user's list
		var configuredSols []string
		for _, sol := range solutions {
			if !stackpointio.StringInSlice(sol.Solution, userSolutionList) {
				// Solution not in user's list, needs to be deleted
				if !sol.Deleteable {
					return fmt.Errorf("Solution, %s, is marked as non-deleteable and cannot be deleted.", sol.Solution)
				}
				if err := config.Client.DeleteSolution(config.OrgID, clusterID, sol.ID); err != nil {
					return err
				}
			} else {
				configuredSols = append(configuredSols, sol.Solution)
			}
		}
		// Loop through user selected solutions, add any that aren't in current cluster
		for _, sol := range userSolutionList {
			if !stackpointio.StringInSlice(sol, configuredSols) {
				// Solution not in cluster, needs to be added
				newSolution := stackpointio.Solution{Solution: sol}
				solution, err := config.Client.AddSolution(config.OrgID, clusterID, newSolution)
				if err != nil {
					return err
				}
				// Wait until installed
				timeout := int(d.Timeout("Update").Seconds())
				if v, ok := d.GetOk("timeout"); ok {
					timeout = v.(int)
				}
				config.Client.WaitSolutionInstalled(config.OrgID, clusterID, solution.ID, timeout)
				log.Printf("[DEBUG] Added solution %s\n", sol)
			}
		}
	}
	return resourceStackPointClusterRead(d, meta)
}

func resourceStackPointClusterDelete(d *schema.ResourceData, meta interface{}) error {
	clusterID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	config := meta.(*Config)
	if err = config.Client.DeleteCluster(config.OrgID, clusterID); err != nil {
		return err
	}
	timeout := int(d.Timeout("Delete").Seconds())
	if v, ok := d.GetOk("timeout"); ok {
		timeout = v.(int)
	}
	if err = config.Client.WaitClusterDeleted(config.OrgID, clusterID, timeout); err != nil {
		return err
	}
	log.Println("[DEBUG] Cluster deletion complete")
	d.SetId("")
	return nil
}
