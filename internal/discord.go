package internal

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"strings"
)

type DiscordBot struct {
	Session   *discordgo.Session
	ChannelId string
	Lobby     *Lobby // Add reference to Lobby
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

	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent

	bot := &DiscordBot{
		Session:   session,
		ChannelId: os.Getenv("DISCORD_CHANNEL_ID"),
	}

	session.AddHandler(bot.makeMessageHandler())

	if err = session.Open(); err != nil {
		log.Fatalf("Failed to open Discord session: %v", err)
	}

	return bot
}

func (d *DiscordBot) makeMessageHandler() func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		parts := strings.Fields(m.Content)

		if len(parts) >= 2 && parts[0] == "!where" {
			id := parts[1]
			log.Printf("Received Where command from discord. ID: %s", id)

			if d.Lobby == nil {
				s.ChannelMessageSend(m.ChannelID, "Lobby not initialized.")
				return
			}

			mazeDisplay := d.Lobby.DisplayMaze(id)
			_, err := s.ChannelMessageSend(m.ChannelID, mazeDisplay)
			if err != nil {
				log.Printf("Error sending maze display: %v", err)
			}
			return
		}

		if len(parts) == 1 && parts[0] == "!where" {
			_, err := s.ChannelMessageSend(m.ChannelID, "Usage: `!where <ID>`")
			if err != nil {
				log.Printf("Error sending usage: %v", err)
			}
		}
	}
}

func (d *DiscordBot) SetLobby(lobby *Lobby) {
	d.Lobby = lobby
}

func (d *DiscordBot) Close() {
	err := d.Session.Close()
	if err != nil {
		return
	}
}

func (d *DiscordBot) SendMessage(message string) {
	_, err := d.Session.ChannelMessageSend(d.ChannelId, message)
	if err != nil {
		panic(err)
	}
}
