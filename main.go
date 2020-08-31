package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
)

//Some globally used variables
var bearerToken string = "x3oyoq4z6291bw6iuqhnnwyzusnasq"
var clientID string = "5qzrrwjwsuum85molwsalbv1q5v1f5"
var streamInfo stream

//Definitions of all the flags
var streamName string
var retryInterval int

func main() {
	parseFlags()

	startTime := time.Now()
	elapsed := int(time.Now().Sub(startTime))
	_ = elapsed //Bypass unused error

	getStreamInfoWithOnlineCheck()

	oldGameID := streamInfo.Data.GameID
	oldGameName := getGameInfo(oldGameID).Data.Name
	newGameID := streamInfo.Data.GameID

	if oldGameName == "" {
		fmt.Printf("Game name not found. Retrying in: " + fmt.Sprint(retryInterval) + "s\n")
		for oldGameName == "" {
			waitRetryInterval()
			oldGameName = getGameInfo(oldGameID).Data.Name
			startTime = time.Now()
			elapsed = int(time.Now().Sub(startTime))
		}
		fmt.Printf("Game name found\n")
	}

	fmt.Printf("Channel to monitor: " + streamInfo.Data.UserName + "\n")
	fmt.Printf("Current Category: " + oldGameName + "\nWaiting for change (Checking every " + fmt.Sprint(retryInterval) + "s)\n")
	for oldGameID == newGameID {
		waitRetryInterval()
		getStreamInfoWithOnlineCheck()
		newGameID = streamInfo.Data.GameID
	}

	newGameName := getGameInfo(newGameID).Data.Name
	if newGameName == "" {
		fmt.Printf("Game name not found. Retrying in: " + fmt.Sprint(retryInterval) + "s\n")
		for newGameName == "" {
			waitRetryInterval()
			newGameName = getGameInfo(oldGameID).Data.Name
			startTime = time.Now()
			elapsed = int(time.Now().Sub(startTime))
		}
		fmt.Printf("Game name found\n")
	}

	fmt.Printf("Category changed to: " + newGameName)

	audioFile, err := os.Open("juntos.ogg")
	if err != nil {
		fmt.Printf("Error opening audio file juntos.ogg\n")
		os.Exit(1)
	}
	audioStreamer, audioFormat, err := vorbis.Decode(audioFile)
	if err != nil {
		fmt.Printf("Error decoding audio file juntos.ogg\n")
		os.Exit(1)
	}
	defer audioStreamer.Close()
	speaker.Init(audioFormat.SampleRate, audioFormat.SampleRate.N(time.Second/10))
	done := make(chan bool)
	speaker.Play(beep.Seq(audioStreamer, beep.Callback(func() {
		done <- true
	})))
	<-done
}

func getStreamInfoWithOnlineCheck() {
	streamInfo = getStreamInfo(streamName)

	if streamInfo.Data.StartedAt == "" {
		fmt.Printf("Channel offline. Waiting for change (Checking every " + fmt.Sprint(retryInterval) + "s)\n")
		for streamInfo.Data.StartedAt == "" {
			waitRetryInterval()
			streamInfo = getStreamInfo(streamName)
		}
		fmt.Printf("Channel went online\n")
	}
}

func waitRetryInterval() {
	startTime := time.Now()
	elapsed := int(time.Now().Sub(startTime))
	for elapsed < retryInterval {
		time.Sleep(1000000000) //sleep for 1000000000ns = 1000ms = 1s
		elapsed = int(time.Now().Sub(startTime).Seconds())
	}
}

func parseFlags() {
	flag.StringVar(&streamName, "c", "xqcow", "provide the name of the twitch channel")
	flag.IntVar(&retryInterval, "t", 10, "provide the interval (in seconds) in which to refresh the stream's information")
	flag.Parse()
}

func getGameInfo(gameID string) game {
	httpClient := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.twitch.tv/helix/games?id="+gameID, nil)
	req.Header.Add("Authorization", "Bearer "+bearerToken)
	req.Header.Add("Client-ID", clientID)
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("Error retrieving game info\n")
		os.Exit(1)
	}
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
	if err != nil {
		fmt.Printf("Error retrieving stream info\n")
		os.Exit(1)
	}
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
