provider "nks" {
  /* Set environment variable SPC_API_TOKEN with your API token from StackPointCloud    
     Set environment variable SPC_BASE_API_URL with API endpoint,   
     defaults to StackPointCloud production enviroment */
}

data "nks_organization" "org" {}

data "nks_keyset" "keyset_default" {
  category = "provider"
  entity   = "azure"

  /* You can specify a custom orgID here,   
     or the system will find and use your default organization ID */
}

data "nks_keyset" "ssh_key" {
  category = "user"
  name     = "default"

  /* You can specify a custom orgID here,   
     or the system will find and use your default organization ID */
}

data "nks_instance_specs" "master-specs" {
  provider_code = "${var.azure_code}"
  node_size     = "${var.azure_master_size}"
}

data "nks_instance_specs" "worker-specs" {
  provider_code = "${var.azure_code}"
  node_size     = "${var.azure_worker_size}"
}

resource "nks_cluster" "terraform-cluster" {
  org_id                            = "${data.nks_organization.org.id}"
  cluster_name                      = "Test Azure Cluster TerraForm"
  provider_code                     = "${var.azure_code}"
  provider_keyset                   = "${data.nks_keyset.keyset_default.id}"
  region                            = "${var.azure_region}"
  k8s_version                       = "${var.azure_k8s_version}"
  startup_master_size               = "${data.nks_instance_specs.master-specs.node_size}"
  startup_worker_count              = 2
  startup_worker_size               = "${data.nks_instance_specs.worker-specs.node_size}"
  provider_resource_group_requested = "${var.azure_resource_group}"
  rbac_enabled                      = true
  dashboard_enabled                 = true
  etcd_type                         = "classic"
  platform                          = "${var.azure_platform}"
  channel                           = "stable"
  ssh_keyset                        = "${data.nks_keyset.ssh_key.id}"
}

resource "nks_solution" "efk" {
  org_id     = "${data.nks_organization.org.id}"
  cluster_id = "${nks_cluster.terraform-cluster.id}"
  solution   = "efk"
}
