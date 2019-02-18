---
layout: "nks"
page_title: "NKS: nks_workspace"
sidebar_current: "docs-nks-resource-workspace"
description: |-
  Installs and manages a workspace in an organization
---

# nks\_workspace

Installs and manages a workspace in an organization in NKS's system

## Example Usage

```hcl
resource "nks_solution" "my_workspace" {
  name                 = "My workspace"
  org_id               = "${data.nks_keysets.keyset_default.org_id}"
  slug                 = "my_workspace"
  default              = "false"
}
```

## Argument reference

* `name` - (Required)[string] Workspace name, can be anything you choose
* `org_id` - (Required)[int] Organization ID, usually populated by a reference to a keyset datasource value
* `slug` - (Optional)[string] Slug, can be anything you choose
* `default` - (Optional)[bool] Default, can be default workspace