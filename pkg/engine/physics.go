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
		input := *char.Input
		x, y := char.X, char.Y
		lx, ly := char.LX, char.LY
		char.Mu.Unlock()

		dir := helper.Vec2{}
		if input.Left {
			dir.X -= 1
		}
		if input.Right {
			dir.X += 1
		}
		if input.Up {
			dir.Y -= 1
		}
		if input.Down {
			dir.Y += 1
		}
		dir = dir.Normalize()

		moveDist := p.charSpeed * delta

		nx := x + dir.X*moveDist
		ny := y + dir.Y*moveDist

		if nx != lx || ny != ly {
			char.Mu.Lock()
			char.X, char.Y = nx, ny
			char.LX, char.LY = nx, ny
			char.Mu.Unlock()

			msg := []byte(fmt.Sprintf("move|%s|%.2f|%.2f", client.ID, nx, ny))
			select {
			case UpdateCh <- msg:
			default:
				logs.Warnf("Update Channel is full, dropping update for %s", client.ID)
			}
		}
	}
}
