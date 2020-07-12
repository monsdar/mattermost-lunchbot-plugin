package main

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"testing"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetPairingForUserID(t *testing.T) {
	t.Run("Empty channel, empty data", func(t *testing.T) {
		lunchbotData := &LunchbotData{}
		reqBodyBytes := new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(lunchbotData)

		users := []*model.User{
			&model.User{
				Id: "1337",
			},
		}

		plugin := &Plugin{}
		api := &plugintest.API{}
		api.On("GetUsersInChannel", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(users, nil)
		api.On("KVGet", mock.AnythingOfType("string")).Return(reqBodyBytes.Bytes(), nil)
		api.On("GetUserStatus", mock.AnythingOfType("string")).Return(&model.Status{Status: "online"}, nil)
		plugin.SetAPI(api)

		_, err := plugin.GetPairingForUserID("", "1337")
		assert.Equal(t, "Cannot find a user to pair with in this channel...", err.Message)
	})

	t.Run("Some users, empty data", func(t *testing.T) {
		lunchbotData := &LunchbotData{}
		reqBodyBytes := new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(lunchbotData)

		users := []*model.User{
			&model.User{
				Id: "1",
			},
			&model.User{
				Id: "1337",
			},
		}

		plugin := &Plugin{}
		api := &plugintest.API{}
		api.On("GetUsersInChannel", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(users, nil)
		api.On("KVGet", mock.AnythingOfType("string")).Return(reqBodyBytes.Bytes(), nil)
		api.On("GetUserStatus", mock.AnythingOfType("string")).Return(&model.Status{Status: "online"}, nil)
		plugin.SetAPI(api)

		user, _ := plugin.GetPairingForUserID("", "1337")
		assert.Equal(t, "1", user.Id)
	})

	t.Run("Bot user", func(t *testing.T) {
		lunchbotData := &LunchbotData{}
		reqBodyBytes := new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(lunchbotData)

		users := []*model.User{
			&model.User{
				Id:    "1",
				IsBot: true,
			},
			&model.User{
				Id: "1337",
			},
		}

		plugin := &Plugin{}
		api := &plugintest.API{}
		api.On("GetUsersInChannel", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(users, nil)
		api.On("KVGet", mock.AnythingOfType("string")).Return(reqBodyBytes.Bytes(), nil)
		api.On("GetUserStatus", mock.AnythingOfType("string")).Return(&model.Status{Status: "online"}, nil)
		plugin.SetAPI(api)

		_, err := plugin.GetPairingForUserID("", "1337")
		assert.Equal(t, "Cannot find a user to pair with in this channel...", err.Message)
	})

	t.Run("Offline user", func(t *testing.T) {
		lunchbotData := &LunchbotData{}
		reqBodyBytes := new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(lunchbotData)

		users := []*model.User{
			&model.User{
				Id: "1",
			},
			&model.User{
				Id: "1337",
			},
		}

		plugin := &Plugin{}
		api := &plugintest.API{}
		api.On("GetUsersInChannel", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(users, nil)
		api.On("KVGet", mock.AnythingOfType("string")).Return(reqBodyBytes.Bytes(), nil)
		api.On("GetUserStatus", mock.AnythingOfType("string")).Return(&model.Status{Status: "offline"}, nil)
		plugin.SetAPI(api)

		_, err := plugin.GetPairingForUserID("", "1337")
		assert.Equal(t, "Cannot find a user to pair with in this channel...", err.Message)
	})

	t.Run("Blacklisted user", func(t *testing.T) {
		lunchbotData := &LunchbotData{
			Blacklists: map[string]map[string]struct{}{
				"1337": map[string]struct{}{
					"1": struct{}{},
				},
			},
		}
		reqBodyBytes := new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(lunchbotData)

		users := []*model.User{
			&model.User{
				Id: "1",
			},
			&model.User{
				Id: "1337",
			},
		}

		plugin := &Plugin{}
		api := &plugintest.API{}
		api.On("GetUsersInChannel", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(users, nil)
		api.On("KVGet", mock.AnythingOfType("string")).Return(reqBodyBytes.Bytes(), nil)
		api.On("GetUserStatus", mock.AnythingOfType("string")).Return(&model.Status{Status: "online"}, nil)
		plugin.SetAPI(api)

		_, err := plugin.GetPairingForUserID("", "1337")
		assert.Equal(t, "Cannot find a user to pair with in this channel...", err.Message)
	})

	t.Run("Triggering user is blacklisted", func(t *testing.T) {
		lunchbotData := &LunchbotData{
			Blacklists: map[string]map[string]struct{}{
				"1": map[string]struct{}{
					"1337": struct{}{},
				},
			},
		}
		reqBodyBytes := new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(lunchbotData)

		users := []*model.User{
			&model.User{
				Id: "1",
			},
			&model.User{
				Id: "1337",
			},
		}

		plugin := &Plugin{}
		api := &plugintest.API{}
		api.On("GetUsersInChannel", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(users, nil)
		api.On("KVGet", mock.AnythingOfType("string")).Return(reqBodyBytes.Bytes(), nil)
		api.On("GetUserStatus", mock.AnythingOfType("string")).Return(&model.Status{Status: "online"}, nil)
		plugin.SetAPI(api)

		_, err := plugin.GetPairingForUserID("", "1337")
		assert.Equal(t, "Cannot find a user to pair with in this channel...", err.Message)
	})
}
func TestGetPairingForUserID_weightedTests(t *testing.T) {
	t.Run("Many users, checking weighted choice", func(t *testing.T) {
		lunchbotData := &LunchbotData{
			LastPairings: map[string][]string{
				"1337": []string{
					"SomeTimeAgo",
					"MediumRecent",
					"MostRecent",
				},
			},
		}
		reqBodyBytes := new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(lunchbotData)

		users := []*model.User{
			&model.User{
				Id: "SomeTimeAgo",
			},
			&model.User{
				Id: "MediumRecent",
			},
			&model.User{
				Id: "MostRecent",
			},
			&model.User{
				Id: "NeverPairedWith",
			},
			&model.User{
				Id: "1337",
			},
		}

		plugin := &Plugin{}
		api := &plugintest.API{}
		api.On("GetUsersInChannel", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(users, nil)
		api.On("KVGet", mock.AnythingOfType("string")).Return(reqBodyBytes.Bytes(), nil)
		api.On("GetUserStatus", mock.AnythingOfType("string")).Return(&model.Status{Status: "online"}, nil)
		plugin.SetAPI(api)

		rand.Seed(1337) //initing rand for this test to be deterministic
		//the user that has never been paired has by far the highest weight, so it will be chosen in almost every case
		user, _ := plugin.GetPairingForUserID("", "1337")
		assert.Equal(t, "NeverPairedWith", user.Id)
		user, _ = plugin.GetPairingForUserID("", "1337")
		assert.Equal(t, "NeverPairedWith", user.Id)
		user, _ = plugin.GetPairingForUserID("", "1337")
		assert.Equal(t, "NeverPairedWith", user.Id)
		user, _ = plugin.GetPairingForUserID("", "1337")
		assert.Equal(t, "NeverPairedWith", user.Id)
		user, _ = plugin.GetPairingForUserID("", "1337")
		assert.Equal(t, "NeverPairedWith", user.Id)
		user, _ = plugin.GetPairingForUserID("", "1337")
		assert.Equal(t, "NeverPairedWith", user.Id)
		user, _ = plugin.GetPairingForUserID("", "1337")
		assert.Equal(t, "NeverPairedWith", user.Id)
		user, _ = plugin.GetPairingForUserID("", "1337")
		assert.Equal(t, "NeverPairedWith", user.Id)
		user, _ = plugin.GetPairingForUserID("", "1337")
		assert.Equal(t, "NeverPairedWith", user.Id)
		user, _ = plugin.GetPairingForUserID("", "1337")
		assert.Equal(t, "NeverPairedWith", user.Id)

		//add "NeverPairedWith" as most recent pairing and see what happens
		lunchbotData = &LunchbotData{
			LastPairings: map[string][]string{
				"1337": []string{
					"SomeTimeAgo",
					"MediumRecent",
					"MostRecent",
					"NeverPairedWith",
				},
			},
		}
		reqBodyBytes = new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(lunchbotData)
		api = &plugintest.API{}
		api.On("GetUsersInChannel", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(users, nil)
		api.On("KVGet", mock.AnythingOfType("string")).Return(reqBodyBytes.Bytes(), nil)
		api.On("GetUserStatus", mock.AnythingOfType("string")).Return(&model.Status{Status: "online"}, nil)
		plugin.SetAPI(api)

		user, _ = plugin.GetPairingForUserID("", "1337")
		assert.Equal(t, "NeverPairedWith", user.Id)
		user, _ = plugin.GetPairingForUserID("", "1337")
		assert.Equal(t, "SomeTimeAgo", user.Id)
		user, _ = plugin.GetPairingForUserID("", "1337")
		assert.Equal(t, "MediumRecent", user.Id)
		user, _ = plugin.GetPairingForUserID("", "1337")
		assert.Equal(t, "SomeTimeAgo", user.Id)
		user, _ = plugin.GetPairingForUserID("", "1337")
		assert.Equal(t, "SomeTimeAgo", user.Id)
		user, _ = plugin.GetPairingForUserID("", "1337")
		assert.Equal(t, "NeverPairedWith", user.Id)
		user, _ = plugin.GetPairingForUserID("", "1337")
		assert.Equal(t, "MediumRecent", user.Id)
		user, _ = plugin.GetPairingForUserID("", "1337")
		assert.Equal(t, "MediumRecent", user.Id)
		user, _ = plugin.GetPairingForUserID("", "1337")
		assert.Equal(t, "SomeTimeAgo", user.Id)
		user, _ = plugin.GetPairingForUserID("", "1337")
		assert.Equal(t, "MostRecent", user.Id)
	})
}
