package server

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

type DiscordBot struct {
	session   *discordgo.Session
	channelId string
}

func NewDiscordBot(discordServer string, discordChannelId string) *DiscordBot {
	token := discordServer

	session, err := discordgo.New("Bot " + token)
	log.Println("New Discord bot created")
	if err != nil {
		panic(err)
	}

	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent

	bot := &DiscordBot{
		session:   session,
		channelId: discordChannelId,
	}

	if err = session.Open(); err != nil {
		panic(err)
	}

	return bot
}

func (d *DiscordBot) Close() {
	err := d.session.Close()
	if err != nil {
		log.Printf("Failed to close Discord session: %v\n", err)
		return
	}
}

func (d *DiscordBot) SendMessage(message string) {
	_, err := d.session.ChannelMessageSend(d.channelId, message)
	if err != nil {
		log.Printf("Failed to send message to Discord: %v\n", err)
		return
	}
}
