package main

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mroth/weightedrand"
)

// GetRandomTopicsMsg return a random topic for the given userIDs
// Return nothing when there's no topics stored for the userIDs
func (p *Plugin) GetRandomTopicsMsg(userIDs []string) string {

	randomTopics := []string{}
	for _, userID := range userIDs {
		data := p.ReadFromStorage()
		if topics, ok := data.UserTopics[userID]; ok {
			topicKeys := reflect.ValueOf(topics).MapKeys()
			if len(topicKeys) > 0 {
				rand.Shuffle(len(topicKeys), func(i, j int) {
					topicKeys[i], topicKeys[j] = topicKeys[j], topicKeys[i]
				})
				randomTopics = append(randomTopics, fmt.Sprintf("%s", topicKeys[0]))
			}
		}
	}

	if len(randomTopics) <= 0 {
		return ""
	}

	return fmt.Sprintf("You could talk about %s", strings.Join(randomTopics, " or "))
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

//SendGroupMessage sends the given message to the given userIDs
func (p *Plugin) SendGroupMessage(message string, userIDs []string) *model.CommandResponse {
	userIDs = append(userIDs, p.botID)
	channel, err := p.API.GetGroupChannel(userIDs)
	if err != nil {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf("Error: Cannot get as group channel to message %s", userIDs),
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

	//read the users data for blacklist and weightedrandom
	data := p.ReadFromStorage()

	weightedUsers := []weightedrand.Choice{} //list of users, sorted by weight
	for _, user := range users {
		//is this the triggering user?
		if user.Id == userID {
			continue
		}
		//is this a bot?
		if user.IsBot {
			continue
		}
		//is the user already paired?
		if _, ok := data.ActivePairings[user.Id]; ok {
			continue
		}
		//is this user offline?
		status, err := p.API.GetUserStatus(user.Id)
		if (err != nil) || (status.Status == "offline") {
			continue
		}
		//is this user on a blacklist? Is the triggering user on the users blacklist?
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

		//check if the user has already been paired lately. Add him with a weight according to how recent the pairing has been
		//by iterating in reverse we make sure that users that appear multiple times in the list will not mess up the weights
		isNewUser := true
		if data.LastPairings != nil {
			for index := len(data.LastPairings[userID]) - 1; index >= 0; index-- {
				currentUserID := data.LastPairings[userID][index]
				if currentUserID == user.Id {
					userWeight := uint(math.Abs(float64(index - len(data.LastPairings[userID]))))
					weightedUsers = append(weightedUsers, weightedrand.Choice{Weight: userWeight, Item: user})
					isNewUser = false
					break
				}
			}
			if !isNewUser { //if the user has been found within our recentpairing we can continue the loop
				continue
			}
		}

		//Finally... this is a brand-new user that has never paired with our triggering user. Add him with a very high weight, so he'll be chosen with a high possibility
		weightedUsers = append(weightedUsers, weightedrand.Choice{Weight: 1000, Item: user})
	}

	if len(weightedUsers) > 0 {
		chooser := weightedrand.NewChooser(weightedUsers...)
		user, ok := chooser.Pick().(*model.User)
		if ok {
			return user, nil
		}
	}

	return nil, &model.AppError{
		Message: "Cannot find a user to pair with in this channel...",
	}
}
