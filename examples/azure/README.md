# Azure Cluster Example

This example will show you how to create a cluster on Azure with NetApp Kubernetes Service.

This example does the following:

- Finds an organization
- Finds an Azure keyset
- Finds an SSH keyset
- Creates a cluster

[Keyset examples](/examples/keysets#adding-a-cloud-provider-keyset-for-azure) shows how to add a keyset to NKS.

## Run the example

From inside of this directory:

```bash
export NKS_API_TOKEN=<this is a secret>
export NKS_API_URL=https://api.stackpoint.io/
terraform init
terraform plan
terraform apply
```

## Remove the example

```bash
terraform destroy
