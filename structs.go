package main

import (
	"github.com/bwmarrin/discordgo"
)

type game struct {
	Data []struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		BoxArtURL string `json:"box_art_url"`
	} `json:"data"`
}

type stream struct {
	Data []struct {
		BroadcasterLanguage string   `json:"broadcaster_language"`
		DisplayName         string   `json:"display_name"`
		GameID              string   `json:"game_id"`
		ID                  string   `json:"id"`
		Islive              bool     `json:"is_live"`
		TagIDs              []string `json:"tag_ids"`
		ThumbnailURL        string   `json:"thumbnail_url"`
		Title               string   `json:"title"`
		StartedAt           string   `json:"started_at"`
	} `json:"data"`
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

type config struct {
	BearerToken                string   `json:"BearerToken"`
	ClientID                   string   `json:"ClientID"`
	ClientSecret               string   `json:"ClientSecret"`
	SoundFile                  string   `json:"SoundFile"`
	UseCategoryWhitelist       bool     `json:"UseCategoryWhitelist"`
	Categories                 []string `json:"Categories"`
	NotifyOnOfflineTitleChange bool     `json:"NotifyOnOfflineTitleChange"`
	NotifyOnOnlineTitleChange  bool     `json:"NotifyOnOnlineTitleChange"`
	DiscordBotToken            string   `json:"DiscordBotToken"`
}

type dcChannels struct {
	ChannelID  string
	StreamName string
	Session    *discordgo.Session
}
