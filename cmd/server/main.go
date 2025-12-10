package main

import (
	"github.com/rhydori/biggulus/pkg/engine"
	"github.com/rhydori/biggulus/pkg/server"
	"github.com/rhydori/biggulus/pkg/session"
)

func main() {
	cs := session.NewClientStore()
	e := engine.NewEngine(16, 300.0, cs)
	s := server.NewServer("[::]:8080", e, cs)
	s.StartServer()
	e.StartEngine()
	select {}
}
