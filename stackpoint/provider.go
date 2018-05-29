package stackpoint

import (
	"fmt"
	"github.com/StackPointCloud/stackpoint-sdk-go/stackpointio"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"time"
)

// Provider returns a schema.Provider for StackPoint
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SPC_API_TOKEN", nil),
				Description: "The token key for API operations.",
			},
			"endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SPC_BASE_API_URL", nil),
				Description: "The endpoint URL for API operations.",
			},
			"org_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"ssh_keyset": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"stackpoint_cluster":     resourceStackPointCluster(),
			"stackpoint_master_node": resourceStackPointMasterNode(),
			"stackpoint_nodepool":    resourceStackPointNodePool(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"stackpoint_instance_specs": dataSourceStackPointInstanceSpecs(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	if _, ok := d.GetOk("token"); !ok {
		return nil, fmt.Errorf("StackPoint token has not been provided.")
	}
	if _, ok := d.GetOk("endpoint"); !ok {
		return nil, fmt.Errorf("StackPoint endpoint has not been provided.")
	}
	config := Config{
		Token:     d.Get("token").(string),
		EndPoint:  d.Get("endpoint").(string),
		Client:    stackpointio.NewClient(d.Get("token").(string), d.Get("endpoint").(string)),
		OrgID:     d.Get("org_id").(int),
		SSHKeyset: d.Get("ssh_keyset").(int),
	}
	return &config, nil
}

var resourceDefaultTimeouts = schema.ResourceTimeout{
	Create:  schema.DefaultTimeout(30 * time.Minute),
	Update:  schema.DefaultTimeout(30 * time.Minute),
	Delete:  schema.DefaultTimeout(30 * time.Minute),
	Default: schema.DefaultTimeout(30 * time.Minute),
}
