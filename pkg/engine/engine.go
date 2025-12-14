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
	e.cs.Mu.Lock()
	clients := make([]*session.Client, 0, len(e.cs.Clients))
	for _, c := range e.cs.Clients {
		clients = append(clients, c)
	}
	e.cs.Mu.Unlock()

	for _, client := range clients {
		dir := helper.Vec2{}
		if client.Input.Left {
			dir.X -= 1
		}
		if client.Input.Right {
			dir.X += 1
		}
		if client.Input.Up {
			dir.Y -= 1
		}
		if client.Input.Down {
			dir.Y += 1
		}
		dir = dir.Normalize()
		client.X += dir.X * e.speed
		client.Y += dir.Y * e.speed
		if client.X != client.LastX || client.Y != client.LastY {
			msg := []byte(fmt.Sprintf("move|%s|%f|%f", client.ID, client.X, client.Y))
			select {
			case e.UpdateCh <- msg:
			default:
				logs.Warnf("Update Channel is full, dropping update for %s", client.ID)
			}

			client.LastX = client.X
			client.LastY = client.Y
		}
	}
}
