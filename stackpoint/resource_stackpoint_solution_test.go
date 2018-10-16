package stackpoint

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/StackPointCloud/stackpoint-sdk-go/stackpointio"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccStackPointSolution_basic(t *testing.T) {
	_, exists := os.LookupEnv("TF_ACC_SOLUTION_LOCK")
	if !exists {
		t.Skip("`TF_ACC_SOLUTION_LOCK` isn't specified - skipping since test will increase test time significantly")
	}

	var solution stackpointio.Solution
	nodeSize := "standard_f1"
	clusterName := "TerraForm AccTest"
	region := "eastus"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDStackPointSolutionDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccStackPointSolution_basic, nodeSize, clusterName, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("stackpoint_solution.efk", "solution", "efk"),
					testAccCheckStackPointSolutionExists("stackpoint_solution.efk", &solution),
				),
			},
		},
	})
}

func testAccCheckDStackPointSolutionDestroyCheck(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "stackpoint_solution" {
			continue
		}
		client := stackpointio.NewClient(os.Getenv("SPC_API_TOKEN"), os.Getenv("SPC_BASE_API_URL"))
		orgID, err := strconv.Atoi(rs.Primary.Attributes["org_id"])
		clID, err := strconv.Atoi(rs.Primary.Attributes["cluster_id"])
		if err != nil {
			return err
		}
		slID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}
		err = client.WaitSolutionDeleted(orgID, clID, slID, 3600)
		if err != nil {
			return fmt.Errorf("Error while waiting for cluster at ID %s to delete: %s", rs.Primary.ID, err)
		}
	}
	return nil
}

func testAccCheckStackPointSolutionExists(n string, sl *stackpointio.Solution) resource.TestCheckFunc {
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
		client := stackpointio.NewClient(os.Getenv("SPC_BASE_API_URL"), os.Getenv("SPC_BASE_API_URL"))
		solution, err := client.GetSolution(orgID, clID, slID)
		if err != nil {
			return fmt.Errorf("error occured while fetching cluster with ID %s: %s\ntoken: %s\nendpoint: %s\n",
				rs.Primary.ID, err, os.Getenv("token"), os.Getenv("endpoint"))
		}
		sl = solution

		return nil
	}
}

const testAccStackPointSolution_basic = `
data "stackpoint_keysets" "keyset_default" {

}
data "stackpoint_instance_specs" "master-specs" {
  provider_code = "azure"
  node_size     = "%s"
}
data "stackpoint_instance_specs" "worker-specs" {
  provider_code = "azure"
  node_size     = "${data.stackpoint_instance_specs.master-specs.node_size}"
}
resource "stackpoint_cluster" "terraform-cluster" {
  org_id                  = "${data.stackpoint_keysets.keyset_default.org_id}"
  cluster_name            = "%s"
  provider_code           = "azure"
  provider_keyset         = "${data.stackpoint_keysets.keyset_default.azure_keyset}"
  region                  = "%s"
  k8s_version             = "v1.9.6"
  startup_master_size     = "${data.stackpoint_instance_specs.master-specs.node_size}"
  startup_worker_count    = 2
  startup_worker_size     = "${data.stackpoint_instance_specs.worker-specs.node_size}"
  provider_network_cidr   = "10.0.0.0/16"
  provider_subnet_cidr    = "10.0.0.0/24"
  rbac_enabled            = true
  dashboard_enabled       = true
  etcd_type               = "classic"
  platform                = "coreos"
  channel                 = "stable"
  timeout                 = 1800
  ssh_keyset              = "${data.stackpoint_keysets.keyset_default.user_ssh_keyset}"
}

resource "stackpoint_solution" "efk"{
	org_id     = "${data.stackpoint_keysets.keyset_default.org_id}"
	cluster_id = "${stackpoint_cluster.terraform-cluster.id}"
	solution   = "efk"
}
`
