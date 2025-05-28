// package main
//
// import (
//
//	"gbccsclub/octopod-challenge/internal"
//	"github.com/gin-gonic/gin"
//	"github.com/joho/godotenv"
//	"log"
//	"os"
//	"strconv"
//
// )
//
//	func main() {
//		_ = godotenv.Load(".env")
//
//		router := gin.Default()
//		lobby := internal.NewLobby(10, 10)
//
//		router.GET("/", func(c *gin.Context) {
//			content := "Octapod Challenge Server" + "\n"
//			content += "-----------------------" + "\n"
//			content += "Maze:\n" + lobby.DisplayMaze("") + "\n"
//			content += "Octapods:\n"
//			for id, oct := range lobby.Octapods {
//				position := oct.position
//				content += id + " (" + strconv.Itoa(int(position.X())) + "," + strconv.Itoa(int(position.Y())) + ")\n"
//			}
//
//			c.String(200, content)
//		})
//		router.GET("/join", lobby.HandleJoin)
//		// For chron job on render to prevent sleep
//		router.GET("/ping", func(c *gin.Context) {
//			c.String(200, ".")
//		})
//
//		var port = "3000"
//		if os.Getenv("PORT") != "" {
//			port = os.Getenv("PORT")
//		}
//
//		log.Println("Starting a lobby server on port", port)
//		err := router.Run(":" + port)
//		if err != nil {
//			log.Fatal(err)
//		}
//	}
package main

import (
	"gbccsclub/octopod-challenge/internal/server"
	"gbccsclub/octopod-challenge/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	_ = godotenv.Load(".env")

	router := gin.Default()
	templ := web.NewTemplates()
	config := server.NewConfig()
	lobby := server.NewLobby(config)

	// ==================== Static Routes ====================

	router.GET("/", func(c *gin.Context) {
		templ.Render(c.Writer, "index", nil)
	})

	router.GET("/admin", func(c *gin.Context) {
		lobby.AdminHandler.HandleGetConfig(c, templ)
	})

	router.POST("/admin/update", func(c *gin.Context) {
		lobby.AdminHandler.HandleUpdateConfig(c, templ, lobby)
	})

	// ==================== Websocket Routes ====================

	router.GET("/join", func(c *gin.Context) {
		lobby.OctapodHandler.HandleJoin(c)
	})

	lobby.Start()

	var port = "3000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	log.Println("Starting a lobby server on port", port)
	err := router.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}
