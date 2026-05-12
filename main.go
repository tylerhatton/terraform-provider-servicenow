package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/tylerhatton/terraform-provider-servicenow/servicenow"
)

// Version is the provider version, set at build time via -ldflags "-X main.Version=...".
var Version = "dev"

func main() {
	servicenow.SetVersion(Version)
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: servicenow.Provider,
	})
}
