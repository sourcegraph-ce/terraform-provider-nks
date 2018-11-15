---
layout: "nks"
page_title: "NKS: nks_cluster"
sidebar_current: "docs-nks-resource-cluster"
description: |-
  Creates and manages a cluster.
---

# nks\_cluster

Creates and manages a cluster in NKS's system

## Example Usage

```hcl
resource "nks_cluster" "terraform-cluster" {
  org_id                            = "${data.nks_keysets.keyset_default.org_id}"
  cluster_name                      = "Test AWS Cluster TerraForm"
  provider_code                     = "aws"
  provider_keyset                   = "${data.nks_keysets.keyset_default.aws_keyset}"
  region                            = "us-east-2"
  k8s_version                       = "v1.8.3"
  startup_master_size               = "${data.nks_instance_specs.master-specs.node_size}"
  startup_worker_count              = 2
  startup_worker_size               = "${data.nks_instance_specs.worker-specs.node_size}"
  zone                              = "us-east-2a"
  provider_network_id_requested     = "__new__"
  provider_network_cidr             = "10.0.0.0/16"
  provider_subnet_id_requested      = "__new__"
  provider_subnet_cidr              = "10.0.0.0/24"
  provider_resource_group_requested = "__new__"
  project_id                        = "someproject"
  rbac_enabled                      = true
  dashboard_enabled                 = true
  etcd_type                         = "classic"
  platform                          = "coreos"
  channel                           = "stable"
  ssh_keyset                        = "${data.nks_keysets.keyset_default.user_ssh_keyset}"
}
```

## Argument reference

* `channel` - (Required)[string] Branch of OS to use
* `cluster_name` - (Required)[string] Cluster name, can be anything you choose
* `dashboard_enabled` - (Required)[bool] Enable Kubernetes dashboard
* `etcd_type` - (Required)[string] Etcd type, classic is recommended
* `k8s_version` - (Required)[string] Kubernetes version to use for cluster build
* `org_id` - (Required)[int] Organization ID, usually populated by a reference to a keyset datasource value
* `platform` - (Required)[string] Operating system of container
* `project_id` - (Required for Packet)[string] Packet project ID
* `provider_code` - (Required)[string] Cloud provider code string
* `provider_keyset` - (Required)[int] Cloud provider keyset, usually populated by a reference to a keyset datasource value
* `provider_network_cidr` - (Required for AWS/Azure)[string] CIDR of VPC network
* `provider_network_id_requested` - (Optional, will default to "__new__")[string] VPC ID to use for cluster network
* `provider_resource_group_requested` - (Optional, will default to "__new__")[string] Azure resource group name
* `provider_subnet_cidr` - (Required for AWS/Azure)[string] CIDR of Subnet network
* `provider_subnet_id_requested` - (Optional, will default to "__new__")[string] Subnet ID of network to use for cluster
* `rbac_enabled` - (Required)[bool] Enable RBAC for cluster
* `region` - (Required)[string] Cloud provider region where cluster will be built
* `ssh_keyset` - (Required)[int] SSH keyset to drop into cluster, usually populated by a reference to a keyset datasource value
* `startup_master_size` - (Required)[string] Master node size, usually populated by a reference to an instance spec value
* `startup_worker_count` - (Required)[int] Number of worker nodes to start the cluster with, minimum is 2
* `startup_worker_size` - (Required)[string] Worker node size, usually populated by a reference to an instance spec value
* `zone` - (Required for AWS)[string] Cloud provider zone where cluster will be built

## Attributes Reference

 * `cluster_id` - Cluster ID, used by other resources to reference created cluster
 * `provider_network_id` - VPC ID, might have been newly created if `"__new__"` was used in provider_network_id_requested
 * `provider_subnet_id` - VPC subnet ID, might have been newly created if `"__new__"` was used in provider_subnet_id_requested