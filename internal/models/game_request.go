package models

import "github.com/zhi-minn/dota-sketch-backend/internal/enums"

type GameRequest struct {
	Categories            []enums.GameCategory `json:"categories"`
	DrawingTime           int                  `json:"drawingTime"`
	Rounds                int                  `json:"rounds"`
	MaxPlayers            int                  `json:"maxPlayers"`
	FirstGuessDelay       int                  `json:"firstGuessDelay"`
	ReduceTimeWhenGuessed bool                 `json:"reduceTimeWhenGuessed"`
	AllowWordReroll       bool                 `json:"allowWordReroll"`
	Public                bool                 `json:"public"`
}
