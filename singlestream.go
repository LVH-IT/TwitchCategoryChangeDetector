package main

import (
	crand "crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func singleStream() {
	oldGameID = streamInfo.Data[0].GameID
	oldTitle = streamInfo.Data[0].Title
	newTitle = oldTitle

	if len(oldGameInfo.Data) == 0 {
		oldGameName = "Game name not found"
		newGameID = "1"
	} else {
		oldGameName = oldGameInfo.Data[0].Name
		newGameID = streamInfo.Data[0].GameID
	}

	if !silentMode {
		println("----------------------------------------------------")
		println("Channel to monitor: " + streamInfo.Data[0].DisplayName)
		println("Current Category: " + oldGameName)
		println("Current Title: " + oldTitle)
	}

	startTime := time.Now()
	elapsed := time.Since(startTime)
	newStreamInfo = streamInfo
	for oldGameID == newGameID && newStreamInfo.Data[0].Islive && oldTitle == newTitle {
		for int(elapsed.Seconds()) < retryInterval {
			if !silentMode {
				fmt.Printf("\rWaiting for change (Rechecking in %ds) ", (retryInterval - int(elapsed.Seconds())))
			}
			time.Sleep(1e9) //sleep for 1000000000ns = 1000ms = 1s
			elapsed = time.Since(startTime)
		}
		routineDone = false

		go func() {
			if initOnline == 6 {
				newStreamInfo = getStreamInfo()
				newGameID = newStreamInfo.Data[0].GameID
				newTitle = newStreamInfo.Data[0].Title
			}
			if initOnline < 6 {
				initOnline++
			}
			routineDone = true
		}()

		for !routineDone {
			for i := 1; i <= 6; i++ {
				if !routineDone {
					if !silentMode {
						if i == 1 {
							fmt.Printf("\rWaiting for change (Rechecking      ) ")
						}
						if i == 2 || i == 6 {
							fmt.Printf("\rWaiting for change (Rechecking  .   ) ")
						}
						if i == 3 || i == 5 {
							fmt.Printf("\rWaiting for change (Rechecking  ..  ) ")
						}
						if i == 4 {
							fmt.Printf("\rWaiting for change (Rechecking  ... ) ")
						}
					}
					time.Sleep(2e8)
				}
			}
		}

		startTime = time.Now()
		elapsed = time.Since(startTime)

	}

	if !silentMode {
		clearCLI()
	}
	soundPlayed := false

	var embedMsg discordgo.MessageEmbed
	var field1 discordgo.MessageEmbedField //Field 1 is for Category changes
	var field2 discordgo.MessageEmbedField //Field 2 is for Title changes

	if oldGameID != newGameID {
		if useCategoryWhitelist {
			newGameInfo := getGameInfo(newStreamInfo)
			for _, b := range whitelistGameIDs {
				if newGameID == b {
					if !silentMode {
						playSound()
					}
					if discordMode {
						//dcChannel.Session.ChannelMessageSend(dcChannel.ChannelID, streamName+" - Category changed to: "+newGameName)
						field1.Inline = true
						field1.Name = "Category changed to"
						field1.Value = newGameName
						embedMsg.Fields = append(embedMsg.Fields, &field1)
					}
					soundPlayed = true
					continue
				}
			}
			if len(newGameInfo.Data) == 0 {
				newGameName = "Game name not found"
			} else {
				newGameName = newGameInfo.Data[0].Name
			}
			if !silentMode {
				println("Category changed to: " + newGameName)
			}

		} else {
			newGameInfo := getGameInfo(newStreamInfo)
			if len(newGameInfo.Data) == 0 {
				newGameName = "Game name not found"
			} else {
				newGameName = newGameInfo.Data[0].Name
			}
			if !silentMode {
				println("Category changed to: " + newGameName)
				playSound()
			}
			if discordMode {
				//dcChannel.Session.ChannelMessageSend(dcChannel.ChannelID, streamName+" - Category changed to: "+newGameName)
				field1.Inline = true
				field1.Name = "Category changed to"
				field1.Value = newGameName
				embedMsg.Fields = append(embedMsg.Fields, &field1)
			}
			soundPlayed = true
		}
	}
	if oldTitle != newTitle {
		if !silentMode {
			println("Title changed to: " + newTitle)
			if !soundPlayed && notifyOnOnlineTitleChange {
				playSound()
			}
		}
		if discordMode {
			//dcChannel.Session.ChannelMessageSend(dcChannel.ChannelID, streamName+" - Title changed to: "+newTitle)
			field2.Inline = true
			field2.Name = "Title changed to"
			field2.Value = newTitle
			embedMsg.Fields = append(embedMsg.Fields, &field2)
		}
	}
	if discordMode && newStreamInfo.Data[0].Islive {
		var thumbnail discordgo.MessageEmbedThumbnail
		trash := make([]byte, 8)
		crand.Read(trash)
		thumbnail.URL = "https://static-cdn.jtvnw.net/previews-ttv/live_user_" + streamName + "-1280x720.jpg" + "?trash=" + hex.EncodeToString(trash)
		thumbnail.Width = 1280
		thumbnail.Height = 720
		embedMsg.Thumbnail = &thumbnail

		embedMsg.Title = streamName

		_, err = dcChannel.Session.ChannelMessageSendEmbed(dcChannel.ChannelID, &embedMsg)
		checkError(err)
	}
}
