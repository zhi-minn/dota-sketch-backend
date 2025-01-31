package enums

import (
	"encoding/json"
	"fmt"
)

type GameRole int

const (
	PLAYER GameRole = iota
	ADMIN
)

func (g GameRole) String() string {
	return [...]string{"PLAYER", "ADMIN"}[g]
}

func (g GameRole) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.String())
}
func (g *GameRole) UnmarshalJSON(data []byte) error {
	var progressStr string
	if err := json.Unmarshal(data, &progressStr); err != nil {
		return err
	}

	switch progressStr {
	case "PLAYER":
		*g = PLAYER
	case "ADMIN":
		*g = ADMIN
	default:
		return fmt.Errorf("invalid GameProgress: %s", progressStr)
	}
	return nil
}
