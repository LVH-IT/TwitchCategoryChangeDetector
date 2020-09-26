package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
)

//Some globally used variables
var bearerToken string        //loaded via loadConfig()
var clientID string           //loaded via loadConfig()
var soundFile string          //loaded via loadConfig()
var clientSecret string       //loaded via loadConfig()
var useCategoryWhitelist bool //loaded via loadConfig()
var categories []string       //loaded via loadConfig()
var streamInfo stream
var routineDone bool
var wasOnline bool = false
var oldGameName string
var newGameID string
var oldGameID string
var oldGameInfo game
var newGameName string
var newStreamInfo stream
var whitelistGameIDs []string

//Definitions of all the flags
var streamName string
var retryInterval int = 10

func main() {
	loadConfig()
	parseFlags()
	checkAPIToken()
	streamInfo = getStreamInfoWithOnlineCheck()
	wasOnline = true

	oldGameInfo = getGameInfo(streamInfo)
	for len(oldGameInfo.Data) == 0 {
		time.Sleep(1e9)
		oldGameInfo = getGameInfo(streamInfo)
	}
	if useCategoryWhitelist {
		whitelistGameIDs = getWhitelistGameIDs()
	}
	singleStream()
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

func getStreamInfoWithOnlineCheck() stream {
	var tempStreamInfo (stream)
	tempStreamInfo = getStreamInfo()
	if len(tempStreamInfo.Data) != 0 {
		if tempStreamInfo.Data[0].Islive == false {
			if wasOnline == true {
				time.Sleep(10e9)                   //Wait a bit and check again to make sure the stream actually went offline (The API needs some seconds to update every server's cache)
				tempStreamInfo = getStreamInfo()   //
				if len(tempStreamInfo.Data) != 0 { //
					if tempStreamInfo.Data[0].Islive == false {
						playSound()
						fmt.Println("\n" + streamName + " just went offline\n")
						for len(tempStreamInfo.Data) == 0 {
							waitRetryInterval()
							tempStreamInfo = getStreamInfo()
						}
					}
				}
			} else {
				fmt.Printf(streamName + " is currently offline. Waiting for change (Checking every " + fmt.Sprint(retryInterval) + "s)\n")
				for tempStreamInfo.Data[0].Islive == false {
					waitRetryInterval()
					tempStreamInfo = getStreamInfo()
				}
				fmt.Printf(streamName + " just went online\n")
				playSound()
			}
		}
	}
	return tempStreamInfo
}

func waitRetryInterval() {
	startTime := time.Now()
	elapsed := int(time.Since(startTime).Seconds())
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

func checkError(errorVar error) {
	if errorVar != nil {
		log.Fatal(errorVar)
		os.Exit(1)
	}
}
