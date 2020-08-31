package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

func getGameInfo(gameID string) game {
	httpClient := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.twitch.tv/helix/games?id="+gameID, nil)
	req.Header.Add("Authorization", "Bearer "+bearerToken)
	req.Header.Add("Client-ID", clientID)
	resp, err := httpClient.Do(req)
	checkError(err)
	body, _ := ioutil.ReadAll(resp.Body)
	gameJSON := string(body)
	gameJSON = strings.ReplaceAll(gameJSON, "[", "")
	gameJSON = strings.ReplaceAll(gameJSON, "]", "")
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
	Data gameData
}

func getStreamInfo(streamerName string) stream {
	httpClient := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.twitch.tv/helix/streams?user_login="+streamerName, nil)
	req.Header.Add("Authorization", "Bearer "+bearerToken)
	req.Header.Add("Client-ID", clientID)
	resp, err := httpClient.Do(req)
	checkError(err)
	body, _ := ioutil.ReadAll(resp.Body)
	streamJSON := string(body)
	streamJSON = strings.ReplaceAll(streamJSON, "[", "")
	streamJSON = strings.ReplaceAll(streamJSON, "]", "")
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
	Data streamData
}
