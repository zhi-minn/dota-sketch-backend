package services

import (
	"github.com/zhi-minn/dota-sketch-backend/internal/enums"
	"github.com/zhi-minn/dota-sketch-backend/internal/models"
	"log"
	"math/rand"
	"time"
)

type GameService struct{}

func NewGameService() *GameService {
	return &GameService{}
}

func (gs *GameService) CreateGame(req models.GameRequest) *models.Game {
	gameCode := generateGameCode()

	settings := &models.GameSettings{
		Categories:            req.Categories,
		DrawingTime:           req.DrawingTime,
		Rounds:                req.Rounds,
		MaxPlayers:            req.MaxPlayers,
		FirstGuessDelay:       req.FirstGuessDelay,
		ReduceTimeWhenGuessed: req.ReduceTimeWhenGuessed,
		AllowWordReroll:       req.AllowWordReroll,
		Public:                req.Public,
	}

	game := &models.Game{
		Code:      gameCode,
		Players:   make(map[string]*models.Player),
		Status:    enums.WAITING,
		Rounds:    []models.Round{},
		Settings:  settings,
		CreatedAt: time.Now(),
	}

	models.Mutex.Lock()
	defer models.Mutex.Unlock()
	models.ActiveGames[gameCode] = game
	log.Println("Created game", game)

	return game
}

func generateGameCode() string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	code := make([]byte, 4)
	for i := range code {
		code[i] = letters[rand.Intn(len(letters))]
	}
	return string(code)
}
