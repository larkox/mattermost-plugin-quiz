package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/larkox/mattermost-plugin-badges/badgesmodel"
)

func (p *Plugin) EnsureBadges() {
	badges := []badgesmodel.Badge{
		{
			Name:        AchievementNameContentCreator,
			Description: "Create a quiz",
			Image:       p.getStaticURL() + "/contentcreator.png", // TODO use correct url
			ImageType:   badgesmodel.ImageTypeAbsoluteURL,
			Multiple:    false,
		},
		{
			Name:        AchievementNameWinner,
			Description: "Get the highest score in a party game",
			Image:       p.getStaticURL() + "/winner.png", // TODO use correct url
			ImageType:   badgesmodel.ImageTypeAbsoluteURL,
			Multiple:    false,
		},
		{
			Name:        AchievementNameHardWorker,
			Description: "Finish a solo game",
			Image:       p.getStaticURL() + "/hardworker.png", // TODO use correct url
			ImageType:   badgesmodel.ImageTypeAbsoluteURL,
			Multiple:    false,
		},
	}

	reqBody := badgesmodel.EnsureBadgesRequest{
		Badges: badges,
		BotID:  p.BotUserID,
	}
	b, err := json.Marshal(reqBody)
	if err != nil {
		p.mm.Log.Debug("Cannot marshal ensure request", "err", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, badgesmodel.PluginPath+badgesmodel.PluginAPIPath+badgesmodel.PluginAPIPathEnsure, bytes.NewReader(b))
	if err != nil {
		p.mm.Log.Debug("Cannot create http request", "err", err)
		return
	}

	resp := p.mm.Plugin.HTTP(req)
	if resp.StatusCode != http.StatusOK {
		p.mm.Log.Debug("Plugin request failed", "req", req, "resp", resp)
		return
	}

	var newBadges []badgesmodel.Badge
	err = json.NewDecoder(resp.Body).Decode(&newBadges)
	if err != nil {
		p.mm.Log.Debug("Cannot unmarshal response", "err", err)
		return
	}

	p.badgesMap = map[string]badgesmodel.BadgeID{}
	for _, badge := range newBadges {
		p.badgesMap[badge.Name] = badge.ID
	}
}

func (p *Plugin) GrantBadge(name string, userID string) {
	if p.badgesMap == nil {
		p.mm.Log.Debug("No badges map")
		return
	}

	badgeID, ok := p.badgesMap[name]
	if !ok {
		p.mm.Log.Debug("Achievement not recognized")
		return
	}

	grantReq := badgesmodel.GrantBadgeRequest{
		BadgeID: badgeID,
		UserID:  userID,
		BotID:   p.BotUserID,
	}

	b, err := json.Marshal(grantReq)
	if err != nil {
		p.mm.Log.Debug("Cannot marshal grant request")
		return
	}

	req, err := http.NewRequest(http.MethodPost, badgesmodel.PluginPath+badgesmodel.PluginAPIPath+badgesmodel.PluginAPIPathGrant, bytes.NewReader(b))
	if err != nil {
		p.mm.Log.Debug("Cannot create request")
	}

	resp := p.mm.Plugin.HTTP(req)
	if resp.StatusCode != http.StatusOK {
		p.mm.Log.Debug("Plugin request failed", "req", req, "resp", resp)
		return
	}

	p.mm.Log.Debug("Achievement granted", "badgeID", badgeID, "userID", userID, "botID", p.BotUserID)
}
