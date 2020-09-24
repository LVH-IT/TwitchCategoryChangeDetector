package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/m7shapan/njson"
)

func getGameInfo(gameID string) game {
	gameJSON := getAPIInfo("games?id=" + gameID)
	var gameInfo game
	njson.Unmarshal([]byte(gameJSON), &gameInfo)
	return gameInfo
}

type game struct {
	ID        int    `njson:"data.0.id"`
	Name      string `njson:"data.0.name"`
	BoxArtURL string `njson:"data.0.box_art_url"`
}

func getStreamInfo(streamerName string) stream {
	streamJSON := getAPIInfo("streams?user_login=" + streamerName)
	var streamInfo stream
	njson.Unmarshal([]byte(streamJSON), &streamInfo)
	return streamInfo
}

type stream struct {
	ID           string   `njson:"data.0.id"`
	UserID       string   `njson:"data.0.user_id"`
	UserName     string   `njson:"data.0.user_name"`
	GameID       string   `njson:"data.0.game_id"`
	Type         string   `njson:"data.0.type"`
	Title        string   `njson:"data.0.title"`
	ViewerCount  string   `njson:"data.0.viewer_count"`
	StartedAt    string   `njson:"data.0.started_at"`
	Language     string   `njson:"data.0.language"`
	ThumbnailURL string   `njson:"data.0.thumbnail_url"`
	TagIDs       []string `njson:"data.0.tag_ids"`
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
	Error     string   `njson:"error"`
	Status    int      `njson:"status"`
	Message   string   `njson:"message"`
	ClientID  string   `njson:"client_id"`
	Scopes    []string `njson:"scopes"`
	ExpiresIn int      `njson:"expires_in"`
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
	njson.Unmarshal([]byte(apiJSON), &APIInfo)

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
