package main

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/larkox/mattermost-plugin-badges/badgesmodel"
	pluginapi "github.com/mattermost/mattermost-plugin-api"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	mm        *pluginapi.Client
	BotUserID string
	store     Store
	router    *mux.Router
	badgesMap map[string]badgesmodel.BadgeID
}

// ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	r.Header.Add("Mattermost-Plugin-ID", c.SourcePluginId)
	w.Header().Set("Content-Type", "application/json")

	p.router.ServeHTTP(w, r)
}

func (p *Plugin) OnActivate() error {
	p.mm = pluginapi.NewClient(p.API)

	botID, err := p.mm.Bot.EnsureBot(&model.Bot{
		Username:    BotUserName,
		DisplayName: BotDisplayName,
		Description: BotDescription,
	})
	if err != nil {
		return errors.Wrap(err, "failed to ensure quiz bot")
	}
	p.BotUserID = botID
	p.store = NewStore(p.mm)
	p.initializeAPI()

	p.EnsureBadges()
	return p.mm.SlashCommand.Register(p.getCommand())
}
