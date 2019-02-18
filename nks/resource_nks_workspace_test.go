package nks

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccResourceWorkspace_Basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccResourceNKSWorkspace_Basic,
			},
		},
	})
}

const testAccResourceNKSWorkspace_Basic = `
data "nks_organization" "org"{
	
}

resource "nks_workspace" "my-workspace" {
    org_id          = "${data.nks_organization.org.id}"
	default         = false
	name            = "My Test Workspace"
}
`
