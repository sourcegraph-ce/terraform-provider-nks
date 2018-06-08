Terraform StackPointCloud Provider
============================

- Website: https://www.stackpoint.io

Maintainers
-----------

This provider plugin is maintained by:

* Justin Hopper ([@justinhopper](https://github.com/justinhopper))

Requirements
------------

-       [Terraform](https://www.terraform.io/downloads.html) 0.10.x
-       [Go](https://golang.org/doc/install) 1.9 (to build the provider plugin)
-	[StackPointCloud Account](http://www.stackpoint.io)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/StackPointCloud/terraform-provider-stackpoint`

```sh
$ mkdir -p $GOPATH/src/github.com/StackPointCloud; cd $GOPATH/src/github.com/StackPointCloud
$ git clone https://github.com/StackPointCloud/terraform-provider-stackpoint
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/StackPointCloud/terraform-provider-stackpoint
$ make build
```

Using The Provider
------------------

The stackpoint provider plugin comes with several example configuration files for working with different cloud platforms that include:

- 1&1
- AWS
- Azure
- Digital Ocean
- Google Compute Engine
- Google Kontainer Engine
- Packet

To use one of the example configuration files, e.g. AWS, copy the example file to TerraForm's main config file:
```sh
$ cp main.tf.aws_example main.tf
```

All of the example configuration files use values in the variables.tf file. Edit this file and replace `# YOUR ID HERE` with your ID key from the StackPointCloud system for that provider. If you plan to only use one provider, such as AWS, you only need to insert your ID key for AWS and can leave the rest as they are.

The default configuration in each example will:

- Load the stackpoint provider plugin (be sure to insert your Organization ID and SSH Keyset ID from the StackPointCloud system)
- Configure a master node size (defaulted to look in the variables.tf file)
- Configure a worker node size (defaulted to look in the variables.tf file)
- Configure a cluster with one master and two workers at the cloud platform specified (most of these options can be left at default settings, or you can customize as you see fit)

The next two blocks are commented out. The first allows for adding a second master node, making the cluster HA (highly available). The second allows for adding an additional worker node pool.
