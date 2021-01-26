package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func loadConfig() {
	configFile, err := ioutil.ReadFile("config.json")
	checkError(err)
	var configData config
	json.Unmarshal([]byte(string(configFile)), &configData)

	//Making the json file pretty
	configBytes, err := json.MarshalIndent(configData, "", "    ")
	checkError(err)
	err = ioutil.WriteFile("config.json", configBytes, 0644)
	checkError(err)

	bearerToken = configData.BearerToken
	clientID = configData.ClientID
	clientSecret = configData.ClientSecret
	soundFile = configData.SoundFile
	useCategoryWhitelist = configData.UseCategoryWhitelist
	categories = configData.Categories
	notifyOnOfflineTitleChange = configData.NotifyOnOfflineTitleChange
	notifyOnOnlineTitleChange = configData.NotifyOnOnlineTitleChange
	discordBotToken = configData.DiscordBotToken

	if len(clientID) != 30 {
		println("The provided Client ID is invalid")
		os.Exit(1)
	}
	if len(clientSecret) != 30 {
		println("The provided Client Secret is invalid")
		os.Exit(1)
	}

	checkAPIToken()

	if !silentMode {
		audioFile, err := os.Open(soundFile)
		checkError(err)
		audioFile.Close()
	}
}
