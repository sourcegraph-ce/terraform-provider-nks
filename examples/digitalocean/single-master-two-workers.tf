provider "stackpoint" {
  /* Set environment variable SPC_API_TOKEN with your API token from StackPointCloud    
     Set environment variable SPC_BASE_API_URL with API endpoint,   
     defaults to StackPointCloud production enviroment */
}

data "stackpoint_keysets" "keyset_default" {
  /* You can specify a custom orgID here,   
     or the system will find and use your default organization ID */
}

data "stackpoint_instance_specs" "master-specs" {
  provider_code = "${var.digitalocean_code}"
  node_size     = "${var.digitalocean_master_size}"
}

data "stackpoint_instance_specs" "worker-specs" {
  provider_code = "${var.digitalocean_code}"
  node_size     = "${var.digitalocean_worker_size}"
}

resource "stackpoint_cluster" "terraform-cluster" {
  org_id               = "${data.stackpoint_keysets.keyset_default.org_id}"
  cluster_name         = "Test DO Cluster TerraForm"
  provider_code        = "${var.digitalocean_code}"
  provider_keyset      = "${data.stackpoint_keysets.keyset_default.do_keyset}"
  region               = "${var.digitalocean_region}"
  k8s_version          = "${var.digitalocean_k8s_version}"
  startup_master_size  = "${data.stackpoint_instance_specs.master-specs.node_size}"
  startup_worker_count = 2
  startup_worker_size  = "${data.stackpoint_instance_specs.worker-specs.node_size}"
  rbac_enabled         = true
  dashboard_enabled    = true
  etcd_type            = "classic"
  platform             = "${var.digitalocean_platform}"
  channel              = "stable"
  ssh_keyset           = "${data.stackpoint_keysets.keyset_default.user_ssh_keyset}"
}

resource "stackpoint_master_node" "master2" {
  org_id        = "${data.stackpoint_keysets.keyset_default.org_id}"
  cluster_id    = "${stackpoint_cluster.terraform-cluster.id}"
  provider_code = "${var.digitalocean_code}"
  platform      = "${var.digitalocean_platform}"
  node_size     = "${data.stackpoint_instance_specs.master-specs.node_size}"
}
