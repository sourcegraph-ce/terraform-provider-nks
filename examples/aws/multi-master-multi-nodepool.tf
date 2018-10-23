provider "nks" {
  /* Set environment variable SPC_API_TOKEN with your API token from StackPointCloud    
     Set environment variable SPC_BASE_API_URL with API endpoint,   
     defaults to StackPointCloud production enviroment */
}

data "nks_keysets" "keyset_default" {
  /* You can specify a custom orgID here,   
     or the system will find and use your default organization ID */
}

data "nks_instance_specs" "master-specs" {
  provider_code = "${var.aws_code}"
  node_size     = "${var.aws_master_size}"
}

data "nks_instance_specs" "worker-specs" {
  provider_code = "${var.aws_code}"
  node_size     = "${var.aws_worker_size}"
}

resource "nks_cluster" "terraform-cluster" {
  org_id                = "${data.nks_keysets.keyset_default.org_id}"
  cluster_name          = "Test AWS Cluster TerraForm"
  provider_code         = "${var.aws_code}"
  provider_keyset       = "${data.nks_keysets.keyset_default.aws_keyset}"
  region                = "${var.aws_region}"
  k8s_version           = "${var.aws_k8s_version}"
  startup_master_size   = "${data.nks_instance_specs.master-specs.node_size}"
  startup_worker_count  = 2
  startup_worker_size   = "${data.nks_instance_specs.worker-specs.node_size}"
  zone                  = "${var.aws_zone}"
  provider_network_cidr = "${var.aws_network_cidr}"
  provider_subnet_cidr  = "${var.aws_subnet_cidr}"
  rbac_enabled          = true
  dashboard_enabled     = true
  etcd_type             = "classic"
  platform              = "${var.aws_platform}"
  channel               = "stable"
  ssh_keyset            = "${data.nks_keysets.keyset_default.user_ssh_keyset}"
}

resource "nks_master_node" "master2" {
  org_id               = "${data.nks_keysets.keyset_default.org_id}"
  cluster_id           = "${nks_cluster.terraform-cluster.id}"
  provider_code        = "${var.aws_code}"
  platform             = "${var.aws_platform}"
  zone                 = "${var.aws_zone2}"
  provider_subnet_cidr = "${var.aws_subnet_cidr2}"
  node_size            = "${data.nks_instance_specs.master-specs.node_size}"
}

resource "nks_nodepool" "nodepool2" {
  org_id               = "${data.nks_keysets.keyset_default.org_id}"
  cluster_id           = "${nks_cluster.terraform-cluster.id}"
  provider_code        = "${var.aws_code}"
  platform             = "${var.aws_platform}"
  zone                 = "${var.aws_zone2}"
  provider_subnet_cidr = "${nks_master_node.master2.provider_subnet_cidr}"
  worker_count         = 1
  worker_size          = "${data.nks_instance_specs.worker-specs.node_size}"
}
