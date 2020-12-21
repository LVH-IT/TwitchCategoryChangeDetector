package main

import (
	"fmt"
	"time"
)

func singleStream() {
	oldGameID = streamInfo.Data[0].GameID
	oldTitle = streamInfo.Data[0].Title
	newTitle = oldTitle
	if len(oldGameInfo.Data) == 0 {
		oldGameName = "Game name not found"
	} else {
		oldGameName = oldGameInfo.Data[0].Name
		newGameID = streamInfo.Data[0].GameID
	}
	println("----------------------------------------------------")
	println("Channel to monitor: " + streamInfo.Data[0].DisplayName)
	println("Current Category: " + oldGameName)
	println("Current Title: " + oldTitle)

	startTime := time.Now()
	elapsed := time.Since(startTime)
	newStreamInfo = streamInfo
	for oldGameID == newGameID && newStreamInfo.Data[0].Islive && oldTitle == newTitle {
		for int(elapsed.Seconds()) < retryInterval {
			fmt.Printf("\rWaiting for change (Rechecking in %ds) ", (retryInterval - int(elapsed.Seconds())))
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
					time.Sleep(2e8)
				}
			}
		}

		startTime = time.Now()
		elapsed = time.Since(startTime)
	}

	clearCLI()
	soundPlayed := false

	if oldGameID != newGameID {
		if useCategoryWhitelist {
			newGameInfo := getGameInfo(newStreamInfo)
			for _, b := range whitelistGameIDs {
				if newGameID == b {
					playSound()
					soundPlayed = true
					continue
				}
			}
			if len(newGameInfo.Data) == 0 {
				newGameName = "Game name not found"
			} else {
				newGameName = newGameInfo.Data[0].Name
			}
			println("Category changed to: " + newGameName)
		} else {
			newGameInfo := getGameInfo(newStreamInfo)
			if len(newGameInfo.Data) == 0 {
				newGameName = "Game name not found"
			} else {
				newGameName = newGameInfo.Data[0].Name
			}

			println("Category changed to: " + newGameName)
			playSound()
			soundPlayed = true
		}
	}
	if oldTitle != newTitle {
		println("Title changed to: " + newTitle)
		if !soundPlayed && notifyOnOnlineTitleChange {
			playSound()
		}
	}
}
