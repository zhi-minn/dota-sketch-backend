package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zhi-minn/dota-sketch-backend/internal/api"
	"github.com/zhi-minn/dota-sketch-backend/internal/services"
	"github.com/zhi-minn/dota-sketch-backend/internal/ws"
	"github.com/zhi-minn/dota-sketch-backend/pkg/middleware"
	"log"
)

func main() {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	gameService := services.NewGameService()
	gameHandler := api.NewGameHandler(gameService)
	r.POST("/api/games", gameHandler.CreateGameHandler)
	r.POST("/api/games/:gameCode", gameHandler.JoinGame)

	r.GET("/metrics/live-games", gameHandler.LiveGamesHandler)

	hub := ws.NewHub()
	r.GET("/ws", ws.ServeWs(hub))

	log.Println("Listening on port 8081")
	r.Run(":8081")
}
