---
layout: "nks"
page_title: "NKS: nks_nodepool"
sidebar_current: "docs-nks-resource-nodepool"
description: |-
  Creates and manages an additional nodepool of worker nodes.
---

# nks\_nodepool

Creates and manages an additional nodepool of worker nodes in NKS's system

## Example Usage

```hcl
resource "nks_nodepool" "nodepool2" {
  org_id               = "${data.nks_keysets.keyset_default.org_id}"
  cluster_id           = "${nks_cluster.terraform-cluster.id}"
  provider_code        = "aws"
  platform             = "coreos"
  zone                 = "us-east-2b"
  provider_subnet_cidr = "10.2.0.0/24"
  worker_count         = 1
  worker_size          = "${data.nks_instance_specs.worker-specs.node_size}"
}
```

## Argument reference

* `cluster_id` - (Required)[int] Cluster ID, usually populated by a reference to a cluster resource value
* `org_id` - (Required)[int] Organization ID, usually populated by a reference to a keyset datasource value
* `platform` - (Required)[string] Operating system of container
* `provider_code` - (Required)[string] Cloud provider code string
* `provider_subnet_cidr` - (Required for AWS/Azure)[string] CIDR of Subnet network
* `worker_count` - (Required)[int] Number of nodes to build in new nodepool
* `worker_size` - (Required)[string] Node size, usually populated by a reference to an instance spec value
* `zone` - (Required for AWS)[string] Cloud provider zone where cluster will be built
