---
layout: "stackpoint"
page_title: "StackPoint: stackpoint_masternode"
sidebar_current: "docs-stackpoint-resource-masternode"
description: |-
  Creates and manages an additional master node.
---

# stackpoint\_masternode

Creates and manages an additional master node in StackPointCloud's system for high-availability

## Example Usage

```hcl
resource "stackpoint_master_node" "master2" {
  org_id               = "${data.stackpoint_keysets.keyset_default.org_id}"
  cluster_id           = "${stackpoint_cluster.terraform-cluster.id}"
  provider_code        = "aws"
  platform             = "coreos"
  zone                 = "us-east-2b"
  provider_subnet_cidr = "10.1.0.0/24"
  node_size            = "${data.stackpoint_instance_specs.master-specs.node_size}"
}
```

## Argument reference

* `cluster_id` - (Required)[int] Cluster ID, usually populated by a reference to a cluster resource value
* `node_size` - (Required)[string] Node size, usually populated by a reference to an instance spec value
* `org_id` - (Required)[int] Organization ID, usually populated by a reference to a keyset datasource value
* `platform` - (Required)[string] Operating system of container
* `provider_code` - (Required)[string] Cloud provider code string
* `provider_subnet_cidr` - (Required for AWS/Azure)[string] CIDR of Subnet network
* `zone` - (Required for AWS)[string] Cloud provider zone where cluster will be built
