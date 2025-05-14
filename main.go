package main

import (
	"gbccsclub/octopod-challenge/internal"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	_ = godotenv.Load(".env")

	router := gin.Default()
	lobby := internal.NewLobby(20, 20)

	router.GET("/join", lobby.HandleJoin)

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
