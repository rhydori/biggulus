package engine

import (
	"fmt"

	"github.com/rhydori/biggulus/pkg/helper"
	"github.com/rhydori/biggulus/pkg/session"
	"github.com/rhydori/logs"
)

type Physics struct {
	charSpeed float64
}

func NewPhysics(charSpeed float64) *Physics {
	return &Physics{
		charSpeed: charSpeed,
	}
}

func (p *Physics) ProcessMovement(clients []*session.Client, delta float64, UpdateCh chan<- []byte) {
	for _, client := range clients {
		char := client.Char

		char.Mu.Lock()
		dir := helper.Vec2{}
		if char.Input.Left {
			dir.X -= 1
		}
		if char.Input.Right {
			dir.X += 1
		}
		if char.Input.Up {
			dir.Y -= 1
		}
		if char.Input.Down {
			dir.Y += 1
		}

		dir = dir.Normalize()

		moveDist := p.charSpeed * delta

		char.X += dir.X * moveDist
		char.Y += dir.Y * moveDist

		if char.X != char.LX || char.Y != char.LY {
			msg := []byte(fmt.Sprintf("move|%s|%.2f|%.2f", client.ID, char.X, char.Y))
			select {
			case UpdateCh <- msg:
			default:
				logs.Warnf("Update Channel is full, dropping update for %s", client.ID)
			}

			char.LX = char.X
			char.LY = char.Y
		}
		char.Mu.Unlock()
	}
}
