package plugin

import (
	"github.com/cloudquery/cloudquery/plugins/source/0xdb/client"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
)

var Version = "Development"

func Plugin() *plugin.Plugin {
	return plugin.NewPlugin(
		"0xdb",
		Version,
		client.Configure,
	)
}
