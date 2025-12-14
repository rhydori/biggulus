package session

import (
	"sync"

	"github.com/rhydori/logs"
)

type Character struct {
	Input  *Input
	X, Y   float64
	LX, LY float64

	Mu sync.Mutex
}

type Input struct {
	Left  bool
	Right bool
	Up    bool
	Down  bool
}

func NewCharacter() *Character {
	return &Character{
		Input: &Input{},
	}
}

func (char *Character) HandleCharacter(msg []string) {
	action := msg[1]
	if action != "move" {
		logs.Warnf("HandleCharacter: Action not found - %s", action)
		return
	}
	switch action {
	case "move":
		char.UpdateCharacterXY(msg)
	}
}

func (char *Character) UpdateCharacterXY(msg []string) {
	char.Mu.Lock()
	defer char.Mu.Unlock()

	key := msg[2]
	state := msg[3]

	switch key {
	case "up":
		b := InputGetState(state)
		char.Input.Up = b
	case "down":
		b := InputGetState(state)
		char.Input.Down = b
	case "left":
		b := InputGetState(state)
		char.Input.Left = b
	case "right":
		b := InputGetState(state)
		char.Input.Right = b
	}

}

func InputGetState(state string) bool {
	switch state {
	case "pressed":
		return true
	case "released":
		return false
	default:
		logs.Warnf("InputGetState: State not found '%s'", state)
		return false
	}
}
