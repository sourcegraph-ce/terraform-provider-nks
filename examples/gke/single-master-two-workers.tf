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
  provider_code = "${var.gke_code}"
  node_size     = "${var.gke_master_size}"
}

data "stackpoint_instance_specs" "worker-specs" {
  provider_code = "${var.gke_code}"
  node_size     = "${var.gke_worker_size}"
}

resource "stackpoint_cluster" "terraform-do-cluster" {
  org_id               = "${data.stackpoint_keysets.keyset_default.org_id}"
  cluster_name         = "Test GKE Cluster TerraForm"
  provider_code        = "${var.gke_code}"
  provider_keyset      = "${data.stackpoint_keysets.keyset_default.gke_keyset}"
  region               = "${var.gke_region}"
  k8s_version          = "${var.gke_k8s_version}"
  startup_master_size  = "${data.stackpoint_instance_specs.master-specs.node_size}"
  startup_worker_count = 2
  startup_worker_size  = "${data.stackpoint_instance_specs.worker-specs.node_size}"
  region               = "${var.gke_region}"
  rbac_enabled         = true
  dashboard_enabled    = true
  etcd_type            = "classic"
  platform             = "${var.gke_platform}"
  channel              = "stable"
  ssh_keyset           = "${data.stackpoint_keysets.keyset_default.user_ssh_keyset}"
}
