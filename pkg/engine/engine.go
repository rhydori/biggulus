package engine

import (
	"time"

	"github.com/rhydori/biggulus/pkg/session"
	"github.com/rhydori/logs"
)

type Engine struct {
	tick time.Duration

	cs   *session.ClientStore
	phys *Physics

	UpdateCh chan []byte
}

func NewEngine(tickInterval time.Duration, cs *session.ClientStore, physic *Physics) *Engine {
	return &Engine{
		tick:     tickInterval,
		cs:       cs,
		phys:     physic,
		UpdateCh: make(chan []byte, 256),
	}
}

func (e *Engine) StartEngine() {
	logs.Info("Engine started at ", e.tick)

	ticker := time.NewTicker(e.tick)
	defer ticker.Stop()

	lastTime := time.Now()

	for range ticker.C {
		now := time.Now()
		delta := now.Sub(lastTime).Seconds()
		lastTime = now

		e.updateLoop(delta)
	}
}

func (e *Engine) updateLoop(delta float64) {
	clients := e.cs.ClientStoreSnapshot()

	e.phys.ProcessMovement(clients, delta, e.UpdateCh)
}
