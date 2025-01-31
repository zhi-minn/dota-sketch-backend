package models

import (
	"github.com/zhi-minn/dota-sketch-backend/internal/enums"
	"sync"
	"time"
)

type Game struct {
	Code      string             `json:"code"`
	Players   map[string]*Player `json:"players"`
	Status    enums.GameStatus   `json:"status"`
	Rounds    []Round            `json:"rounds"`
	Settings  *GameSettings      `json:"settings"`
	CreatedAt time.Time          `json:"createdAt"`
}

type Player struct {
	Nickname string
	Score    int
	Role     enums.GameRole
}

type Round struct {
	Word       string
	Sketcher   string
	RoundState string
}

type GameSettings struct {
	Categories            []enums.GameCategory
	DrawingTime           int
	Rounds                int
	MaxPlayers            int
	FirstGuessDelay       int
	ReduceTimeWhenGuessed bool
	AllowWordReroll       bool
	Public                bool
}

var ActiveGames = make(map[string]*Game)
var Mutex sync.Mutex
