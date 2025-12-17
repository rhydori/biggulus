package main

import (
	"time"

	"github.com/rhydori/biggulus/pkg/auth"
	"github.com/rhydori/biggulus/pkg/database"
	"github.com/rhydori/biggulus/pkg/engine"
	"github.com/rhydori/biggulus/pkg/server"
	"github.com/rhydori/biggulus/pkg/session"
	"github.com/rhydori/logs"
)

const (
	tickInterval = 16 * time.Millisecond
	charSpeed    = 300.0
)

func main() {
	db := database.OpenSQLite("./sqlite/game.db")
	if db == nil {
		logs.Fatal("Database is nil.")
	}
	userRep := auth.NewSQLiteUserRepo(db)
	tokenRep := auth.NewSQLiteTokenRepo(db)
	authService := auth.NewService(userRep, tokenRep)

	clientStore := session.NewClientStore()

	physics := engine.NewPhysics(charSpeed)
	engine := engine.NewEngine(tickInterval, clientStore, physics)

	server := server.NewServer("[::]:8080", engine, clientStore, authService)

	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for range ticker.C {
			tokenRep.DeleteExpired()
		}
	}()

	server.StartServer()
	engine.StartEngine()
}
