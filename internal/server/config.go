package server

import (
	"os"
	"sync"
)

type Config struct {
	mu sync.RWMutex

	// Loop
	TickInterval int

	// Competition settings
	MaxExplorationSteps int
	MaxSolvingSteps     int

	// Maze
	MazeSize int

	// Discord
	DiscordBotToken  string
	DiscordChannelId string
}

func NewConfig() *Config {
	return &Config{
		TickInterval:        3000,
		MaxExplorationSteps: 2 * 10 * 10,
		MaxSolvingSteps:     5 * 10,
		MazeSize:            10,
		DiscordBotToken:     os.Getenv("DISCORD_BOT_TOKEN"),
		DiscordChannelId:    os.Getenv("DISCORD_CHANNEL_ID"),
	}
}
