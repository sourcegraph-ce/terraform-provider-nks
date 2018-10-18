package nks

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceKeysets_lookup(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNKSKeysets_lookup,
			},
		},
	})
}

const testAccDataSourceNKSKeysets_lookup = `
data "nks_keyset" "keyset-default" {
	category = "provider"
	entity = "azure"
}
`
