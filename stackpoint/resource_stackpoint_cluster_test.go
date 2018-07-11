package stackpoint

import (
	"fmt"
	"github.com/StackPointCloud/stackpoint-sdk-go/stackpointio"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"strconv"
	"testing"
)

func TestAccStackPointCluster_basic(t *testing.T) {
	var cluster stackpointio.Cluster
	nodeSize := "n1-standard-1"
	clusterName := "TerraForm AccTest"
	region := "us-west1-a"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDStackPointClusterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccStackPointCluster_basic, nodeSize, clusterName, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.stackpoint_instance_specs.master-specs", "node_size", nodeSize),
					resource.TestCheckResourceAttr("data.stackpoint_instance_specs.worker-specs", "node_size", nodeSize),
					resource.TestCheckResourceAttr("stackpoint_cluster.terraform-cluster", "cluster_name", clusterName),
					resource.TestCheckResourceAttr("stackpoint_cluster.terraform-cluster", "region", region),
					testAccCheckStackPointClusterExists("stackpoint_cluster.terraform-cluster", &cluster),
				),
			},
		},
	})
}

func testAccCheckDStackPointClusterDestroyCheck(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "stackpoint_cluster" {
			continue
		}
		client := stackpointio.NewClient(os.Getenv("SPC_API_TOKEN"), os.Getenv("SPC_BASE_API_URL"))
		orgID, err := strconv.Atoi(rs.Primary.Attributes["org_id"])
		if err != nil {
			return err
		}
		clID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}
		err = client.WaitClusterDeleted(orgID, clID, 3600)
		if err != nil {
			return fmt.Errorf("Error while waiting for cluster at ID %s to delete: %s", rs.Primary.ID, err)
		}
	}
	return nil
}

func testAccCheckStackPointClusterExists(n string, cl *stackpointio.Cluster) resource.TestCheckFunc {
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
		clID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}
		client := stackpointio.NewClient(os.Getenv("SPC_BASE_API_URL"), os.Getenv("SPC_BASE_API_URL"))
		cluster, err := client.GetCluster(orgID, clID)
		if err != nil {
			return fmt.Errorf("Error occured while fetching cluster with ID %s: %s\ntoken: %s\nendpoint: %s\n",
				rs.Primary.ID, err, os.Getenv("token"), os.Getenv("endpoint"))
		}
		cl = cluster

		return nil
	}
}

const testAccStackPointCluster_basic = `
data "stackpoint_keysets" "keyset_default" {

}
data "stackpoint_instance_specs" "master-specs" {
  provider_code = "gce"
  node_size     = "%s"
}
data "stackpoint_instance_specs" "worker-specs" {
  provider_code = "gce"
  node_size     = "${data.stackpoint_instance_specs.master-specs.node_size}"
}
resource "stackpoint_cluster" "terraform-cluster" {
  org_id                = "${data.stackpoint_keysets.keyset_default.org_id}"
  cluster_name          = "%s"
  provider_code         = "gce"
  provider_keyset       = "${data.stackpoint_keysets.keyset_default.gce_keyset}"
  region                = "%s"
  k8s_version           = "v1.9.6"
  startup_master_size   = "${data.stackpoint_instance_specs.master-specs.node_size}"
  startup_worker_count  = 2
  startup_worker_size   = "${data.stackpoint_instance_specs.worker-specs.node_size}"
  rbac_enabled          = true
  dashboard_enabled     = true
  etcd_type             = "classic"
  platform              = "coreos"
  channel               = "stable"
  ssh_keyset            = "${data.stackpoint_keysets.keyset_default.user_ssh_keyset}"
}
`
