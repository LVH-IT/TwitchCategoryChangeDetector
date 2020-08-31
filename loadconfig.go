package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func loadConfig() {
	configFile, err := ioutil.ReadFile("config.json")
	checkError(err)
	var configData config
	json.Unmarshal([]byte(string(configFile)), &configData)

	bearerToken = configData.BearerToken
	clientID = configData.ClientID
	soundFile = configData.SoundFile

	if len(bearerToken) != 30 {
		fmt.Printf("The provided Bearer Token is invalid")
		os.Exit(1)
	}
	if len(clientID) != 30 {
		fmt.Printf("The provided Client ID is invalid")
		os.Exit(1)
	}
	audioFile, err := os.Open(soundFile)
	checkError(err)
	audioFile.Close()
}

type config struct {
	BearerToken string `json:"BearerToken"`
	ClientID    string `json:"ClientID"`
	SoundFile   string `json:"SoundFile"`
}
