package main

import (
        "github.com/hashicorp/terraform/plugin"
	"github.com/StackPointCloud/terraform-provider-stackpoint/stackpoint"
)

func main() {
        plugin.Serve(&plugin.ServeOpts{
                ProviderFunc: stackpoint.Provider})
}
