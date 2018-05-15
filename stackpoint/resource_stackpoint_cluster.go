package stackpoint

import (
	"encoding/json"
	"github.com/StackPointCloud/stackpoint-sdk-go/stackpointio"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"fmt"
	"strconv"
	"strings"
	"time"
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
			"provider_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provider_keyset": {
				Type:     schema.TypeInt,
				Required: true,
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
			"master_count": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"master_size": {
				Type:     schema.TypeString,
				Required: true,
			},
                        "worker_size": {
                                Type:     schema.TypeString,
                                Required: true,
                        },
			"worker_count": {
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
			"provider_resource_group": {
				Type:	 schema.TypeString,
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
			"k8s_version": {
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
                        "dashboard_installed": {
                                Type:     schema.TypeBool,
                                Computed: true,
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
			"solutions": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
}

func resourceStackPointClusterCreate(d *schema.ResourceData, meta interface{}) error {
	// Get client for API
        client := meta.(*stackpointio.APIClient)

        // Grab machine size values for provider
        mOptions, err := client.GetInstanceSpecs(d.Get("provider_name").(string))
        if err != nil {
                return err
        }
        // Validate worker node size
        if !stackpointio.InstanceInList(mOptions, d.Get("worker_size").(string)) {
                return fmt.Errorf("Invalid machine size for worker node: %s\n", d.Get("worker_size").(string))
        }
        // Validate master node size
        if !stackpointio.InstanceInList(mOptions, d.Get("master_size").(string)) {
                return fmt.Errorf("Invalid machine size for master node: %s\n", d.Get("master_size").(string))
        }
	// Set up cluster structure based on input from user
	newCluster := stackpointio.Cluster{
		Name:              d.Get("cluster_name").(string),
		Provider:          d.Get("provider_name").(string),
		ProviderKey:       d.Get("provider_keyset").(int),
		MasterCount:       d.Get("master_count").(int),
		MasterSize:        d.Get("master_size").(string),
		WorkerCount:	   d.Get("worker_count").(int),
		WorkerSize:	   d.Get("worker_size").(string),
		KubernetesVersion: d.Get("k8s_version").(string),
		RbacEnabled:       d.Get("rbac_enabled").(bool),
		DashboardEnabled:  d.Get("dashboard_enabled").(bool),
		EtcdType:          d.Get("etcd_type").(string),
		Platform:          d.Get("platform").(string),
		Channel:           d.Get("channel").(string),
		SSHKeySet:         d.Get("ssh_keyset").(int),
		Solutions: 	   []stackpointio.Solution{},  // helm_tiller will get automatically installed
	}
	// Grab provider-specific fields
	if d.Get("provider_name").(string) == "aws" {
                if _, ok := d.GetOk("region"); !ok {
                        return fmt.Errorf("StackPoint needs region for AWS clusters.")
                }
		if _, ok := d.GetOk("zone"); !ok {
			return fmt.Errorf("StackPoint needs zone for AWS clusters.")
		}
		if _, ok := d.GetOk("provider_network_id"); !ok {
			return fmt.Errorf("StackPoint needs provider_network_id for AWS clusters.")
		}
		if _, ok := d.GetOk("provider_network_cidr"); !ok {
			return fmt.Errorf("StackPoint needs provider_network_cidr for AWS clusters.")
		}
		if _, ok := d.GetOk("provider_subnet_id"); !ok {
			return fmt.Errorf("StackPoint needs provider_subnet_id for AWS clusters.")
		}
		if _, ok := d.GetOk("provider_subnet_cidr"); !ok {
			return fmt.Errorf("StackPoint needs provider_subnet_cidr for AWS clusters.")
		}
		newCluster.Region = d.Get("region").(string)
		newCluster.Zone = d.Get("zone").(string)
		newCluster.ProviderNetworkID = d.Get("provider_network_id").(string)
		newCluster.ProviderNetworkCdr = d.Get("provider_network_cidr").(string)
		newCluster.ProviderSubnetID = d.Get("provider_subnet_id").(string)
		newCluster.ProviderSubnetCidr = d.Get("provider_subnet_cidr").(string)
	} else if d.Get("provider_name").(string) == "do" || 
		d.Get("provider_name").(string) == "gce" || d.Get("provider_name").(string) == "gke" {
                if _, ok := d.GetOk("region"); !ok {
                        return fmt.Errorf("StackPoint needs region for DigitalOcean/GCE/GKE clusters.")
                }
		newCluster.Region = d.Get("region").(string)
	} else if d.Get("provider_name").(string) == "azure" {
                if _, ok := d.GetOk("provider_resource_group"); !ok {
                        return fmt.Errorf("StackPoint needs provider_resource_group for Azure clusters.")
                }
                if _, ok := d.GetOk("region"); !ok {
                        return fmt.Errorf("StackPoint needs region for Azure clusters.")
                }
                if _, ok := d.GetOk("provider_network_id"); !ok {
                        return fmt.Errorf("StackPoint needs provider_network_id for Azure clusters.")
                }
                if _, ok := d.GetOk("provider_network_cidr"); !ok {
                        return fmt.Errorf("StackPoint needs provider_network_cidr for Azure clusters.")
                }
                if _, ok := d.GetOk("provider_subnet_id"); !ok {
                        return fmt.Errorf("StackPoint needs provider_subnet_id for Azure clusters.")
                }
                if _, ok := d.GetOk("provider_subnet_cidr"); !ok {
                        return fmt.Errorf("StackPoint needs provider_subnet_cidr for Azure clusters.")
                }
		newCluster.ProviderResourceGp = d.Get("provider_resource_group").(string)
		newCluster.Region = d.Get("region").(string)
                newCluster.ProviderNetworkID = d.Get("provider_network_id").(string)
                newCluster.ProviderNetworkCdr = d.Get("provider_network_cidr").(string)
                newCluster.ProviderSubnetID = d.Get("provider_subnet_id").(string)
                newCluster.ProviderSubnetCidr = d.Get("provider_subnet_cidr").(string)
        }
	// Do cluster creation call
	cluster, err := client.CreateCluster(d.Get("org_id").(int), newCluster)

	reqJSON, _ := json.Marshal(newCluster)
	resJSON, _ := json.Marshal(cluster)

	log.Println("[DEBUG] Cluster create request", string(reqJSON))
	log.Println("[DEBUG] Cluster create response", string(resJSON))

	// Don't bail until request and response are logged above
	if err != nil {
		return err
	}

	// Wait until provisioned
	err = client.WaitClusterProvisioned(d.Get("org_id").(int), cluster.ID)
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
	client := meta.(*stackpointio.APIClient)
	cluster, err := client.GetCluster(d.Get("org_id").(int), clusterID)
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
  	d.Set("org_id", cluster.OrganizationKey)
  	d.Set("provider_name", cluster.Provider)
  	d.Set("provider_keyset", cluster.ProviderKey)
  	d.Set("provider_keyset_name", cluster.ProviderKeyName)
  	d.Set("region", cluster.Region)
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
  	d.Set("user_ssh_keyset", cluster.SSHKeySet)

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
        client := meta.(*stackpointio.APIClient)

	if d.HasChange("master_count") {
		oldV, newV := d.GetChange("worker_count")

		// Only allow 1 master to be deleted or added at a time (could change this if API will handle it)
		if oldV.(int) > newV.(int) {
			// User asks to reduce master count
			// Only allow 1 master to be deleted at a time (could change this if API will handle it)
			if (oldV.(int) - newV.(int)) > 1 {
				return fmt.Errorf("Only a single reduction to master count is allowed at this time.")
			}
			nodes, err := client.GetNodes(d.Get("org_id").(int), clusterID)
                        if err != nil {
                                log.Println("[DEBUG] Cluster update got error when getting node list")
                                return err
                        }
                        if len(nodes) == 0 {
                                return fmt.Errorf("No nodes found to reduce master count.")
                        }
			if len(nodes) < 2 {
                                return fmt.Errorf("Cannot reduce master count below 1.")
                        }
                        for i := 0; i < len(nodes); i++ {
                                // Sort nodes by master class and running state
                                if nodes[i].Role == "master" && nodes[i].State == "running" {
                                        if err := client.DeleteNode(d.Get("org_id").(int), clusterID, nodes[i].ID); err != nil {
                                                log.Printf("[DEBUG] Cluster update got error when deleting node at ID: %d\n",
                                                        nodes[i].ID)
                                                return err
                                        }
					// Allow some time before eventual state read call
					time.Sleep(10)
                                        break
                                }
			}
		} else {
                        // User asks to increase master count
			// Only allow 1 master to be added at a time (could change this if API will handle it)
                        if (newV.(int) - oldV.(int)) > 1 {
                                return fmt.Errorf("Only a single addition to master count is allowed at this time.")
                        }
        		// Set up new master node
        		newNode := stackpointio.NodeAdd {
                		Count: 1,
                		Role:  "master",
                		Size:  d.Get("master_size").(string),
			}
			if d.Get("provider_name").(string) == "aws" {
                		if _, ok := d.GetOk("zone"); !ok {
                        		return fmt.Errorf("StackPoint needs zone for AWS clusters.")
                		}
                		if _, ok := d.GetOk("provider_subnet_id"); !ok {
                        		return fmt.Errorf("StackPoint needs provider_subnet_id for AWS clusters.")
                		}
                		if _, ok := d.GetOk("provider_subnet_cidr"); !ok {
                        		return fmt.Errorf("StackPoint needs provider_subnet_cidr for AWS clusters.")
                		}
                		newNode.Zone = d.Get("zone").(string)
                		newNode.ProviderSubnetID = d.Get("provider_subnet_id").(string)
                		newNode.ProviderSubnetCidr = d.Get("provider_subnet_cidr").(string)
			}
                        log.Println("[DEBUG] Cluster update attempting to add master node\n")
                        nodes, err := client.AddNode(d.Get("org_id").(int), clusterID, newNode)
                        if err != nil {
                                return err
                        }
                        for _, node := range nodes {
                                if err := client.WaitNodeProvisioned(d.Get("org_id").(int), clusterID, node.ID); err != nil {
                                        return err
                                }
                        }
			// Allow some time before eventual state read call
			time.Sleep(10)
                }
	}
	if d.HasChange("worker_count") {
		oldV, newV := d.GetChange("worker_count")

		// Don't allow worker count to reduce to less than 1
		if newV.(int) < 1 {
                        return fmt.Errorf("Cannot reduce to less than 1 worker node.")
                }
		if oldV.(int) > newV.(int) {
			// User asks to reduce worker count, get list of nodes
			nodes, err := client.GetNodes(d.Get("org_id").(int), clusterID)
        		if err != nil {
				log.Println("[DEBUG] Cluster update got error when getting node list")
                		return err
        		}
		        if len(nodes) == 0 {
                		return fmt.Errorf("No nodes found to reduce worker count.")
        		}
			num_nodes_deleted := 0
		        for i := 0; i < len(nodes); i++ {
				// Sort nodes by worker class and running state
				if nodes[i].Role == "worker" && nodes[i].State == "running" {
					if err := client.DeleteNode(d.Get("org_id").(int), clusterID, nodes[i].ID); err != nil {
						log.Printf("[DEBUG] Cluster update got error when deleting node at ID: %d\n", 
							nodes[i].ID)
						return err
					}
					// Pause for a couple seconds for node state to reflect that it's deleting
					time.Sleep(10)
					num_nodes_deleted = num_nodes_deleted + 1
					if oldV.(int) - newV.(int) - num_nodes_deleted == 0 {
						// Number of nodes deleted should be reached now
						break;
					} 
				}
			}
			if oldV.(int) - newV.(int) - num_nodes_deleted != 0 {
				return fmt.Errorf("Error deleting nodes, node count is off after deletion")
			}
		} else {
			// User asks to increase worker count, get nodepool list, add to first available one
		        nps, err := client.GetNodePools(d.Get("org_id").(int), clusterID)
        		if err != nil {
				log.Println("[DEBUG] Cluster update got error when getting nodepool list")
                		return err
        		}
        		if len(nps) == 0 {
				// No nodepools found to add to, throw error
                		return fmt.Errorf("No nodepools found to add worker node to")
        		}
			// Make sure nodepool exists and is active
                	//if nps[0].State != "active" {
				//return fmt.Errorf("No active nodepools found")
			//}
        		newNode := stackpointio.NodeAddToPool {
				Count:      (newV.(int) - oldV.(int)),
                		Role:       "worker",
                		NodePoolID:  nps[0].ID,
			}
			log.Printf("[DEBUG] Cluster update attempting to add %d worker node(s)\n", (newV.(int) - oldV.(int)))
        		nodes, err := client.AddNodesToNodePool(d.Get("org_id").(int), clusterID, nps[0].ID, newNode)
        		if err != nil {
                		return err
        		}
			for _, node := range nodes {
				if err := client.WaitNodeProvisioned(d.Get("org_id").(int), clusterID, node.ID); err != nil {
					return err
				}
			}
			// Pause for a couple seconds for new node to appear
			time.Sleep(10)
		}
        }
	if d.HasChange("solutions") {
		_, newV := d.GetChange("solutions")
		userIntList := newV.([]interface{})
		var userSolutionList []string
		for _, item := range userIntList {
			userSolutionList = append(userSolutionList, item.(string))
		}
		solutions, err := client.GetSolutions(d.Get("org_id").(int), clusterID)
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
				if err := client.DeleteSolution(d.Get("org_id").(int), clusterID, sol.ID); err != nil {
					return err
				}
				// Pause for a few seconds for solution deletion to report
				time.Sleep(10)
			} else {
				configuredSols = append(configuredSols, sol.Solution)
			}
        	}
		// Loop through user selected solutions, add any that aren't in current cluster
		for _, sol := range userSolutionList {
			if !stackpointio.StringInSlice(sol, configuredSols) {
                                // Solution not in cluster, needs to be added
        			newSolution := stackpointio.Solution { Solution: sol }
				_, err := client.AddSolution(d.Get("org_id").(int), clusterID, newSolution);
				if err != nil {
                                        return err
                                }
				// Pause for a few seconds to let new solution show up before state read call
				time.Sleep(10)
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
	client := meta.(*stackpointio.APIClient)
	err = client.DeleteCluster(d.Get("org_id").(int), clusterID)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
