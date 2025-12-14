package engine

import (
	"fmt"
	"time"

	"github.com/rhydori/biggulus/pkg/helper"
	"github.com/rhydori/biggulus/pkg/session"
	"github.com/rhydori/logs"
)

type Engine struct {
	tick  time.Duration
	speed float64

	cs       *session.ClientStore
	UpdateCh chan []byte
}

func NewEngine(tickInterval time.Duration, characterSpeed float64, cs *session.ClientStore) *Engine {
	return &Engine{
		tick:     tickInterval,
		speed:    characterSpeed * tickInterval.Seconds(),
		cs:       cs,
		UpdateCh: make(chan []byte, 1024),
	}
}

func (e *Engine) StartEngine() {
	logs.Info("Engine started at ", e.tick)

	ticker := time.NewTicker(e.tick)
	for range ticker.C {
		e.updateLoop()
	}
}

func (e *Engine) updateLoop() {
	clients := e.cs.ClientStoreSnapshot()

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
		char.X += dir.X * e.speed
		char.Y += dir.Y * e.speed
		if char.X != char.LX || char.Y != char.LY {
			msg := []byte(fmt.Sprintf("move|%s|%f|%f", client.ID, char.X, char.Y))
			select {
			case e.UpdateCh <- msg:
			default:
				logs.Warnf("Update Channel is full, dropping update for %s", client.ID)
			}

			char.LX = char.X
			char.LY = char.Y
		}
		char.Mu.Unlock()
	}
}
