NKS Provider for Terraform
==========================

NetApp Kubernetes Service (NKS) is a universal control plane for creating and managing Kubernetes clusters.

- Website: https://cloud.netapp.com/kubernetes-service

Requirements
------------

- [NKS Account](https://cloud.netapp.com/kubernetes-service)
- [Terraform](https://www.terraform.io/downloads.html) 0.10.x
- [Go](https://golang.org/doc/install) 1.9 (to build the provider plugin)

Building the Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/NetApp/terraform-provider-nks`

```sh
$ mkdir -p $GOPATH/src/github.com/NetApp; cd $GOPATH/src/github.com/NetApp
$ git clone git@github.com:NetApp/terraform-provider-nks
```

Enter the provider directory and build the provider.

```sh
$ cd $GOPATH/src/github.com/NetApp/terraform-provider-nks
$ make build
```

Using the Provider
------------------
If you're building the provider, follow the instructions to [install it as a plugin](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin). After placing the binary into your plugins directory, run `terraform init` to initialize it.

Examples
--------

The examples cover basic scenarios for using NetApp Kubernetes Services. Each example has a related README with more details on the topic. The following examples are availabe:

- [Add a keyset](examples/keysets)
- [Create an AKS cluster](examples/aks)
- [Create an AWS cluster](examples/aws)
- [Create an Azure cluster](examples/azure)
- [Create an EKS cluster](examples/eks)
- [Create a GCE cluster](examples/gce)
- [Create a GKE cluster](examples/gke)
- [Create a cluster and install solutions from the gallery](examples/solutions)

Developing the Provider
-----------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.9+ is *required*). You will also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-nks
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources and often cost money to run.

```sh
$ make testacc
```

If you need to add a new package in the vendor directory under `github.com/NetApp/nks-sdk-go`, create a separate PR handling _only_ the update of the vendor for your new requirement. Make sure to pin your dependency to a specific version, and that all versions of `github.com/NetApp/nks-sdk-go/*` are pinned to the same version.

Terraform Resources
-------------------

- Website: https://www.terraform.io/
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)
