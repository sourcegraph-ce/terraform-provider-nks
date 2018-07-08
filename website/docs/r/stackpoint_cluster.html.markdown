---
layout: "stackpoint"
page_title: "StackPoint: stackpoint_cluster"
sidebar_current: "docs-stackpoint-resource-cluster"
description: |-
  Creates and manages a cluster.
---

# stackpoint\_cluster

Creates and manages a cluster in StackPointCloud's system

## Example Usage

```hcl
resource "stackpoint_cluster" "terraform-cluster" {
  org_id                = "${data.stackpoint_keysets.keyset_default.org_id}"
  cluster_name          = "Test AWS Cluster TerraForm"
  provider_code         = "aws"
  provider_keyset       = "${data.stackpoint_keysets.keyset_default.aws_keyset}"
  region                = "us-east-2"
  k8s_version           = "v1.8.3"
  startup_master_size   = "${data.stackpoint_instance_specs.master-specs.node_size}"
  startup_worker_count  = 2
  startup_worker_size   = "${data.stackpoint_instance_specs.worker-specs.node_size}"
  zone                  = "us-east-2a"
  provider_netword_id   = "__new__"
  provider_network_cidr = "10.0.0.0/16"
  provider_subnet_id    = "__new__"
  provider_subnet_cidr  = "10.0.0.0/24"
  rbac_enabled          = true
  dashboard_enabled     = true
  etcd_type             = "classic"
  platform              = "coreos"
  channel               = "stable"
  ssh_keyset            = "${data.stackpoint_keysets.keyset_default.user_ssh_keyset}"
}
```

## Argument reference

* `org_id` - (Required)[int] Organization ID, usually populated by a reference to a keyset datasource value
* `cluster_name` - (Required)[string] Cluster name, can be anything you choose
* `provider_code` - (Required)[string] Cloud provider code string
* `provider_keyset` - (Required)[int] Cloud provider keyset, usually populated by a reference to a keyset datasource value
* `region` - (Required)[string] Cloud provider region where cluster will be built
* `k8s_version` - (Required)[string] Kubernetes version to use for cluster build
* `startup_master_size` - (Required)[string] Master node size, usually populated by a reference to an instance spec value
* `startup_worker_count` - (Required)[int] Number of worker nodes to start the cluster with, minimum is 2
* `startup_worker_size` - (Required)[string] Worker node size, usually populated by a reference to an instance spec value
* `zone` - (Required for AWS)[string] Cloud provider zone where cluster will be built
* `provider_network_id` - (Optional, will default to "__new__")[string] VPC ID to use for cluster network
* `provider_network_cidr` - (Required for AWS/Azure)[string] CIDR of VPC network
* `provider_subnet_id` - (Optional, will default to "__new__")[string] Subnet ID of network to use for cluster
* `provider_subnet_cidr` - (Required)[string] CIDR of Subnet network
* `rbac_enabled` - (Required)[bool] Enable RBAC for cluster
* `dashboard_enabled` - (Required)[bool] Enable Kubernetes dashboard
* `etcd_type` - (Required)[string] Etcd type, classic is recommended
* `platform` - (Required)[string] Operating system of container
* `channel` - (Required)[string] Branch of OS to use
* `ssh_keyset` - (Required)[int] SSH keyset to drop into cluster, usually populated by a reference to a keyset datasource value