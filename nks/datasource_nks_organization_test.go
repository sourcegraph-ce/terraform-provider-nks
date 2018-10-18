package nks

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceOrganization_lookup(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNKSOrganization_lookup,
			},
		},
	})
}

const testAccDataSourceNKSOrganization_lookup = `
data "nks_organization" "default" {
}
`
