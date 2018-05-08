package stackpoint

import (
	"encoding/json"
	"github.com/StackPointCloud/stackpoint-sdk-go/stackpointio"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strconv"
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
				SuppressDiff: false,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"provider_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"provider_key": {
				Type:     schema.TypeInt,
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
			"worker_count": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"worker_size": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
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
			"solutions": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},
	}
}

func resourceStackPointClusterCreate(d *schema.ResourceData, meta interface{}) error {
	newCluster := stackpointio.Cluster{
		Name:              d.Get("cluster_name").(string),
		Provider:          d.Get("provider_name").(string),
		ProviderKey:       d.Get("provider_key").(int),
		MasterCount:       d.Get("master_count").(int),
		MasterSize:        d.Get("master_size").(string),
		WorkerCount:       d.Get("worker_count").(int),
		WorkerSize:        d.Get("worker_size").(string),
		Region:            d.Get("region").(string),
		KubernetesVersion: d.Get("k8s_version").(string),
		RbacEnabled:       d.Get("rbac_enabled").(bool),
		DashboardEnabled:  d.Get("dashboard_enabled").(bool),
		EtcdType:          d.Get("etcd_type").(string),
		Platform:          d.Get("platform").(string),
		Channel:           d.Get("channel").(string),
		SSHKeySet:         d.Get("ssh_keyset").(int),
		Solutions:         []stackpointio.Solution{},
	}
	client := meta.(*stackpointio.APIClient)
	cluster, err := client.CreateCluster(d.Get("org_id").(int), newCluster)

	reqJSON, _ := json.Marshal(newCluster)
	resJSON, _ := json.Marshal(cluster)

	log.Println("[DEBUG] Cluster create request", string(reqJSON))
	log.Println("[DEBUG] Cluster create response", string(resJSON))

	// Don't bail until request and response are logged above
	if err != nil {
		return err
	}
	// Use following code for solutions list:
	//	if nRaw, ok := d.GetOk("nic"); ok {
	//		nicRaw := nRaw.(*schema.Set).List()

	// Wait until provisioned (until "state" is "running")
	for i := 1; ; i++ {
		state, err := client.GetClusterState(d.Get("org_id").(int), cluster.ID)
		if err != nil {
			return err
		}
		if state == "running" {
			d.SetId(strconv.Itoa(cluster.ID))
			break
		}
		time.Sleep(time.Second)
	}
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
		return err
	}
	d.Set("state", cluster.State)
	d.Set("instanceID", cluster.InstanceID)
	return nil
}

func resourceStackPointClusterUpdate(d *schema.ResourceData, meta interface{}) error {
        return nil                                             
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
