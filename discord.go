package main

import (
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func startDiscordBot() {
	println("Starting in Discord Bot mode")
	bot, err = discordgo.New("Bot " + discordBotToken)
	if err != nil {
		println("error creating Discord session,", err)
		os.Exit(1)
	}

	bot.AddHandler(messageCreate)
	bot.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsDirectMessages)

	err = bot.Open()
	if err != nil {
		println("error opening discord bot connection,", err)
		os.Exit(1)
	}

	println("Discord Bot is running.  Press CTRL-C to exit.")
}

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}
	if strings.HasPrefix(message.Content, "!start") {
		if dcChannel.ChannelID == "" {
			dcChannel.Session = session
			dcChannel.ChannelID = message.ChannelID
			dcChannel.StreamName = strings.TrimPrefix(message.Content, "!start")
			streamName = strings.ReplaceAll(dcChannel.StreamName, " ", "")
			session.ChannelMessageSend(message.ChannelID, "Now monitoring "+streamName)

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
			session.ChannelMessageSend(message.ChannelID, "You already monitor someone")
		}
	}
}
