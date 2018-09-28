package stackpoint

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccDataSourceInstanceSpecs_lookup(t *testing.T) {
	nodeSize := "n1-standard-1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceInstanceSpecs_lookup, nodeSize),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.stackpoint_instance_specs.master-specs", "node_size", nodeSize),
				),
			},
		},
	})

}

const testAccDataSourceInstanceSpecs_lookup = `
data "stackpoint_instance_specs" "master-specs" {
  provider_code = "gce"
  node_size     = "%s"
}
`
