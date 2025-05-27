package handler

import (
	"gbccsclub/octopod-challenge/internal/server"
	"gbccsclub/octopod-challenge/internal/web"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strconv"
)

type AdminHandler struct {
	config *server.Config
}

func NewAdminHandler(c *server.Config) *AdminHandler {
	return &AdminHandler{
		config: c,
	}
}

func (ah *AdminHandler) HandleGetConfig(c *gin.Context, templ *web.Templates) {
	tickInterval, timeoutInterval, mazeSize := ah.config.Get()

	props := map[string]interface{}{
		"TickInterval":    tickInterval,
		"TimeoutInterval": timeoutInterval,
		"MazeSize":        mazeSize,
	}

	templ.Render(c.Writer, "admin", props)
}

func (ah *AdminHandler) HandleUpdateConfig(c *gin.Context, templ *web.Templates) {
	password := c.PostForm("password")
	hashedPassword := os.Getenv("HASHED_PASSWORD")
	if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) != nil {
		templ.Render(c.Writer, "error_message", map[string]interface{}{
			"Message": "Invalid password",
		})
		return
	}

	tickInterval, err := strconv.Atoi(c.PostForm("tick_interval"))
	if err != nil {
		templ.Render(c.Writer, "error_message", map[string]interface{}{
			"Message": "Invalid tick interval",
		})
		return
	}

	timeoutInterval, err := strconv.Atoi(c.PostForm("timeout_interval"))
	if err != nil {
		templ.Render(c.Writer, "error_message", map[string]interface{}{
			"Message": "Invalid timeout interval",
		})
		return
	}

	mazeSize, err := strconv.Atoi(c.PostForm("maze_size"))
	if err != nil {
		templ.Render(c.Writer, "error_message", map[string]interface{}{
			"Message": "Invalid maze size",
		})
		return
	}

	ah.config.Set(tickInterval, timeoutInterval, mazeSize)
	// TODO: restart lobby

	templ.Render(c.Writer, "success_message", map[string]interface{}{
		"Message": "Configuration updated successfully",
	})
}
