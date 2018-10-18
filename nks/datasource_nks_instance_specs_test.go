package nks

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceInstanceSpecs_lookup(t *testing.T) {
	nodeSize := "standard_f1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceInstanceSpecs_lookup, nodeSize),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.nks_instance_specs.master-specs", "node_size", nodeSize),
				),
			},
		},
	})

}

const testAccDataSourceInstanceSpecs_lookup = `
data "nks_instance_specs" "master-specs" {
  provider_code = "azure"
  node_size     = "%s"
}
`
