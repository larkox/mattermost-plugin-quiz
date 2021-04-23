package main

import (
	"embed"

	"github.com/mattermost/mattermost-server/v5/plugin"
)

//go:embed static
var staticAssets embed.FS //nolint: gochecknoglobals

func main() {
	plugin.ClientMain(&Plugin{})
}
