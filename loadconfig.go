package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type config struct {
	BearerToken                string   `json:"BearerToken"`
	ClientID                   string   `json:"ClientID"`
	ClientSecret               string   `json:"ClientSecret"`
	SoundFile                  string   `json:"SoundFile"`
	UseCategoryWhitelist       bool     `json:"UseCategoryWhitelist"`
	Categories                 []string `json:"Categories"`
	NotifyOnOfflineTitleChange bool     `json:"NotifyOnOfflineTitleChange"`
	NotifyOnOnlineTitleChange  bool     `json:"NotifyOnOnlineTitleChange"`
}

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

	if len(clientID) != 30 {
		println("The provided Client ID is invalid")
		os.Exit(1)
	}
	if len(clientSecret) != 30 {
		println("The provided Client Secret is invalid")
		os.Exit(1)
	}

	checkAPIToken()

	audioFile, err := os.Open(soundFile)
	checkError(err)
	audioFile.Close()
}
