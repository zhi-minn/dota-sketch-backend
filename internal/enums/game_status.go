package enums

import (
	"encoding/json"
	"fmt"
)

type GameStatus int

const (
	WAITING GameStatus = iota
	IN_PROGRESS
	ENDED
)

func (g GameStatus) String() string {
	return [...]string{"WAITING", "IN_PROGRESS", "ENDED"}[g]
}

func (g GameStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.String())
}
func (g *GameStatus) UnmarshalJSON(data []byte) error {
	var progressStr string
	if err := json.Unmarshal(data, &progressStr); err != nil {
		return err
	}

	switch progressStr {
	case "WAITING":
		*g = WAITING
	case "IN_PROGRESS":
		*g = IN_PROGRESS
	case "ENDED":
		*g = ENDED
	default:
		return fmt.Errorf("invalid GameProgress: %s", progressStr)
	}
	return nil
}
