package stackpoint

import (
	"fmt"
	"github.com/StackPointCloud/stackpoint-sdk-go/stackpointio"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceStackPointInstanceSpecs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceStackPointInstanceSpecsRead,
		Schema: map[string]*schema.Schema{
			"provider_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_size": {
				Type:     schema.TypeString,
				Required: true,
			},
                        "endpoint": {
                                Type:     schema.TypeString,
                                Optional: true,
                        },
		},
	}
}

func dataSourceStackPointInstanceSpecsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// Grab machine size values for provider
	mOptions, err := config.Client.GetInstanceSpecs(d.Get("provider_code").(string), d.Get("endpoint").(string))
	if err != nil {
		return err
	}
	if !stackpointio.InstanceInList(mOptions, d.Get("node_size").(string)) {
		return fmt.Errorf("Invalid machine size for node: %s\n", d.Get("node_size").(string))
	}
	d.Set("provider_code", d.Get("provider_code").(string))
	d.Set("size", d.Get("node_size").(string))
	d.SetId("1") // This is just a holder for now, there are no numerical values for instances in our system
	return nil
}
