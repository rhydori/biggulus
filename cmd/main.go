package main

import (
	"time"

	"github.com/rhydori/biggulus/pkg/engine"
	"github.com/rhydori/biggulus/pkg/server"
	"github.com/rhydori/biggulus/pkg/session"
)

const (
	tickInterval = 16 * time.Millisecond
	charSpeed    = 300.0
)

func main() {
	cs := session.NewClientStore()
	p := engine.NewPhysics(charSpeed)

	e := engine.NewEngine(tickInterval, cs, p)
	s := server.NewServer("[::]:8080", e, cs)

	s.StartServer()
	e.StartEngine()

	select {}
}
