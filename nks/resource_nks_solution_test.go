package nks

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/NetApp/nks-sdk-go/nks"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNKSSolution_basic(t *testing.T) {
	_, exists := os.LookupEnv("TF_ACC_SOLUTION_LOCK")
	if !exists {
		t.Skip("`TF_ACC_SOLUTION_LOCK` isn't specified - skipping since test will increase test time significantly")
	}

	var solution nks.Solution
	nodeSize := "standard_f1"
	clusterName := "TerraForm AccTest Solution"
	region := "eastus"
	solutionName := "haproxy"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccNKSSolution_basic, nodeSize, clusterName, region, solutionName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nks_solution.efk", "solution", solutionName),
					testAccCheckNKSSolutionExists("nks_solution.efk", &solution),
				),
			},
		},
	})
}

func testAccCheckNKSSolutionExists(n string, sl *nks.Solution) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		if rs.Primary.Attributes["org_id"] == "" {
			return fmt.Errorf("No Org ID is set")
		}
		orgID, err := strconv.Atoi(rs.Primary.Attributes["org_id"])
		if err != nil {
			return err
		}
		clID, err := strconv.Atoi(rs.Primary.Attributes["cluster_id"])
		if err != nil {
			return err
		}
		slID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}
		client := nks.NewClient(os.Getenv("NKS_API_URL"), os.Getenv("NKS_API_URL"))
		solution, err := client.GetSolution(orgID, clID, slID)
		if err != nil {
			return fmt.Errorf("error occured while fetching solution: %s", err)
		}
		sl = solution

		return nil
	}
}

const testAccNKSSolution_basic = `
data "nks_organization" "org"{

}
data "nks_keyset" "keyset_default" {
	category = "provider"
	entity = "azure"
}

data "nks_keyset" "ssh" {
	category = "user_ssh"
	name = "default"
}

data "nks_instance_specs" "master-specs" {
  provider_code = "azure"
  node_size     = "%s"
}
data "nks_instance_specs" "worker-specs" {
  provider_code = "azure"
  node_size     = "${data.nks_instance_specs.master-specs.node_size}"
}
resource "nks_cluster" "terraform-cluster" {
  org_id                  = "${data.nks_organization.org.id}"
  cluster_name            = "%s"
  provider_code           = "azure"
  provider_keyset         = "${data.nks_keyset.keyset_default.id}"
  region                  = "%s"
  k8s_version             = "v1.9.6"
  startup_master_size     = "${data.nks_instance_specs.master-specs.node_size}"
  startup_worker_count    = 2
  startup_worker_size     = "${data.nks_instance_specs.worker-specs.node_size}"
  provider_network_cidr   = "10.0.0.0/16"
  provider_subnet_cidr    = "10.0.0.0/24"
  rbac_enabled            = true
  dashboard_enabled       = true
  etcd_type               = "classic"
  platform                = "coreos"
  channel                 = "stable"
  timeout                 = 1800
  ssh_keyset              = "${data.nks_keyset.ssh.id}"
}

resource "nks_solution" "efk"{
	org_id     = "${data.nks_organization.org.id}"
	cluster_id = "${nks_cluster.terraform-cluster.id}"
	solution   = "%s"

}
`
