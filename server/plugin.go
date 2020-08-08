package main

import (
	"math/rand"
	"sync"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// botID stores the id of our plguin bot
	botID string

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

//LunchbotData contains all data necessary to be stored for the Lunchbot Plugin
type LunchbotData struct {
	LastPairings map[string][]string            `json:"LastPairings"` //Key: UserID, Value: Ordered list of users that this user has been paired with, most recent user is the latest pairing
	UserTopics   map[string]map[string]struct{} `json:"UserTopics"`   //Key: UserID, Value: Set of topics a user is interested in
	Blacklists   map[string]map[string]struct{} `json:"Blacklists"`   //Key: UserID, Value: Set of users that this user has blacklisted
}

//NumHistoryEntries is the number of last pairings per user that get stored in order to avoid pairing with the same users again and again
const NumHistoryEntries int = 50

// OnActivate is invoked when the plugin is activated.
//
// This demo implementation logs a message to the demo channel whenever the plugin is activated.
// It also creates a demo bot account
func (p *Plugin) OnActivate() error {
	//init the rand
	rand.Seed(1337)

	//register all our commands
	if err := p.registerCommands(); err != nil {
		return errors.Wrap(err, "failed to register commands")
	}

	//make sure the bot exists
	botID, ensureBotError := p.Helpers.EnsureBot(&model.Bot{
		Username:    "lunchbot",
		DisplayName: "LunchBot",
		Description: "A bot to find random people to lunch with",
	}, plugin.ProfileImagePath("/assets/together.png"))
	if ensureBotError != nil {
		return errors.Wrap(ensureBotError, "failed to ensure lunchbot.")
	}
	p.botID = botID

	return nil
}

// See https://developers.mattermost.com/extend/plugins/server/reference/
