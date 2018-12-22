package nks

import (
	"fmt"
	"log"

	"github.com/StackPointCloud/nks-sdk-go/nks"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNKSInstanceSpecs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNKSInstanceSpecsRead,
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

func dataSourceNKSInstanceSpecsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// Grab machine size values for provider, using optional endpoint if needed (sometimes different machines loaded on staging vs prod)
	var endpoint string
	if ep, ok := d.GetOk("endpoint"); ok {
		endpoint = ep.(string)
	}
	mOptions, err := config.Client.GetInstanceSpecs(d.Get("provider_code").(string), endpoint)
	if err != nil {
		log.Printf("[DEBUG] InstanceSpecs GetInstanceSpecs failed: %s\n", err)
		return err
	}
	if !nks.InstanceInList(mOptions, d.Get("node_size").(string)) {
		return fmt.Errorf("invalid machine size for node: %s", d.Get("node_size").(string))
	}
	d.Set("provider_code", d.Get("provider_code").(string))
	d.Set("size", d.Get("node_size").(string))
	d.SetId("1") // This is just a holder for now, there are no numerical values for instances in our system

	return nil
}
