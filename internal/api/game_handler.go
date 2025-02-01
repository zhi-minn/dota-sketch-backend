package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zhi-minn/dota-sketch-backend/internal/enums"
	"github.com/zhi-minn/dota-sketch-backend/internal/models"
	"github.com/zhi-minn/dota-sketch-backend/internal/services"
	"net/http"
	"time"
)

type GameHandler struct {
	gameService *services.GameService
}

func NewGameHandler(gameService *services.GameService) *GameHandler {
	return &GameHandler{gameService}
}

func (h *GameHandler) CreateGameHandler(c *gin.Context) {
	var req models.GameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid game creation request"})
		return
	}

	game := h.gameService.CreateGame(req)

	c.JSON(http.StatusOK, game)
}

func (h *GameHandler) JoinGame(c *gin.Context) {
	gameCode := c.Param("gameCode")
	game, exists := models.ActiveGames[gameCode]
	if !exists {
		c.JSON(400, gin.H{"error": "Game does not exist"})
	}

	var data map[string]string
	if err := json.NewDecoder(c.Request.Body).Decode(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	nick, exists := data["nick"] // Extract the "nick" field
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "'nick' field is required"})
		return
	}

	_, exists = game.Players[nick]
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nickname in use"})
		return
	}

	playerUid := uuid.New().String()
	game.Players[playerUid] = &models.Player{
		Nickname: nick,
		Score:    0,
		Role:     enums.PLAYER,
	}
	c.JSON(http.StatusOK, gin.H{"code": game.Code, "uid": playerUid})
}

func (h *GameHandler) LiveGamesHandler(c *gin.Context) {
	result := make(map[string]string)
	for gameCode, game := range models.ActiveGames {
		result[gameCode] = game.CreatedAt.Format(time.DateTime)
	}

	c.JSON(http.StatusOK, result)
}
