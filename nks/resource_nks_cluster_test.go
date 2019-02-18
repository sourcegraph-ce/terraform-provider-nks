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

func TestAccNKS_basic(t *testing.T) {
	var cluster nks.Cluster
	nodeSize := "standard_f1"
	clusterName := "TerraForm Acceptance Test"
	region := "eastus"
	vpcCIDR := "10.0.0.0/16"
	subnetCIDR := "10.0.0.0/24"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNKSClusterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccNKSCluster_basic, nodeSize, clusterName, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.nks_instance_specs.master-specs", "node_size", nodeSize),
					resource.TestCheckResourceAttr("data.nks_instance_specs.worker-specs", "node_size", nodeSize),
					resource.TestCheckResourceAttr("nks_cluster.terraform-cluster", "cluster_name", clusterName),
					resource.TestCheckResourceAttr("nks_cluster.terraform-cluster", "region", region),
					resource.TestCheckResourceAttr("nks_cluster.terraform-cluster", "provider_network_cidr", vpcCIDR),
					resource.TestCheckResourceAttr("nks_cluster.terraform-cluster", "provider_subnet_cidr", subnetCIDR),
					testAccCheckNKSClusterExists("nks_cluster.terraform-cluster", &cluster),
				),
			},
		},
	})
}

func testAccCheckDNKSClusterDestroyCheck(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nks_cluster" {
			continue
		}
		client := nks.NewClient(os.Getenv("NKS_API_TOKEN"), os.Getenv("NKS_API_URL"))
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

func testAccCheckNKSClusterExists(n string, cl *nks.Cluster) resource.TestCheckFunc {
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
		client := nks.NewClient(os.Getenv("NKS_API_URL"), os.Getenv("NKS_API_URL"))
		cluster, err := client.GetCluster(orgID, clID)
		if err != nil {
			return fmt.Errorf("Error occured while fetching cluster with ID %s: %s\ntoken: %s\nendpoint: %s\n",
				rs.Primary.ID, err, os.Getenv("token"), os.Getenv("endpoint"))
		}
		cl = cluster

		return nil
	}
}

const testAccNKSCluster_basic = `
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
`
