package main

import (
	"fmt"
	"time"
)

func singleStream() {
	oldGameID = streamInfo.Data[0].GameID
	oldTitle := streamInfo.Data[0].Title
	newTitle := oldTitle
	if len(oldGameInfo.Data) == 0 {
		oldGameName = "Game name not found"
	} else {
		oldGameName = oldGameInfo.Data[0].Name
		newGameID = streamInfo.Data[0].GameID
	}

	fmt.Println("----------------------------------------------------")
	fmt.Println("Channel to monitor: " + streamInfo.Data[0].DisplayName)
	fmt.Println("Current Category: " + oldGameName)
	fmt.Println("Current Title: " + oldTitle)

	startTime := time.Now()
	elapsed := time.Since(startTime)
	for oldGameID == newGameID && streamInfo.Data[0].Islive && oldTitle == newTitle {
		for int(elapsed.Seconds()) < retryInterval {
			fmt.Printf("\rWaiting for change (Rechecking in %ds) ", (retryInterval - int(elapsed.Seconds())))
			time.Sleep(1e9) //sleep for 1000000000ns = 1000ms = 1s
			elapsed = time.Since(startTime)
		}
		routineDone = false
		go func() {
			newStreamInfo = getStreamInfoWithOnlineCheck()
			newGameID = newStreamInfo.Data[0].GameID
			if notifyOnOnlineTitleChange {
				newTitle = newStreamInfo.Data[0].Title
			}
			routineDone = true
		}()

		for routineDone == false {
			for i := 1; i <= 6; i++ {
				if routineDone == false {
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
	if oldGameID != newGameID {
		if useCategoryWhitelist {
			for _, b := range whitelistGameIDs {
				if newGameID == b {
					newGameInfo := getGameInfo(newStreamInfo)
					if len(newGameInfo.Data) == 0 {
						newGameName = "Game name not found"
					} else {
						newGameName = newGameInfo.Data[0].Name
					}
					println()
					fmt.Printf("Category changed to: " + newGameName + "\n")
					playSound()
					return
				}
			}
			newGameInfo := getGameInfo(newStreamInfo)
			if len(newGameInfo.Data) == 0 {
				newGameName = "Game name not found"
			} else {
				newGameName = newGameInfo.Data[0].Name
			}
			println()
			fmt.Printf("Category changed to: " + newGameName + "\n")
		} else {
			newGameInfo := getGameInfo(newStreamInfo)
			if len(newGameInfo.Data) == 0 {
				newGameName = "Game name not found"
			} else {
				newGameName = newGameInfo.Data[0].Name
			}
			println()
			fmt.Printf("Category changed to: " + newGameName + "\n")
			playSound()
		}
	}
	if oldTitle != newTitle && notifyOnOnlineTitleChange {
		println()
		fmt.Printf("Title changed to: " + newTitle + "\n")
		playSound()
	}
}
