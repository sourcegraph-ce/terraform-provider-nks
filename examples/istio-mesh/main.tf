provider "nks" {
  # Set environment variable NKS_API_TOKEN with your API token from NKS
  # Set environment variable NKS_API_URL with API endpoint,   
  # defaults to NKS production enviroment.
}

# Organization
data "nks_organization" "default" {
  name = "${var.organization_name}"
}

# Keysets
data "nks_keyset" "keyset_default" {
  # You can specify a custom orgID here or the system will find and use your
  # default organization ID.
  org_id   = "${data.nks_organization.default.id}"
  name     = "${var.provider_keyset_name}"
  category = "provider"
  entity   = "${var.provider_code}"
}

data "nks_keyset" "keyset_ssh" {
  # You can specify a custom orgID here or the system will find and use your
  # default organization ID.
  org_id   = "${data.nks_organization.default.id}"
  category = "user_ssh"
  name     = "${var.ssh_keyset_name}"
}

# Instance specs
data "nks_instance_specs" "master-specs" {
  provider_code = "${var.provider_code}"
  node_size     = "${var.provider_master_size}"
}

data "nks_instance_specs" "worker-specs" {
  provider_code = "${var.provider_code}"
  node_size     = "${var.provider_worker_size}"
}

# Clusters
resource "nks_cluster" "terraform-cluster-a" {
  org_id                            = "${data.nks_organization.default.id}"
  cluster_name                      = "${var.a_cluster_name}"
  provider_code                     = "${var.provider_code}"
  provider_keyset                   = "${data.nks_keyset.keyset_default.id}"
  region                            = "${var.provider_region}"
  k8s_version                       = "${var.provider_k8s_version}"
  startup_master_size               = "${data.nks_instance_specs.master-specs.node_size}"
  startup_worker_count              = 2
  startup_worker_size               = "${data.nks_instance_specs.worker-specs.node_size}"
  provider_resource_group_requested = "${var.provider_resource_group}"
  rbac_enabled                      = true
  dashboard_enabled                 = true
  etcd_type                         = "${var.provider_etcd_type}"
  platform                          = "${var.provider_platform}"
  channel                           = "${var.provider_channel}"
  ssh_keyset                        = "${data.nks_keyset.keyset_ssh.id}"
}

resource "nks_cluster" "terraform-cluster-b" {
  org_id                            = "${data.nks_organization.default.id}"
  cluster_name                      = "${var.b_cluster_name}"
  provider_code                     = "${var.provider_code}"
  provider_keyset                   = "${data.nks_keyset.keyset_default.id}"
  region                            = "${var.provider_region}"
  k8s_version                       = "${var.provider_k8s_version}"
  startup_master_size               = "${data.nks_instance_specs.master-specs.node_size}"
  startup_worker_count              = 2
  startup_worker_size               = "${data.nks_instance_specs.worker-specs.node_size}"
  provider_resource_group_requested = "${var.provider_resource_group}"
  rbac_enabled                      = true
  dashboard_enabled                 = true
  etcd_type                         = "${var.provider_etcd_type}"
  platform                          = "${var.provider_platform}"
  channel                           = "${var.provider_channel}"
  ssh_keyset                        = "${data.nks_keyset.keyset_ssh.id}"
}

# Solutions
resource "nks_solution" "istio-a" {
	org_id     = "${data.nks_organization.default.id}"
	cluster_id = "${nks_cluster.terraform-cluster-a.id}"
	solution   = "istio"
}

resource "nks_solution" "istio-b" {
	org_id     = "${data.nks_organization.default.id}"
	cluster_id = "${nks_cluster.terraform-cluster-b.id}"
	solution   = "istio"
}

# Workspace
data "nks_workspace" "my-workspace" {
	org_id = "${data.nks_organization.org.id}"
}

# Istio mesh
resource "nks_istio_mesh" "terraform-istio-mesh" {
	name        = "${var.istio_mesh_name}"
	mesh_type   = "${var.istio_mesh_type}"
	org_id      = "${data.nks_organization.default.id}"
	workspace	  = "${data.nks_workspace.my-workspace.id}"
	members		  = [
		{
			cluster	= "${nks_cluster.terraform-cluster-a.id}"
			role	  = "host"
      istio_solution_id = "${nks_solution.istio-a.id}"
		},
		{
			cluster	= "${nks_cluster.terraform-cluster-b.id}"
			role	  = "guest"
      istio_solution_id = "${nks_solution.istio-b.id}"
		}
	]
}
