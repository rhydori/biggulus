package engine

import (
	"fmt"
	"math"

	"github.com/rhydori/biggulus/pkg/helper"
	"github.com/rhydori/biggulus/pkg/protocol"
	"github.com/rhydori/biggulus/pkg/session"
	"github.com/rhydori/logs"
)

const eps = 0.01

type Physics struct {
	charSpeed float64
}

func NewPhysics(charSpeed float64) *Physics {
	return &Physics{
		charSpeed: charSpeed,
	}
}

func (p *Physics) ProcessMovement(clients []*session.Client, delta float64, UpdateCh chan<- []byte) {
	if delta <= 0 {
		return
	}
	for _, client := range clients {
		if client == nil || client.Char == nil {
			continue
		}
		charSnap := client.Char.CharacterSnapshot()

		input := charSnap.Input
		pos := charSnap.Position
		last := charSnap.LastPosition

		dir := helper.Vector2{}
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

		move := p.charSpeed * delta

		nx := pos.X + dir.X*move
		ny := pos.Y + dir.Y*move

		if math.Abs(nx-last.X) > eps || math.Abs(ny-last.Y) > eps {
			newPos := helper.Vector2{X: nx, Y: ny}
			client.Char.ApplyPosition(newPos)

			xStr := fmt.Sprintf("%.2f", nx)
			yStr := fmt.Sprintf("%.2f", ny)
			msg := protocol.CreateMessageBytes("character", "move", client.ID, xStr, yStr)
			select {
			case UpdateCh <- msg:
			default:
				logs.Warnf("Update Channel is full, dropping update for %s", client.ID)
			}
		}
	}
}
