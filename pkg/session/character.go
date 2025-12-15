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
	case "right", "left", "up", "down":
	default:
		logs.Warnf("UpdateCharacterXY: InputKey not found: %s", key)
		return
	}

	switch state {
	case "pressed", "released":
	default:
		logs.Warnf("UpdateCharacterXY: InputState not found: %s", state)
		return
	}

	val := state == "pressed"
	switch key {
	case "up":
		char.Input.Up = val
	case "down":
		char.Input.Down = val
	case "left":
		char.Input.Left = val
	case "right":
		char.Input.Right = val
	}

}
