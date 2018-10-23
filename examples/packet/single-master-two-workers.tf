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
  provider_code = "${var.packet_code}"
  node_size     = "${var.packet_master_size}"
}

data "nks_instance_specs" "worker-specs" {
  provider_code = "${var.packet_code}"
  node_size     = "${var.packet_worker_size}"
}

resource "nks_cluster" "terraform-cluster" {
  cluster_name         = "Test Packet Cluster TerraForm"
  provider_code        = "${var.packet_code}"
  provider_keyset      = "${data.nks_keysets.keyset_default.packet_keyset}"
  region               = "${var.packet_region}"
  k8s_version          = "${var.packet_k8s_version}"
  startup_master_size  = "${data.nks_instance_specs.master-specs.node_size}"
  startup_worker_count = 2
  startup_worker_size  = "${data.nks_instance_specs.worker-specs.node_size}"
  project_id           = "YOUR PROJECT ID HERE"
  rbac_enabled         = true
  dashboard_enabled    = true
  etcd_type            = "classic"
  platform             = "${var.packet_platform}"
  channel              = "stable"
  ssh_keyset           = "${data.nks_keysets.keyset_default.user_ssh_keyset}"
}
