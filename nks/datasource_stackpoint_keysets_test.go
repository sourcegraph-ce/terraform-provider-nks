package stackpoint

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceKeysets_lookup(t *testing.T) {
	orgID := os.Getenv("SPC_ORG_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceStackPointKeysets_lookup, orgID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.nks_keysets.keyset-default", "org_id", orgID),
				),
			},
		},
	})
}

const testAccDataSourceStackPointKeysets_lookup = `
data "nks_keysets" "keyset-default" {
  org_id = "%s"
}
`
