package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
)

//Some globally used variables
var bearerToken string //loaded via loadConfig()
var clientID string    //loaded via loadConfig()
var soundFile string   //loaded via loadConfig()
var streamInfo stream

//Definitions of all the flags
var streamName string
var retryInterval int

func main() {
	parseFlags()
	loadConfig()
	getStreamInfoWithOnlineCheck()

	oldGameID := streamInfo.GameID
	oldGameName := getGameInfo(oldGameID).Name
	newGameID := streamInfo.GameID

	if oldGameName == "" {
		fmt.Printf("Game name not found. Retrying in: " + fmt.Sprint(retryInterval) + "s\n")
		for oldGameName == "" {
			waitRetryInterval()
			oldGameName = getGameInfo(oldGameID).Name
		}
		fmt.Printf("Game name found\n")
	}

	fmt.Printf("Channel to monitor: " + streamInfo.UserName + "\n")
	fmt.Printf("Current Category: " + oldGameName + "\nWaiting for change (Checking every " + fmt.Sprint(retryInterval) + "s)\n")
	for oldGameID == newGameID {
		waitRetryInterval()
		getStreamInfoWithOnlineCheck()
		newGameID = streamInfo.GameID
	}

	newGameName := getGameInfo(newGameID).Name
	if newGameName == "" {
		fmt.Printf("Game name not found. Retrying in: " + fmt.Sprint(retryInterval) + "s\n")
		for newGameName == "" {
			waitRetryInterval()
			newGameName = getGameInfo(oldGameID).Name
		}
		fmt.Printf("Game name found\n")
	}

	fmt.Printf("Category changed to: " + newGameName + "\n")
	playSound()
}

func playSound() {
	audioFile, err := os.Open(soundFile)
	checkError(err)
	audioStreamer, audioFormat, err := vorbis.Decode(audioFile)
	checkError(err)
	defer audioStreamer.Close()
	speaker.Init(audioFormat.SampleRate, audioFormat.SampleRate.N(time.Second/10))
	done := make(chan bool)
	speaker.Play(beep.Seq(audioStreamer, beep.Callback(func() {
		done <- true
	})))
	<-done
	audioFile.Close()
}

func getStreamInfoWithOnlineCheck() {
	streamInfo = getStreamInfo(streamName)
	if streamInfo.StartedAt == "" {
		fmt.Printf(streamName + " is currently offline. Waiting for change (Checking every " + fmt.Sprint(retryInterval) + "s)\n")
		for streamInfo.StartedAt == "" {
			waitRetryInterval()
			streamInfo = getStreamInfo(streamName)
		}
		fmt.Printf(streamName + " just went online\n")
		playSound()
	}
}

func waitRetryInterval() {
	startTime := time.Now()
	elapsed := int(time.Since(startTime))
	for elapsed < retryInterval {
		time.Sleep(1000000000) //sleep for 1000000000ns = 1000ms = 1s
		elapsed = int(time.Since(startTime).Seconds())
	}
}

func parseFlags() {
	flag.StringVar(&streamName, "c", "xqcow", "provide the name of the twitch channel")
	flag.IntVar(&retryInterval, "t", 10, "provide the interval (in seconds) in which to refresh the stream's information")
	flag.Parse()
}
