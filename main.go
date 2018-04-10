package main

import (
	"github.com/MissionCriticalCloud/terraform-provider-cosmic/cosmic"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: cosmic.Provider,
	})
}
