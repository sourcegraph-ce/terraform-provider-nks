package main

import (
	"github.com/StackPointCloud/terraform-provider-nks/stackpoint"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: stackpoint.Provider})
}
