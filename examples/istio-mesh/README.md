# Istio mesh example

This example will show you how to create an istio mesh containing two clusters on Azure (with istio solution installed) with NetApp Kubernetes Service.

This example does the following:

- Finds an organization
- Finds an Azure keyset
- Finds an SSH keyset
- Creates two clusters
- Installs istio solution on both of them
- Creates an istio mesh with one of the clusters having a host role, and the other a guest role.

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
