package nks

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/StackPointCloud/nks-sdk-go/nks"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNKSMasterNode_basic(t *testing.T) {
	_, exists := os.LookupEnv("TF_ACC_MASTER_NODE_LOCK")
	if !exists {
		t.Skip("`TF_ACC_MASTER_NODE_LOCK` isn't specified - skipping since test will increase test time significantly")
	}

	var node nks.Node
	nodeSize := "standard_f1"
	clusterName := "TerraForm AccTest Master node"
	region := "eastus"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccNKSMasterNode_basic, nodeSize, clusterName, region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNKSMasterNodeExists("nks_master_node.master", &node),
				),
			},
		},
	})
}

func testAccCheckNKSMasterNodeExists(n string, nd *nks.Node) resource.TestCheckFunc {
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
		npID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}
		client := nks.NewClient(os.Getenv("NKS_API_URL"), os.Getenv("NKS_API_URL"))
		node, err := client.GetNode(orgID, clID, npID)
		if err != nil {
			return fmt.Errorf("error occured while fetching nodepool: %s", err)
		}
		nd = node

		return nil
	}
}

const testAccNKSMasterNode_basic = `
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

resource "nks_master_node" "master" {
	org_id               = "${data.nks_organization.org.id}"
	cluster_id           = "${nks_cluster.terraform-cluster.id}"
	provider_code        = "azure"
	platform             = "coreos"
	zone                 = "us-east-2b"
	provider_subnet_cidr = "10.0.1.0/24"
	node_size            = "${data.nks_instance_specs.master-specs.node_size}"
  }

`
