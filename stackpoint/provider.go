package stackpoint

import (
	"github.com/StackPointCloud/stackpoint-sdk-go/stackpointio"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
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
		},
		ResourcesMap: map[string]*schema.Resource{
			"stackpoint_cluster": resourceStackPointCluster(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	if token, ok := d.GetOk("token"); !ok {
		return nil, fmt.Errorf("StackPoint token has not been provided.")
	}
	if endpoint, ok := d.GetOk("endpoint"); !ok {
		return nil, fmt.Errorf("StackPoint endpoint has not been provided.")
	}
	return stackpointio.NewClient(token, endpoint), nil
}
