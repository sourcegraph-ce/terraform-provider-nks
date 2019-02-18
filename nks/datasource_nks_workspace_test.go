package nks

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceWorkspace_lookup(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNKSWorkspace_lookup,
			},
		},
	})
}

const testAccDataSourceNKSWorkspace_lookup = `
data "nks_organization" "org" {
}

data "nks_workspace" "my-workspace" {
	org_id = "${data.nks_organization.org.id}"
}
`
