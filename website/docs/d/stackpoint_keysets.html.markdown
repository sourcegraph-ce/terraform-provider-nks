---
layout: "stackpoint"
page_title: "StackPoint : stackpoint_keysets"
sidebar_current: "docs-stackpoint-keysets"
description: |-
  Get information on StackPoint keysets
---

# stackpoint\_keysets

The keysets data source can be used to automatically look up your configured cloud provider keysets, based on the API token your used in the provider.  Optionally, you can supply an organization ID as well that will be used.

## Example Usage

```hcl
data "stackpoint_keysets" "keyset_default" {
  org_id  = 111
}
```

## Argument Reference

 * `org_id` - (Optional) Organization ID to use (otherwise the default organization ID is located and used)

## Attributes Reference

 * `aws_keyset` - AWS keyset, if configured, for building resources on AWS infrastructure
 * `do_keyset` - DigitalOcean keyset, if configured, for building resources on DigitalOcean infrastructure
 * `gce_keyset` - GCE keyset, if configured, for building resources on GCE infrastructure
 * `gke_keyset` - GKE keyset, if configured, for building resources on GKE infrastructure
 * `oneandone_keyset` - OneAndOne keyset, if configured, for building resources on OneAndOne infrastructure
 * `org_id` - Organization ID, used for the creation of most resources
 * `packet_keyset` - Packet keyset, if configured, for building resources on Packet infrastructure
 * `user_ssh_keyset` - Your SSH keyset, which will be used on any nodes built
