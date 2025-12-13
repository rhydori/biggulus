package main

import (
	"time"

	"github.com/rhydori/biggulus/pkg/engine"
	"github.com/rhydori/biggulus/pkg/server"
	"github.com/rhydori/biggulus/pkg/session"
)

const (
	tickInterval   = 16 * time.Millisecond
	characterSpeed = 300.0
)

func main() {
	cs := session.NewClientStore()
	e := engine.NewEngine(tickInterval, characterSpeed, cs)
	s := server.NewServer("[::]:8080", e, cs)

	s.StartServer()
	e.StartEngine()

	select {}
}
