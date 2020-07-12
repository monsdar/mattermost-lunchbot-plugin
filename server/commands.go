package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
)

const (
	commandLunchbot                = "lunchbot"
	commandLunchbotBlacklistShow   = commandLunchbot + " blacklist show"
	commandLunchbotBlacklistAdd    = commandLunchbot + " blacklist add"
	commandLunchbotBlacklistRemove = commandLunchbot + " blacklist remove"
	commandLunchbotTopicsShow      = commandLunchbot + " topics show"
	commandLunchbotTopicsAdd       = commandLunchbot + " topics add"
	commandLunchbotTopicsRemove    = commandLunchbot + " topics remove"
)

func (p *Plugin) registerCommands() error {
	commands := [...]model.Command{
		model.Command{
			Trigger:          commandLunchbot,
			AutoComplete:     true,
			AutoCompleteDesc: "Pairs you with a random user",
		},
		model.Command{
			Trigger:          commandLunchbotBlacklistShow,
			AutoComplete:     true,
			AutoCompleteDesc: "Your blacklist is a list of users you do not want to get paired with",
		},
		model.Command{
			Trigger:          commandLunchbotBlacklistAdd,
			AutoComplete:     true,
			AutoCompleteHint: "<username>",
			AutoCompleteDesc: "Add someone to your blacklist by his username",
		},
		model.Command{
			Trigger:          commandLunchbotBlacklistRemove,
			AutoComplete:     true,
			AutoCompleteHint: "<username>",
			AutoCompleteDesc: "Remove someone from your blacklist by his username",
		},
		model.Command{
			Trigger:          commandLunchbotTopicsShow,
			AutoComplete:     true,
			AutoCompleteDesc: "Topics are things you'd like to talk about",
		},
		model.Command{
			Trigger:          commandLunchbotTopicsAdd,
			AutoComplete:     true,
			AutoCompleteHint: "<topic>",
			AutoCompleteDesc: "Add a topic to your list",
		},
		model.Command{
			Trigger:          commandLunchbotTopicsRemove,
			AutoComplete:     true,
			AutoCompleteHint: "<topic>",
			AutoCompleteDesc: "Remove a topic from your list",
		},
	}

	for _, command := range commands {
		if err := p.API.RegisterCommand(&command); err != nil {
			return errors.Wrapf(err, fmt.Sprintf("Failed to register %s command", command.Trigger))
		}
	}

	return nil
}

// ExecuteCommand executes a command that has been previously registered via the RegisterCommand
// API.
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	userCommands := map[string]func(args *model.CommandArgs) (*model.CommandResponse, *model.AppError){
		commandLunchbotBlacklistShow: func(args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
			return p.executeCommandLunchbotBlacklistShow(args), nil
		},
		commandLunchbotBlacklistAdd: func(args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
			return p.executeCommandLunchbotBlacklistAdd(args), nil
		},
		commandLunchbotBlacklistRemove: func(args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
			return p.executeCommandLunchbotBlacklistRemove(args), nil
		},
		commandLunchbotTopicsShow: func(args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
			return p.executeCommandLunchbotTopicsShow(args), nil
		},
		commandLunchbotTopicsAdd: func(args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
			return p.executeCommandLunchbotTopicsAdd(args), nil
		},
		commandLunchbotTopicsRemove: func(args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
			return p.executeCommandLunchbotTopicsRemove(args), nil
		},
	}

	mainCommand := map[string]func(args *model.CommandArgs) (*model.CommandResponse, *model.AppError){
		//this needs to be last, as prefix `/lunchbot` is also part of the above commands
		commandLunchbot: func(args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
			return p.executeCommandLunchbot(args), nil
		},
	}

	trigger := strings.TrimPrefix(args.Command, "/")
	for key, value := range userCommands {
		if strings.HasPrefix(trigger, key) {
			return value(args)
		}
	}
	for key, value := range mainCommand {
		if strings.HasPrefix(trigger, key) {
			return value(args)
		}
	}

	//return an error message when the command has not been detected at all
	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         fmt.Sprintf("Unknown command: " + args.Command),
	}, nil
}

func (p *Plugin) executeCommandLunchbotBlacklistShow(args *model.CommandArgs) *model.CommandResponse {
	message := fmt.Sprintf("Your blacklist is empty. Use '/%s' to add someone to your blacklist.", commandLunchbotBlacklistAdd)

	data := p.ReadFromStorage()
	if data.Blacklists != nil {
		if blacklist, ok := data.Blacklists[args.UserId]; ok {
			message = "Users on your blacklist:\n"
			for entry := range blacklist {
				user, _ := p.API.GetUser(entry)
				message += fmt.Sprintf("  - %s\n", user.GetDisplayName(""))
			}
		}
	}

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         message,
	}
}

func (p *Plugin) executeCommandLunchbotBlacklistAdd(args *model.CommandArgs) *model.CommandResponse {
	givenUserID := strings.TrimPrefix(args.Command, fmt.Sprintf("/%s", commandLunchbotBlacklistAdd))
	givenUserID = strings.TrimPrefix(givenUserID, " ")
	if len(givenUserID) <= 0 {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         "Error: Please enter a user you want to blacklist",
		}
	}

	user := p.GetUser(givenUserID)
	if user == nil {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf("Error: Cannot find the user %s", givenUserID),
		}
	}

	data := p.ReadFromStorage()
	if data.Blacklists == nil {
		data.Blacklists = make(map[string]map[string]struct{})
	}
	if blacklist, ok := data.Blacklists[args.UserId]; ok {
		blacklist[user.Id] = struct{}{}
	} else {
		data.Blacklists[args.UserId] = map[string]struct{}{user.Id: struct{}{}}
	}
	p.WriteToStorage(&data)

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         fmt.Sprintf("Added '%s' to your blacklist", user.GetDisplayName("")),
	}
}

func (p *Plugin) executeCommandLunchbotBlacklistRemove(args *model.CommandArgs) *model.CommandResponse {
	givenUserID := strings.TrimPrefix(args.Command, fmt.Sprintf("/%s", commandLunchbotBlacklistRemove))
	givenUserID = strings.TrimPrefix(givenUserID, " ")
	if len(givenUserID) <= 0 {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         "Error: Please enter a user you want to remove from your blacklist",
		}
	}

	user := p.GetUser(givenUserID)
	if user == nil {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf("Error: Cannot find the user '%s'", givenUserID),
		}
	}

	data := p.ReadFromStorage()
	if data.Blacklists != nil {
		if blacklist, ok := data.Blacklists[args.UserId]; ok {
			delete(blacklist, user.Id)
			p.WriteToStorage(&data)

			return &model.CommandResponse{
				ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
				Text:         fmt.Sprintf("Removed '%s' from your blacklist", user.GetDisplayName("")),
			}
		}
	}

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         fmt.Sprintf("Error: Cannot remove '%s' from your blacklist.", user.GetDisplayName("")),
	}
}

func (p *Plugin) executeCommandLunchbotTopicsShow(args *model.CommandArgs) *model.CommandResponse {
	message := fmt.Sprintf("There are no topics set yet... Use '/%s' to set a topic.", commandLunchbotTopicsShow)

	data := p.ReadFromStorage()
	if data.UserTopics != nil {
		if topics, ok := data.UserTopics[args.UserId]; ok {
			message = "Your topics:\n"
			for entry := range topics {
				message += fmt.Sprintf("  - %s\n", entry)
			}
		}
	}

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         message,
	}
}

func (p *Plugin) executeCommandLunchbotTopicsAdd(args *model.CommandArgs) *model.CommandResponse {
	givenTopic := strings.TrimPrefix(args.Command, fmt.Sprintf("/%s", commandLunchbotTopicsAdd))
	givenTopic = strings.TrimPrefix(givenTopic, " ")
	if len(givenTopic) <= 0 {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         "Error: Please enter a valid topic",
		}
	}

	data := p.ReadFromStorage()
	if data.UserTopics == nil {
		data.UserTopics = make(map[string]map[string]struct{})
	}
	if topics, ok := data.UserTopics[args.UserId]; ok {
		topics[givenTopic] = struct{}{}
	} else {
		data.UserTopics[args.UserId] = map[string]struct{}{givenTopic: struct{}{}}
	}
	p.WriteToStorage(&data)

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         fmt.Sprintf("Added '%s' to your topics", givenTopic),
	}
}

func (p *Plugin) executeCommandLunchbotTopicsRemove(args *model.CommandArgs) *model.CommandResponse {
	givenTopic := strings.TrimPrefix(args.Command, fmt.Sprintf("/%s", commandLunchbotTopicsRemove))
	givenTopic = strings.TrimPrefix(givenTopic, " ")
	if len(givenTopic) <= 0 {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         "Error: Please enter a valid topic",
		}
	}

	data := p.ReadFromStorage()
	if data.UserTopics != nil {
		if topics, ok := data.UserTopics[args.UserId]; ok {
			if _, ok := data.UserTopics[args.UserId][givenTopic]; ok {
				delete(topics, givenTopic)
				p.WriteToStorage(&data)
				return &model.CommandResponse{
					ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
					Text:         fmt.Sprintf("Removed '%s' from your topics", givenTopic),
				}
			}
		}
	}

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         fmt.Sprintf("Error: Cannot remove '%s' from your topics.", givenTopic),
	}
}

func (p *Plugin) executeCommandLunchbot(args *model.CommandArgs) *model.CommandResponse {
	triggerUser, err := p.API.GetUser(args.UserId)
	if err != nil {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         "Error: Cannot get your user...",
		}
	}
	pairedUser, err := p.GetPairingForUserID(args.ChannelId, args.UserId)
	if err != nil {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         "Error: Cannot match you with a user from this channel",
		}
	}

	data := p.ReadFromStorage()
	if data.LastPairings == nil {
		data.LastPairings = map[string][]string{}
	}
	data.LastPairings[triggerUser.Id] = append(data.LastPairings[triggerUser.Id], pairedUser.Id)
	data.LastPairings[pairedUser.Id] = append(data.LastPairings[pairedUser.Id], triggerUser.Id)
	p.WriteToStorage(&data)

	users := []string{triggerUser.Id, pairedUser.Id}
	resp := p.SendGroupMessage("Hey! I think both of you should meet for lunch soon!", users)

	if resp != nil {
		return resp
	}

	topics := p.GetRandomTopicsMsg(users)
	resp = p.SendGroupMessage(topics, users)
	if resp != nil {
		return resp
	}

	return &model.CommandResponse{}
}
