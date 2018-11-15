---
layout: "nks"
page_title: "Provider: NKS"
sidebar_current: "docs-nks-index"
description: |-
  A provider for NKS.
---

# NKS Provider

The NKS provider gives the ability to deploy and configure resources using the NKS API.

Use the navigation to the left to read about the available data sources and resources.


## Usage

The provider needs to be configured with proper credentials before it can be used.


```hcl
$ export NKS_API_TOKEN="nks_api_token"
$ export NKS_API_URL="nks_api_url"
```

Or you can provide your credentials in a `.tf` configuration file as shown in this example.


## Example Usage


```hcl
provider "nks" {
  token    = "nks_api_token"
  endpoint = "nks_api_url"
}

data "nks_keysets" "keyset_default" {
  /* You can specify a custom orgID here,   
     or the system will find and use your default organization ID */
}
```


**Note**: The credentials provided in a `.tf` file will override the credentials from environment variables.

## Configuration Reference

The following arguments are supported:

* `token` - (Required) If omitted, the `NKS_API_TOKEN` environment variable is used.

* `endpoint` - (Optional) If omitted, the `NKS_API_URL` environment variable is used, or it defaults to the current production API URL.


## Resource Timeout

Individual resources may provide a `timeout` value to configure the amount of time a specific operation is allowed to take before being considered an error. Each resource may provide a configurable timeout, measured in seconds. Each resource that supports timeouts will have a default value for that operation.
Users can overwrite the default values for a specific resource in the configuration.

The default `timeouts` values are:

* create  - (Default 30 minutes) Used for creating a resource.
* update  - (Default 30 minutes) Used for updating a resource .
* delete  - (Default 30 minutes) Used for destroying a resource.
* default - (Default 30 minutes) Used for every other action on a resource.

An example of overwriting the default timeout (setting timeout to 6 minutes instead of the default 30 minutes):

```hcl
resource "nks_master_node" "master2" {
  org_id        = "111"
  cluster_id    = "4066"
  provider_code = "azure"
  node_size     = "standard_f1"
  timeout       = 3600
}

```

~> **Note:** Terraform does not automatically rollback in the face of errors.
Instead, your Terraform state file will be partially updated with
any resources that successfully completed.

## Support
You are welcome to contact us with questions or comments at: questions@stackpointcloud.com
