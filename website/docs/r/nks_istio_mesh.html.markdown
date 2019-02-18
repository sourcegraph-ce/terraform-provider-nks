---
layout: "nks"
page_title: "NKS: nks_istio_mesh"
sidebar_current: "docs-nks-resource-istio-mesh"
description: |-
  Installs and manages an istio mesh in an organization
---

# nks\_solution

Installs and manages an istio mesh in an organization in NKS's system

## Example Usage

```hcl
resource "nks_istio_mesh" "my_istio_mesh" {
  name                    = "My Istio Mesh"
  mesh_type               = "cross_cluster"
  workspace               = "${data.nks_workspace.default.id}"
  org_id                  = "${data.nks_keysets.keyset_default.org_id}"
  members = [
        {
            cluster            = "${data.nks_cluster.cluster1.id}"
            role               = "host"
            istio_solution_id  = "${nks_solution.istio-a.id}"
        },
        {
            cluster            = "${data.nks_cluster.cluster2.id}"
            role               = "guest"
            istio_solution_id  = "${nks_solution.istio-b.id}"
        },
  ]
}
```

## Argument reference

* `name` - (Required)[string] Keyset name, can be anything you choose
* `mesh_type` - (Required)[string] Mesh type, represents type of mesh and the following type is available `cross_cluster`
* `workspace` - (Required)[int] Workspace ID, usually populated by a refernece to a workspace resource value
* `org_id` - (Required)[int] Organization ID, usually populated by a reference to a keyset datasource value
* `members` - (Required)[list] Members, a list of `member`

`member` supports the following:
* `role` - (Required)[string] Role, represents type of role and the following types are available `host` and `guest`.
* `cluster` - (Required)[int] Cluster ID, usually populated by a reference to a cluster resource value
* `istio_solution_id` - (Optional)[int] Istio solution ID, usually populated by a reference to istio solution resource value
  previously installed on the same istio mesh member cluster. Terraform parameter `depends_on` can be used as an alternative.