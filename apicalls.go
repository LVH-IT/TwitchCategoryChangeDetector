package main

import (
	"io/ioutil"
	"net/http"

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
