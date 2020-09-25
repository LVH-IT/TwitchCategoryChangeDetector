package main

import (
	"fmt"
	"time"
)

func singleStream() {
	oldGameID = streamInfo.Data[0].GameID
	if len(oldGameInfo.Data) == 0 {
		oldGameName = "Game name not found"
	} else {
		oldGameName = oldGameInfo.Data[0].Name
		newGameID = streamInfo.Data[0].GameID
	}

	fmt.Println("----------------------------------------------------")
	fmt.Println("Channel to monitor: " + streamInfo.Data[0].DisplayName)
	fmt.Println("Current Category: " + oldGameName)

	startTime := time.Now()
	elapsed := time.Since(startTime)
	for oldGameID == newGameID && streamInfo.Data[0].Islive == true {
		for int(elapsed.Seconds()) < retryInterval {
			fmt.Printf("\rWaiting for change (Rechecking in %ds) ", (retryInterval - int(elapsed.Seconds())))
			time.Sleep(1e9) //sleep for 1000000000ns = 1000ms = 1s
			elapsed = time.Since(startTime)
		}
		routineDone = false
		go func() {
			newStreamInfo = getStreamInfoWithOnlineCheck()
			newGameID = newStreamInfo.Data[0].GameID
			routineDone = true
		}()
		for routineDone == false {
			if routineDone == false {
				fmt.Printf("\rWaiting for change (Rechecking      ) ")
				time.Sleep(2e8)
				if routineDone == false {
					fmt.Printf("\rWaiting for change (Rechecking  .   ) ")
					time.Sleep(2e8)
					if routineDone == false {
						fmt.Printf("\rWaiting for change (Rechecking  ..  ) ")
						time.Sleep(2e8)
						if routineDone == false {
							fmt.Printf("\rWaiting for change (Rechecking  ... ) ")
							time.Sleep(2e8)
							if routineDone == false {
								fmt.Printf("\rWaiting for change (Rechecking  ..  ) ")
								time.Sleep(2e8)
								if routineDone == false {
									fmt.Printf("\rWaiting for change (Rechecking  .   ) ")
									time.Sleep(2e8)
								}
							}
						}
					}
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
}
