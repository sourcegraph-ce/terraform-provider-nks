// provider "nks" {
//   /* Set environment variable SPC_API_TOKEN with your API token from StackPointCloud    
//      Set environment variable SPC_BASE_API_URL with API endpoint,   
//      defaults to StackPointCloud production enviroment */
// }
// data "nks_keysets" "keyset_default" {
//   /* You can specify a custom orgID here,   
//      or the system will find and use your default organization ID */
// }
// data "nks_instance_specs" "master-specs" {
//   provider_code = "${var.azure_code}"
//   node_size     = "${var.azure_master_size}"
// }
// data "nks_instance_specs" "worker-specs" {
//   provider_code = "${var.azure_code}"
//   node_size     = "${var.azure_worker_size}"
// }
// resource "nks_cluster" "terraform-cluster" {
//   org_id                            = "${data.nks_keysets.keyset_default.org_id}"
//   cluster_name                      = "Test Azure Cluster TerraForm"
//   provider_code                     = "${var.azure_code}"
//   provider_keyset                   = "${data.nks_keysets.keyset_default.azure_keyset}"
//   region                            = "${var.azure_region}"
//   k8s_version                       = "${var.azure_k8s_version}"
//   startup_master_size               = "${data.nks_instance_specs.master-specs.node_size}"
//   startup_worker_count              = 2
//   startup_worker_size               = "${data.nks_instance_specs.worker-specs.node_size}"
//   provider_resource_group_requested = "${var.azure_resource_group}"
//   rbac_enabled                      = true
//   dashboard_enabled                 = true
//   etcd_type                         = "classic"
//   platform                          = "${var.azure_platform}"
//   channel                           = "stable"
//   ssh_keyset                        = "${data.nks_keysets.keyset_default.user_ssh_keyset}"
// }
// resource "nks_master_node" "master2" {
//   org_id        = "${data.nks_keysets.keyset_default.org_id}"
//   cluster_id    = "${nks_cluster.terraform-cluster.id}"
//   provider_code = "${var.azure_code}"
//   platform      = "${var.azure_platform}"
//   node_size     = "${data.nks_instance_specs.master-specs.node_size}"
// }
// resource "nks_nodepool" "nodepool2" {
//   org_id        = "${data.nks_keysets.keyset_default.org_id}"
//   cluster_id    = "${nks_cluster.terraform-cluster.id}"
//   provider_code = "${var.azure_code}"
//   platform      = "${var.azure_platform}"
//   worker_count  = 2
//   worker_size   = "${data.nks_instance_specs.worker-specs.node_size}"
// }

