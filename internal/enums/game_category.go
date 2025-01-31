package enums

import (
	"encoding/json"
	"fmt"
)

type GameCategory int

const (
	HEROES GameCategory = iota
	ITEMS
	ABILITIES
)

func (g GameCategory) String() string {
	return [...]string{"HEROES", "ITEMS", "ABILITIES"}[g]
}

func (g GameCategory) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.String())
}
func (g *GameCategory) UnmarshalJSON(data []byte) error {
	var categoryStr string
	if err := json.Unmarshal(data, &categoryStr); err != nil {
		return err
	}

	switch categoryStr {
	case "HEROES":
		*g = HEROES
	case "ITEMS":
		*g = ITEMS
	case "ABILITIES":
		*g = ABILITIES
	default:
		return fmt.Errorf("invalid GameState: %s", categoryStr)
	}
	return nil
}
