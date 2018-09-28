---
layout: "stackpoint"
page_title: "StackPoint : stackpoint_instances"
sidebar_current: "docs-stackpoint-datasource-instances"
description: |-
  Get information on a StackPoint Instances
---

# stackpoint\_instances

The instance specs data source can be used to select and validate a node size. You can provide a string for the name of the instance size and our system will validate that it is a valid instance size for nodes on that cloud provider.

## Example Usage

```hcl
data "stackpoint_instance_specs" "master-specs" {
  provider_code = "azure"
  node_size     = "standard_f1"
}
```

## Argument Reference

 * `provider_code` - (Required) Short name for the cloud provider you wish to interact with.
 * `node_size`     - (Required) Name of the instance size you wish to use for a node.
 * `endpoint`      - (Optional) You may override the endpoint used by the provider here, if you are using an endpoint that doesn't validate instance sizes.

## Attributes Reference

 * `node_size` - Instance size string, used for cluster and node creation