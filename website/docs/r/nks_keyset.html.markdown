---
layout: "nks"
page_title: "NKS: nks_keyset"
sidebar_current: "docs-nks-resource-keyset"
description: |-
  Installs and manages a keyset in a cluster
---

# nks\_solution

Installs and manages a solution in a cluster in NKS's system

## Example Usage

```hcl
resource "nks_solution" "aws_keyset" {
  org_id               = "${data.nks_keysets.keyset_default.org_id}"
  name                 = "AWS Keyset"
  category             = "provider"
  entity               = "aws"
  workspaces           = ""
  keys = [
      {
        key_type         = "pub"
        key              = "${var.aws_access_key}"
      },
      {
        key_type         = "pvt"
        key              = "${var.aws_secret_key}"
      }
  ]  
}
```

## Argument reference

* `org_id` - (Required)[int] Organization ID, usually populated by a reference to a keyset datasource value
* `name` - (Required)[string] Keyset name, can be anything you choose
* `category` - (Required)[string] Category, category valid values are provider, user_ssh and solution
* `entity` - (Optional)[string] Entity, usually populated by a reference to a keyset datasource value
* `keys` - (Required)[list] Keys, a list of `key` related to a keyset
* `workspaces` - (Optional)[list] Workspaces, A list of workspace ID

`key` supports the following:
* `key_type` - (Optional)[string] Keytype, represents type of key. Examples of key type commonly used are `pub`, `prv`, `tenant`, `subscription`, `license`, `pull_secret`, `token`, `username`, `password`, `scope`, `other` etc.
* `key` - (Optional)[string] Key, represents a value for specific key type