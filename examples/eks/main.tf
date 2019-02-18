provider "nks" {
  # Set environment variable NKS_API_TOKEN with your API token from NKS
  # Set environment variable NKS_API_URL with API endpoint,   
  # defaults to NKS production enviroment.
}

# Organization
data "nks_organization" "default" {
  name = "${var.organization_name}"
}

# Keyset
data "nks_keyset" "keyset_provider" {
  # You can specify a custom orgID here or the system will find and use your
  # default organization ID.
  org_id = "${data.nks_organization.default.id}"

  name     = "${var.provider_keyset_name}"
  category = "provider"
  entity   = "${var.provider_code}"
}

data "nks_keyset" "keyset_ssh" {
  # You can specify a custom orgID here or the system will find and use your
  # default organization ID.
  org_id = "${data.nks_organization.default.id}"

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

# Cluster
resource "nks_cluster" "terraform-cluster" {
  org_id                        = "${data.nks_organization.default.id}"
  cluster_name                  = "${var.cluster_name}"
  provider_code                 = "${var.provider_code}"
  provider_keyset               = "${data.nks_keyset.keyset_provider.id}"
  region                        = "${var.provider_region}"
  k8s_version                   = "${var.provider_k8s_version}"
  startup_master_size           = "${data.nks_instance_specs.master-specs.node_size}"
  startup_worker_count          = 2
  startup_worker_size           = "${data.nks_instance_specs.worker-specs.node_size}"
  startup_worker_min_count      = 2
  startup_worker_max_count      = 3
  zone                          = "${var.provider_zone}"
  provider_network_id_requested = "${var.provider_network_id}"
  provider_network_cidr         = "${var.provider_network_cidr}"
  provider_subnet_cidr          = "${var.provider_subnet_cidr}"
  rbac_enabled                  = true
  dashboard_enabled             = true
  etcd_type                     = "${var.provider_etcd_type}"
  platform                      = "${var.provider_platform}"
  channel                       = "${var.provider_channel}"
  ssh_keyset                    = "${data.nks_keyset.keyset_ssh.id}"

  network_component {
    cidr           = "10.0.12.0/24"
    component_type = "provider_subnet"
    id             = "__new__"
    vpc_id         = "vpc-1179c777"
    zone           = "us-east-1a"
    provider_id    = "__new__"
  }

  network_component {
    cidr           = "10.0.15.0/24"
    component_type = "provider_subnet"
    id             = "__new__"
    vpc_id         = "vpc-1179c777"
    zone           = "us-east-1b"
    provider_id    = "__new__"
  }
}
