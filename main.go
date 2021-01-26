package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
)

var (
	//Some globally used variables
	bearerToken                string   //
	clientID                   string   //
	soundFile                  string   //
	clientSecret               string   //
	useCategoryWhitelist       bool     //loaded via loadConfig()
	categories                 []string //
	notifyOnOfflineTitleChange bool     //
	notifyOnOnlineTitleChange  bool     //

	streamInfo       stream
	routineDone      bool
	wasOnline        bool = false
	oldGameName      string
	newGameID        string
	oldGameID        string
	oldGameInfo      game
	newGameName      string
	newStreamInfo    stream
	whitelistGameIDs []string
	oldTitle         string
	newTitle         string
	initOnline       = 6
	initOffline      = 6
	hostOS           = runtime.GOOS
	silentMode       = false
	err              error
	//errLogFile *os.File

	//Definitions of all the flags
	streamName    string
	retryInterval int = 10

	//Discord Bot
	discordMode     bool
	discordBotToken string
	bot             *discordgo.Session
	dcChannel       dcChannels
)

func main() {
	checkOS()
	parseFlags()
	loadConfig()
	checkAPIToken()

	if discordMode {
		startDiscordBot()
	}

	go func() {
		ctrlC()
	}()
	if !discordMode {
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
	} else {
		for {
			time.Sleep(1e9)
		}
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
			if !silentMode {
				clearCLI()
				println(streamName + " just went offline. Waiting for change (Checking every " + fmt.Sprint(retryInterval) + "s)")
				playSound()
			}
			if discordMode {
				dcChannel.Session.ChannelMessageSend(dcChannel.ChannelID, streamName+" just went offline")
			}
			for !tempStreamInfo.Data[0].Islive {
				waitRetryInterval()
				if initOffline == 6 {
					tempStreamInfo = getStreamInfo()
					offlineNewTitle = tempStreamInfo.Data[0].Title
					if offlineOldTitle != offlineNewTitle {
						if !silentMode {
							if notifyOnOfflineTitleChange {
								playSound()
							}
							clearCLI()
							println("Title changed to: " + offlineNewTitle)
							println("----------------------------------------------------")
							println(streamName + " just went offline. Waiting for change (Checking every " + fmt.Sprint(retryInterval) + "s)")
						}

						if notifyOnOfflineTitleChange && discordMode {
							dcChannel.Session.ChannelMessageSend(dcChannel.ChannelID, streamName+" - Offline title changed to: "+offlineNewTitle)
						}
						offlineOldTitle = offlineNewTitle
					}
				}
				if initOffline < 6 {
					initOffline++
				}
			}
		} else {
			if !silentMode {
				println(streamName + " is currently offline. Waiting for change (Checking every " + fmt.Sprint(retryInterval) + "s)")
			}
			for !tempStreamInfo.Data[0].Islive {
				waitRetryInterval()
				if initOffline == 6 {
					tempStreamInfo = getStreamInfo()
					offlineNewTitle = tempStreamInfo.Data[0].Title
					if offlineOldTitle != offlineNewTitle {
						if !silentMode {
							if notifyOnOfflineTitleChange {
								playSound()
							}
							clearCLI()
							println("Title changed to: " + offlineNewTitle)
							println("----------------------------------------------------")
							println(streamName + " is currently offline. Waiting for change (Checking every " + fmt.Sprint(retryInterval) + "s)")
						}

						if notifyOnOfflineTitleChange && discordMode {
							dcChannel.Session.ChannelMessageSend(dcChannel.ChannelID, streamName+" - Offline title changed to: "+offlineNewTitle)
						}

						offlineOldTitle = offlineNewTitle
					}
				}
				if initOffline < 6 {
					initOffline++
				}
			}
			wasOnline = true
		}
		initOffline = 6
		initOnline = 6
		if !silentMode {
			playSound()
			clearCLI()
			println(streamName + " just went online")
		}

		if discordMode {
			dcChannel.Session.ChannelMessageSend(dcChannel.ChannelID, streamName+" just went online")
		}

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
	flag.BoolVar(&discordMode, "dcbot", false, "used for discord bot implementation")
	flag.Parse()

	if discordMode {
		silentMode = true
	}
}

func checkError(errorVar error) {
	if errorVar != nil {
		log.Fatal(errorVar)
		os.Exit(1)
		/*logAgain:
		if _, err := os.Stat("error.log"); err == nil {
			errLogFile, err = os.Open("error.log")
			if err == nil {
				log.SetOutput(errLogFile)
				log.Println(errorVar)
				errLogFile.Sync()
			}
		} else if os.IsNotExist(err) {
			errLogFile, err = os.Create("error.log")
			goto logAgain
		}*/
	}
}

func ctrlC() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	println("Exiting...")
	bot.Close()
	os.Exit(1)
}
