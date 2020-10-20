package main

import (
	"github.com/MissionCriticalCloud/terraform-provider-cosmic/cosmic"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: cosmic.New,
	})
}
