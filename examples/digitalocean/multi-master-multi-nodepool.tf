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
  provider_code = "${var.digitalocean_code}"
  node_size     = "${var.digitalocean_master_size}"
}

data "nks_instance_specs" "worker-specs" {
  provider_code = "${var.digitalocean_code}"
  node_size     = "${var.digitalocean_worker_size}"
}

resource "nks_cluster" "terraform-cluster" {
  org_id               = "${data.nks_keysets.keyset_default.org_id}"
  cluster_name         = "Test DO Cluster TerraForm"
  provider_code        = "${var.digitalocean_code}"
  provider_keyset      = "${data.nks_keysets.keyset_default.do_keyset}"
  region               = "${var.digitalocean_region}"
  k8s_version          = "${var.digitalocean_k8s_version}"
  startup_master_size  = "${data.nks_instance_specs.master-specs.node_size}"
  startup_worker_count = 2
  startup_worker_size  = "${data.nks_instance_specs.worker-specs.node_size}"
  rbac_enabled         = true
  dashboard_enabled    = true
  etcd_type            = "classic"
  platform             = "${var.digitalocean_platform}"
  channel              = "stable"
  ssh_keyset           = "${data.nks_keysets.keyset_default.user_ssh_keyset}"
}

resource "nks_master_node" "master2" {
  org_id        = "${data.nks_keysets.keyset_default.org_id}"
  cluster_id    = "${nks_cluster.terraform-cluster.id}"
  provider_code = "${var.digitalocean_code}"
  platform      = "${var.digitalocean_platform}"
  node_size     = "${data.nks_instance_specs.master-specs.node_size}"
}

resource "nks_nodepool" "nodepool2" {
  org_id        = "${data.nks_keysets.keyset_default.org_id}"
  cluster_id    = "${nks_cluster.terraform-cluster.id}"
  provider_code = "${var.digitalocean_code}"
  platform      = "${var.digitalocean_platform}"
  worker_count  = 1
  worker_size   = "${data.nks_instance_specs.worker-specs.node_size}"
}
