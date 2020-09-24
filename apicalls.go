package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func getGameInfo(gameID string) game {
	gameJSON := getAPIInfo("games?id=" + gameID)
	var gameInfo game
	json.Unmarshal([]byte(gameJSON), &gameInfo)
	return gameInfo
}

type game struct {
	Data []struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		BoxArtURL string `json:"box_art_url"`
	} `json:"data"`
}

func getStreamInfo(streamerName string) stream {
	streamJSON := getAPIInfo("streams?user_login=" + streamerName)
	var streamInfo stream
	json.Unmarshal([]byte(streamJSON), &streamInfo)
	return streamInfo
}

type stream struct {
	Data []struct {
		ID           string   `json:"id"`
		UserID       string   `json:"user_id"`
		UserName     string   `json:"user_name"`
		GameID       string   `json:"game_id"`
		Type         string   `json:"type"`
		Title        string   `json:"title"`
		ViewerCount  string   `json:"viewer_count"`
		StartedAt    string   `json:"started_at"`
		Language     string   `json:"language"`
		ThumbnailURL string   `json:"thumbnail_url"`
		TagIDs       []string `json:"tag_ids"`
	} `json:"data"`
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

type apiValidation struct {
	Error     string   `json:"error"`
	Status    int      `json:"status"`
	Message   string   `json:"message"`
	ClientID  string   `json:"client_id"`
	Scopes    []string `json:"scopes"`
	ExpiresIn int      `json:"expires_in"`
}

type apiToken struct {
	AccessToken string `json:"access_token"`
	TokenType   int    `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func checkAPIToken() {
	httpClient := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.twitch.tv/helix/streams?user_login=twitch", nil)
	req.Header.Add("Authorization", "Bearer "+bearerToken)
	req.Header.Add("Client-ID", clientID)
	resp, err := httpClient.Do(req)
	checkError(err)
	body, _ := ioutil.ReadAll(resp.Body)
	apiJSON := string(body)
	var APIInfo apiValidation
	json.Unmarshal([]byte(apiJSON), &APIInfo)

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
	apiJSON := string(body)
	var TokenResponse apiToken
	json.Unmarshal([]byte(apiJSON), &TokenResponse)

	//Check Responsse for validity
	if len(TokenResponse.AccessToken) != 30 {
		fmt.Println("Could not obtain a new Bearer Access Token. Please try again or get one manually")
		os.Exit(1)
	}
	return TokenResponse
}
