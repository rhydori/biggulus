package session

import (
	"sync"

	"github.com/rhydori/biggulus/pkg/helper"
	"github.com/rhydori/logs"
)

type Character struct {
	*Input
	Position     helper.Vector2
	LastPosition helper.Vector2

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
		Input:        &Input{},
		Position:     helper.Vector2{},
		LastPosition: helper.Vector2{},
	}
}

func (char *Character) HandleCharacter(msg []string) {
	if len(msg) < 2 {
		logs.Warnf("HandleCharacter: malformed message: %v", msg)
		return
	}

	action := msg[1]
	if action != "move" {
		logs.Warnf("HandleCharacter: Action not found - %s", action)
		return
	}
	switch action {
	case "move":
		char.UpdateCharacterPosition(msg)
	}
}

func (char *Character) CharacterSnapshot() *Character {
	char.Mu.Lock()
	defer char.Mu.Unlock()

	inputCopy := &Input{
		Left:  char.Left,
		Right: char.Right,
		Up:    char.Up,
		Down:  char.Down,
	}

	return &Character{
		Input:        inputCopy,
		Position:     char.Position,
		LastPosition: char.LastPosition,
	}
}

func (char *Character) ApplyPosition(pos helper.Vector2) {
	char.Mu.Lock()
	defer char.Mu.Unlock()

	char.LastPosition = char.Position
	char.Position = pos
}

func (char *Character) UpdateCharacterPosition(msg []string) {
	if len(msg) < 4 {
		logs.Warnf("UpdateCharacterXY: malformed message: %v", msg)
		return
	}

	key := msg[2]
	state := msg[3]
	switch key {
	case "right", "left", "up", "down":
	default:
		logs.Warnf("UpdateCharacterPosition: InputKey not found: %s", key)
		return
	}
	switch state {
	case "pressed", "released":
	default:
		logs.Warnf("UpdateCharacterPosition: InputState not found: %s", state)
		return
	}

	val := state == "pressed"
	char.Mu.Lock()
	defer char.Mu.Unlock()
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
