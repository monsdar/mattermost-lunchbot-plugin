package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
)

// GetRandomTopic return a random topic for the given userId
// Return nothing when there's no topics stored for the userId
func (p *Plugin) GetRandomTopic(userID string) string {

	data := p.ReadFromStorage()
	if topics, ok := data.UserTopics[userID]; ok {
		topicKeys := reflect.ValueOf(topics).MapKeys()
		if len(topicKeys) > 0 {
			rand.Shuffle(len(topicKeys), func(i, j int) {
				topicKeys[i], topicKeys[j] = topicKeys[j], topicKeys[i]
			})
			randomTopic := topicKeys[0]
			return fmt.Sprintf("You could talk about %s.", randomTopic)
		}
	}

	return ""
}

//GetUser returns a user that is identified by a given string. It tries different ways to get the user.
func (p *Plugin) GetUser(userStr string) *model.User {
	//first try to get the user by username
	user, err := p.API.GetUserByUsername(userStr)
	if err == nil {
		return user
	}

	//then try to get the user by userId
	user, err = p.API.GetUser(userStr)
	if err == nil {
		return user
	}

	//then remove the `@` and check if it's possible to find the user then
	userStr = strings.TrimPrefix(userStr, "@")
	user, err = p.API.GetUserByUsername(userStr)
	if err == nil {
		return user
	}
	user, err = p.API.GetUser(userStr)
	if err == nil {
		return user
	}

	//return nil when there is no user found
	return nil
}

//SendPrivateMessage sends the given message to the given userID
func (p *Plugin) SendPrivateMessage(message string, userID string) *model.CommandResponse {
	channel, err := p.API.GetDirectChannel(p.botID, userID)
	if err != nil {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf("Error: Cannot get as direct channel to message %s", userID),
		}
	}
	post := &model.Post{
		ChannelId: channel.Id,
		UserId:    p.botID,
		Message:   message,
	}
	if _, err = p.API.CreatePost(post); err != nil {
		const errorMessage = "Error: Failed to create post"
		p.API.LogError(errorMessage, "err", err.Error())
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         errorMessage,
		}
	}
	return nil
}

// GetPairingForUserID returns a random user that is found in the given channel and that is not a bot
// This function is limited to 1000 users per channel
func (p *Plugin) GetPairingForUserID(channelID string, userID string) (*model.User, *model.AppError) {
	users, _ := p.API.GetUsersInChannel(channelID, "username", 0, 1000)
	rand.Shuffle(len(users), func(i, j int) {
		users[i], users[j] = users[j], users[i]
	})

	//we need to check the blacklist later on, read here to avoid multiple reads
	data := p.ReadFromStorage()

	targetuser := new(model.User)
	hasUserBeenFound := false
	for _, user := range users {
		if user.Id == userID {
			continue
		}
		if user.IsBot {
			continue
		}

		//check blacklists of triggering and target user as well
		if data.Blacklists != nil {
			if blacklist, ok := data.Blacklists[userID]; ok {
				if _, ok := blacklist[user.Id]; ok {
					continue
				}
			}
			if blacklist, ok := data.Blacklists[user.Id]; ok {
				if _, ok := blacklist[userID]; ok {
					continue
				}
			}
		}

		status, err := p.API.GetUserStatus(user.Id)
		if (err != nil) || (status.Status == "offline") {
			continue
		}

		targetuser = user
		hasUserBeenFound = true
		break
	}

	if !hasUserBeenFound {
		return nil, &model.AppError{
			Message: "Cannot match a user in this channel...",
		}
	}
	return targetuser, nil
}
