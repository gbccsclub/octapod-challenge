package server

import (
	"gbccsclub/octopod-challenge/internal/web"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strconv"
)

type AdminHandler struct {
	config *Config
}

func NewAdminHandler(c *Config) *AdminHandler {
	return &AdminHandler{
		config: c,
	}
}

func (ah *AdminHandler) HandleGetConfig(c *gin.Context, templ *web.Templates) {
	ah.config.mu.Lock()
	defer ah.config.mu.Unlock()

	props := map[string]interface{}{
		"TickInterval":        ah.config.TickInterval,
		"MazeSize":            ah.config.MazeSize,
		"DiscordChannelId":    ah.config.DiscordChannelId,
		"MaxExplorationSteps": ah.config.MaxExplorationSteps,
		"MaxSolvingSteps":     ah.config.MaxSolvingSteps,
	}

	templ.Render(c.Writer, "admin", props)
}

func (ah *AdminHandler) HandleUpdateConfig(c *gin.Context, templ *web.Templates, lobby *Lobby) bool {
	password := c.PostForm("password")
	hashedPassword := os.Getenv("HASHED_PASSWORD")
	if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) != nil {
		templ.Render(c.Writer, "error_message", map[string]interface{}{
			"Message": "Invalid password",
		})
		return false
	}

	tickInterval, err := strconv.Atoi(c.PostForm("tick_interval"))
	if err != nil {
		templ.Render(c.Writer, "error_message", map[string]interface{}{
			"Message": "Invalid tick interval",
		})
		return false
	}

	maxExplorationSteps, err := strconv.Atoi(c.PostForm("max_exploration_steps"))
	if err != nil {
		templ.Render(c.Writer, "error_message", map[string]interface{}{
			"Message": "Invalid max exploration steps",
		})
		return false
	}

	maxSolvingSteps, err := strconv.Atoi(c.PostForm("max_solving_steps"))
	if err != nil {
		templ.Render(c.Writer, "error_message", map[string]interface{}{
			"Message": "Invalid max solving steps",
		})
		return false
	}

	mazeSize, err := strconv.Atoi(c.PostForm("maze_size"))
	if err != nil {
		templ.Render(c.Writer, "error_message", map[string]interface{}{
			"Message": "Invalid maze size",
		})
		return false
	}

	discordBotToken := c.PostForm("discord_bot_token")
	if discordBotToken == "" {
		discordBotToken = os.Getenv("DISCORD_BOT_TOKEN")
	}

	discordChannelId := c.PostForm("discord_channel_id")

	ah.config.mu.Lock()
	defer ah.config.mu.Unlock()

	ah.config.TickInterval = tickInterval
	ah.config.MazeSize = mazeSize
	ah.config.DiscordBotToken = discordBotToken
	ah.config.DiscordChannelId = discordChannelId
	ah.config.MaxExplorationSteps = maxExplorationSteps
	ah.config.MaxSolvingSteps = maxSolvingSteps

	lobby.RequestRestart()

	templ.Render(c.Writer, "success_message", map[string]interface{}{
		"Message": "Configuration updated successfully",
	})

	return true
}
