package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func getWhitelistGameIDs() []string {
	var params (string) = "games?"
	var gameJSON string
	var gameIDs []string
	for a, b := range categories {
		if a == 0 {
			params += "name=" + url.QueryEscape(b)
		} else {
			params += "&name=" + url.QueryEscape(b)
		}
	}
	gameJSON = getAPIInfo(params)
	var gameInfo game
	json.Unmarshal([]byte(gameJSON), &gameInfo)
	for a := range gameInfo.Data {
		gameIDs = append(gameIDs, gameInfo.Data[a].ID)
	}
	return gameIDs
}

func getGameInfo(streamDataJSON stream) game {
	var params (string) = "games?"
	var gameJSON string
	params += "id=" + streamDataJSON.Data[0].GameID
	gameJSON = getAPIInfo(params)
	var gameInfo game
	json.Unmarshal([]byte(gameJSON), &gameInfo)
	return gameInfo
}

func getStreamInfo() stream {
	var params (string) = "search/channels?first=1"
	params += "&query=" + streamName
	streamJSON := getAPIInfo(params)
	json.Unmarshal([]byte(streamJSON), &streamInfo)
	return streamInfo
}

func getAPIInfo(endPointAndParams string) string {
	httpClient := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.twitch.tv/helix/"+endPointAndParams, nil)
	req.Header.Add("Authorization", "Bearer "+bearerToken)
	req.Header.Add("Client-ID", clientID)
	resp, err := httpClient.Do(req)
	checkError(err)
	body, _ := ioutil.ReadAll(resp.Body)
	apiJSON := string(body)
	return apiJSON
}

func checkAPIToken() {
	httpClient := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.twitch.tv/helix/streams?user_login=twitch", nil)
	req.Header.Add("Authorization", "Bearer "+bearerToken)
	req.Header.Add("Client-ID", clientID)
	resp, err := httpClient.Do(req)
	checkError(err)
	body, _ := ioutil.ReadAll(resp.Body)
	apiJSON := body
	var APIInfo apiValidation
	json.Unmarshal(apiJSON, &APIInfo)

	if APIInfo.Error != "" {
		println("Error " + fmt.Sprint(APIInfo.Status) + " (" + APIInfo.Error + "): " + APIInfo.Message)

		var getNewOne string
		println("Do you want to get a new one?")
		println("Type \"y\" for yes or \"n\" for no")
		fmt.Scanf("%s", &getNewOne)

		if getNewOne == "y" {
			//Get new Token
			TokenResponse := getNewAPITokenJSON()

			//Load Config File
			configFile, err := ioutil.ReadFile("config.json")
			checkError(err)
			var configData config
			json.Unmarshal([]byte(string(configFile)), &configData)

			//Change Token and write it to config.json
			configData.BearerToken = TokenResponse.AccessToken
			configBytes, err := json.MarshalIndent(configData, "", "    ")
			checkError(err)
			err = ioutil.WriteFile("config.json", configBytes, 0644)
			checkError(err)

			//Load new Config and continue
			loadConfig()
			println("New Token has been obtained")
			println("Starting to monitor")
			println("")
		} else if getNewOne == "n" {
			os.Exit(1)
		} else {
			println("Invalid input. You can only type \"y\" or \"n\". Exiting")
			os.Exit(1)
		}
	}
}

func getNewAPITokenJSON() apiToken {
	httpClient := &http.Client{}
	req, _ := http.NewRequest("POST", "https://id.twitch.tv/oauth2/token?client_id="+clientID+"&client_secret="+clientSecret+"&grant_type=client_credentials", nil)
	resp, err := httpClient.Do(req)
	checkError(err)
	body, _ := ioutil.ReadAll(resp.Body)
	apiJSON := body
	var TokenResponse apiToken
	json.Unmarshal(apiJSON, &TokenResponse)

	//Check Responsse for validity
	if len(TokenResponse.AccessToken) != 30 {
		println("Could not obtain a new Bearer Access Token. Please try again or get one manually")
		os.Exit(1)
	}
	return TokenResponse
}
