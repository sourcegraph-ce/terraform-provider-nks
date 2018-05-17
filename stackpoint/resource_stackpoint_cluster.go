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
	"os"
)

const operationTimeout = 1200  // 20 minutes

func outDebug(m string) {
	f, err := os.OpenFile("/tmp/tf_debug.txt", os.O_APPEND|os.O_WRONLY, 0644)
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
                        "master_count": {
                                Type:     schema.TypeInt,
                                Required: true,
                        },
                        "master_size": {
                                Type:     schema.TypeString,
                                Required: true,
                        },
                        "nodepool": {
                                Type:     schema.TypeSet,
                                Required: true,
                                Elem: &schema.Resource{
                                        Schema: map[string]*schema.Schema{
						"local_id": {
                                                        Type:     schema.TypeInt,
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
                                        },
                                },
				Set: nodepoolHash,
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
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project_id": {
                                Type:    schema.TypeString,
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
                        "dashboard_installed": {
                                Type:     schema.TypeBool,
                                Computed: true,
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
        // Grab worker node info from nodepool config
        nodepoolList := d.Get("nodepool").(*schema.Set).List()
	var workerSize string
	var workerCount int
outDebug(fmt.Sprintf("len(nodepoolList): %d\n", len(nodepoolList)))
        for i, element := range nodepoolList {
outDebug(fmt.Sprintf("In nodepool loop for creation, i=%d\n", i))
                if i > 0 {
                        return fmt.Errorf("Sorry, the StackPoint plugin only supports a single nodepool at creation time.")
                }
                elementMap := element.(map[string]interface{})
                if elementMap["worker_size"] != nil {
outDebug(fmt.Sprintf("In nodepool loop for creation, i=%d, worker_size=%s\n", i, elementMap["worker_size"].(string)))
                        // Validate worker node size
                        if !stackpointio.InstanceInList(mOptions, elementMap["worker_size"].(string)) {
                                return fmt.Errorf("Invalid machine size for worker node: %s\n", elementMap["worker_size"].(string))
                        }
			workerSize = elementMap["worker_size"].(string)
                }
                if elementMap["worker_count"] != nil {
outDebug(fmt.Sprintf("In nodepool loop for creation, i=%d, worker_count=%d\n", i, elementMap["worker_count"].(int)))
			workerCount = elementMap["worker_count"].(int)
                }
        }
        // Make sure we have workerSize and workerCount for cluster build
	if workerSize == "" || workerCount == 0 {
                return fmt.Errorf("Missing worker_size or worker_count.")
        }
        // Make sure at least 2 worker nodes (currently at least 2 worker nodes are required at creation)
        if workerCount < 2 {
                return fmt.Errorf("Need at least 2 worker nodes to create a cluster.")
        }
        // Validate master node size
        if !stackpointio.InstanceInList(mOptions, d.Get("master_size").(string)) {
                return fmt.Errorf("Invalid machine size for master node: %s\n", d.Get("master_size").(string))
        }
	// Make sure only single master (only single master allowed at creation time currently)
	if d.Get("master_count").(int) > 1 {
		return fmt.Errorf("Only a single master node is allowed at creation time currently.")
	}
	// Set up cluster structure based on input from user
	newCluster := stackpointio.Cluster{
		Name:              d.Get("cluster_name").(string),
		Provider:          d.Get("provider_name").(string),
		ProviderKey:       d.Get("provider_keyset").(int),
		MasterCount:       d.Get("master_count").(int),
		MasterSize:        d.Get("master_size").(string),
		WorkerCount:	   workerCount,
		WorkerSize:	   workerSize,
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
	} else if d.Get("provider_name").(string) == "do" || d.Get("provider_name").(string) == "gce" || 
		d.Get("provider_name").(string) == "gke" {
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
        } else if d.Get("provider_name").(string) == "packet" {
                if _, ok := d.GetOk("region"); !ok {
                        return fmt.Errorf("StackPoint needs region for Packet clusters.")
                }
                if _, ok := d.GetOk("project_id"); !ok {
                        return fmt.Errorf("StackPoint needs project_id for Packet clusters.")
                }
                newCluster.Region = d.Get("region").(string)
		newCluster.ProjectID = d.Get("provider_id").(string)
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
	if v, ok := d.GetOk("zone"); ok {
		d.Set("zone", v.(string))
	}
	if v, ok := d.GetOk("project_id"); ok {
		d.Set("project_id", v.(string))
	}
	if v, ok := d.GetOk("provider_resource_group"); ok {
        	d.Set("provider_resource_group", v.(string))
	}
	if v, ok := d.GetOk("provider_network_id"); ok {
        	d.Set("provider_network_id", v.(string))
	}
	if v, ok := d.GetOk("provider_network_cidr"); ok {
        	d.Set("provider_network_cidr", v.(string))
	}
	if v, ok := d.GetOk("provider_subnet_id"); ok {
        	d.Set("provider_subnet_id", v.(string))
	}
	if v, ok := d.GetOk("provider_subnet_cidr"); ok {
        	d.Set("provider_subnet_cidr", v.(string))
	}
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

	// Collect solutions for TF state
	var solArray []string
	for _, sol := range cluster.Solutions {
		solArray = append(solArray, sol.Solution)
	}
	d.Set("solutions", solArray)

	// Collect nodepool info from cluster
        nps, err := client.GetNodePools(d.Get("org_id").(int), clusterID)
        if err != nil {
                return err
        }
        nodepoolMap := nodepoolsToMap(nps)
	d.Set("nodepool", nodepoolMap)

	return nil
}

func resourceStackPointClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	clusterID, err := strconv.Atoi(d.Id())
        if err != nil {
                return err
        }
        client := meta.(*stackpointio.APIClient)

	if d.HasChange("master_count") {
		oldV, newV := d.GetChange("master_count")
		oldVi, newVi := oldV.(int), newV.(int)

		// Only allow 1 master to be deleted or added at a time (could change this if API will handle it)
		if oldVi > newVi {
			// User asks to reduce master count
			// Only allow 1 master to be deleted at a time (could change this if API will handle it)
			if (oldVi - newVi) > 1 {
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
                                        // Wait for node to delete
                                        if err = client.WaitNodeDeleted(d.Get("org_id").(int), clusterID, nodes[i].ID); err != nil {
                                                return err
                                        }
                                        break
                                }
			}
		} else {
                        // User asks to increase master count
			// Only allow 1 master to be added at a time (could change this if API will handle it)
                        if (newVi - oldVi) > 1 {
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
                }
	}
	if d.HasChange("nodepool") {
outDebug(fmt.Sprintf("In nodepool has change\n"))
                oldV, newV := d.GetChange("nodepool")
                oldVSet, newVSet := oldV.(*schema.Set), newV.(*schema.Set)
                nodepoolsToAdd := newVSet.Difference(oldVSet)
                nodepoolsToRemove := oldVSet.Difference(newVSet)
                for i, rawNP := range nodepoolsToAdd.List() {
outDebug(fmt.Sprintf("In nodepoolToAdd, i=%d\n", i))
                        rawNPMap := rawNP.(map[string]interface{})
outDebug(fmt.Sprintf("rawNPMap[worker_count]: %d\n", rawNPMap["worker_count"].(int)))
			newNodepool := stackpointio.NodePool {
				Name: fmt.Sprintf("TerraForm NodePool %d", i + 1),
				NodeCount: rawNPMap["worker_count"].(int),
				Size:      rawNPMap["worker_size"].(string),
				Platform:  d.Get("platform").(string),
			}
			// Create new nodepool
			pool, err := client.CreateNodePool(d.Get("org_id").(int), clusterID, newNodepool)
			if err != nil {
				log.Fatal(err)
			}
			client.WaitNodePoolProvisioned(d.Get("org_id").(int), clusterID, pool.ID)
                }
                for i, rawNP := range nodepoolsToRemove.List() {
			// Currently impossible to delete nodepools, so just log this for now
			log.Println("[DEBUG] Cluster update attempting to add master node\n")
outDebug(fmt.Sprintf("In nodepoolToRemove, i=%d\n", i))
                        rawNPMap := rawNP.(map[string]interface{})
outDebug(fmt.Sprintf("rawNPMap[worker_count]: %d\n", rawNPMap["worker_count"].(int)))
outDebug(fmt.Sprintf("rawNPMap[instance_id]: %d\n", rawNPMap["instance_id"].(string)))
                }
        }
	if d.HasChange("worker_count") {
outDebug(fmt.Sprintf("In worker_count has change???????\n"))
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
				// Pause for a few seconds for solution deletion to report (no way to wait for state on this)
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
				solution, err := client.AddSolution(d.Get("org_id").(int), clusterID, newSolution);
				if err != nil {
                                        return err
                                }
				// Wait until installed
				client.WaitSolutionInstalled(d.Get("org_id").(int), clusterID, solution.ID)
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
	if err = client.DeleteCluster(d.Get("org_id").(int), clusterID); err != nil {
		return err
	}
	if err = client.WaitClusterDeleted(d.Get("org_id").(int), clusterID); err != nil {
		return err
	}
        log.Println("[DEBUG] Cluster deletion complete")
	d.SetId("")
	return nil
}

func nodepoolsToMap(nps []stackpointio.NodePool) []map[string]interface{} {
        nodepoolMap := make([]map[string]interface{}, len(nps))
        for i, np := range nps {
outDebug(fmt.Sprintf("In nodepoolsToMap loop, i=%d, np.Name=%s, np.InstanceID=%s, np.NodeCount=%d\n", i, np.Name, np.InstanceID, np.NodeCount))
                nodepoolMap[i] = map[string]interface{} {
                        // local_id will start on 1, so increase i by 1
                        "local_id": i+1,
                        "name": np.Name,
                        "instance_id": np.InstanceID,
                        "autoscaled": np.Autoscaled,
                        "autoscale_min_count": np.MinCount,
                        "autoscale_max_count": np.MaxCount,
                        "worker_count": np.NodeCount,
                        "worker_size": np.Size,
                        "state": np.State,
                }
        }
	return nodepoolMap
}

func nodepoolHash(v interface{}) int {
	m := v.(map[string]interface{})
	return m["local_id"].(int)
}
