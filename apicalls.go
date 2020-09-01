package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

func getGameInfo(gameID string) game {
	gameJSON := getAPIInfo("games?id=" + gameID)
	var gameInfo game
	json.Unmarshal([]byte(gameJSON), &gameInfo)
	return gameInfo
}

type gameData struct {
	BoxArtURL string `json:"box_art_url"`
	ID        int    `json:"id"`
	Name      string `json:"name"`
}

type game struct {
	Data gameData `json:"data"`
}

func getStreamInfo(streamerName string) stream {
	streamJSON := getAPIInfo("streams?user_login=" + streamerName)
	var streamInfo stream
	json.Unmarshal([]byte(streamJSON), &streamInfo)
	return streamInfo
}

type streamData struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	UserName     string `json:"user_name"`
	GameID       string `json:"game_id"`
	Type         string `json:"type"`
	Title        string `json:"title"`
	ViewerCount  string `json:"viewer_count"`
	StartedAt    string `json:"started_at"`
	Language     string `json:"language"`
	ThumbnailURL string `json:"thumbnail_url"`
	TagIDs       string `json:"tag_ids"`
}

type stream struct {
	Data streamData `json:"data"`
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
	apiJSON = strings.ReplaceAll(apiJSON, "[", "")
	apiJSON = strings.ReplaceAll(apiJSON, "]", "")
	return apiJSON
}
