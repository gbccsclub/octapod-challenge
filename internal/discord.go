package internal

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
)

type DiscordBot struct {
	Session   *discordgo.Session
	ChannelId string
}

func NewDiscordBot() *DiscordBot {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		panic("DISCORD_BOT_TOKEN is not set")
	}
	session, err := discordgo.New("Bot " + token)
	log.Println("New Discord bot created")
	if err != nil {
		panic(err)
	}

	return &DiscordBot{
		Session:   session,
		ChannelId: os.Getenv("DISCORD_CHANNEL_ID"),
	}
}

func (d *DiscordBot) Open() {
	err := d.Session.Open()
	if err != nil {
		panic(err)
	}
}

func (d *DiscordBot) Close() {
	d.Session.Close()
}

func (d *DiscordBot) SendMessage(message string) {
	_, err := d.Session.ChannelMessageSend(d.ChannelId, message)
	if err != nil {
		panic(err)
	}
}
