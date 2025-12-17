package session

import (
	"sync"

	"github.com/rhydori/biggulus/pkg/helper"
	"github.com/rhydori/biggulus/pkg/protocol"
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

func (char *Character) HandleCharacter(msg *protocol.Message) {
	if msg == nil {
		logs.Warnf("HandleCharacter: nil message received")
		return
	}

	action := msg.Action
	switch msg.Action {
	case "move":
		char.UpdateCharacterPosition(msg.Params)
	default:
		logs.Warnf("HandleCharacter: Action not found - %s", action)
		return
	}
}

func (char *Character) CharacterSnapshot() Character {
	char.Mu.Lock()
	defer char.Mu.Unlock()

	inputCopy := *char.Input
	return Character{
		Input:        &inputCopy,
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
	if len(msg) < 2 {
		logs.Warnf("UpdateCharacterPosition: Malformed message '%s'", msg)
		return
	}
	key := msg[0]
	state := msg[1]
	switch key {
	case "right", "left", "up", "down":
	default:
		logs.Warnf("UpdateCharacterPosition: Invalid Entity '%s'", key)
		return
	}
	switch state {
	case "pressed", "released":
	default:
		logs.Warnf("UpdateCharacterPosition: Invalid InputState '%s'", state)
		return
	}

	// NEED TO IMPLEMENT HEARTBEAT ON PRESSED IN CASE CLIENT CRASHES WHILE PRESSING
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
