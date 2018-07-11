package stackpoint

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccDataSourceKeysets_lookup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceStackPointKeysets_lookup,
			},
		},
	})

}

const testAccDataSourceStackPointKeysets_lookup = `
data "stackpoint_keysets" "keyset_default" {

}
`
