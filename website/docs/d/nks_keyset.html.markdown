---
layout: "nks"
page_title: "NKS : nks_keyset"
sidebar_current: "docs-nks-keyset"
description: |-
  Get information on NKS keyset
---

# nks\_keyset

The keysets data source can be used to automatically look up your configured cloud provider keysets, based on the API token your used in the provider.  Optionally, you can supply an organization ID as well that will be used.

## Example Usage

```hcl
data "nks_keyset" "keyset_default" {
  org_id   = 111
  category = "provider"
  entity   = "aws"
}
```

## Argument Reference

 * `category` - (Required) Indicates in which group of keysets to search either "provider" or "user_ssh"
 * `entity` - (Required) If chosen category is 'provider' the 'entity' is required. Options are 'aws', 'azure', 'packet' ...
 * `name` - (Optional) Search by name or part of the name of the keyset. Case insensitive
 * `org_id` - (Optional) Organization ID to use (otherwise the default organization ID is located and used)

## Attributes Reference

 * `id` - ID of the keyset
 