package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
)

//Some globally used variables
var bearerToken string              //loaded via loadConfig()
var clientID string                 //loaded via loadConfig()
var soundFile string                //loaded via loadConfig()
var clientSecret string             //loaded via loadConfig()
var useCategoryWhitelist bool       //loaded via loadConfig()
var categories []string             //loaded via loadConfig()
var notifyOnOfflineTitleChange bool //loaded via loadConfig()
var notifyOnOnlineTitleChange bool  //loaded via loadConfig()

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
var oldTitle string
var newTitle string
var initOnline = 1
var initOffline = 1
var hostOS = runtime.GOOS

//Definitions of all the flags
var streamName string
var retryInterval int = 10

func main() {
	checkOS()
	loadConfig()
	parseFlags()
	checkAPIToken()
	for {
		streamInfo = getStreamInfoWithOnlineCheck()
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
}

func checkOS() {
	if hostOS != "windows" && hostOS != "linux" {
		println("Your OS is not supported. Exiting in 10s")
		time.Sleep(10e9)
		os.Exit(1)
	}
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
	var tempStreamInfo stream
	var offlineOldTitle string
	var offlineNewTitle string
	tempStreamInfo = getStreamInfo()
	offlineOldTitle = tempStreamInfo.Data[0].Title
	if tempStreamInfo.Data[0].Islive {
		wasOnline = true
	} else {
		if wasOnline {
			playSound()
			clearCLI()
			println(streamName + " just went offline. Waiting for change (Checking every " + fmt.Sprint(retryInterval) + "s)")
			for !tempStreamInfo.Data[0].Islive {
				waitRetryInterval()
				if initOffline == 6 {
					tempStreamInfo = getStreamInfo()
					offlineNewTitle = tempStreamInfo.Data[0].Title
					if offlineOldTitle != offlineNewTitle {
						if notifyOnOfflineTitleChange {
							playSound()
						}
						clearCLI()
						println("Title changed to: " + offlineNewTitle)
						println("----------------------------------------------------")
						println(streamName + " just went offline. Waiting for change (Checking every " + fmt.Sprint(retryInterval) + "s)")
						offlineOldTitle = offlineNewTitle
					}
				}
				if initOffline < 6 {
					initOffline++
				}
			}
		} else {
			println(streamName + " is currently offline. Waiting for change (Checking every " + fmt.Sprint(retryInterval) + "s)")
			for !tempStreamInfo.Data[0].Islive {
				waitRetryInterval()
				if initOffline == 6 {
					tempStreamInfo = getStreamInfo()
					offlineNewTitle = tempStreamInfo.Data[0].Title
					if offlineOldTitle != offlineNewTitle {
						if notifyOnOfflineTitleChange {
							playSound()
						}
						clearCLI()
						println("Title changed to: " + offlineNewTitle)
						println("----------------------------------------------------")
						println(streamName + " is currently offline. Waiting for change (Checking every " + fmt.Sprint(retryInterval) + "s)")
						offlineOldTitle = offlineNewTitle
					}
				}
				if initOffline < 6 {
					initOffline++
				}
			}
			wasOnline = true
		}
		initOffline = 1
		initOnline = 1
		playSound()
		clearCLI()
		println(streamName + " just went online")
	}
	return tempStreamInfo
}

func clearCLI() {
	var clear *exec.Cmd
	if hostOS == "windows" {
		clear = exec.Command("cmd", "/c", "cls")
	}
	if hostOS == "linux" {
		clear = exec.Command("clear")
	}
	clear.Stdout = os.Stdout
	clear.Run()
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
