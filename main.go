package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: servicenow.Provider,
	})
}
