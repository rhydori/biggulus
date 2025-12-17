package server

import (
	"fmt"

	"github.com/rhydori/biggulus/pkg/protocol"
	"github.com/rhydori/biggulus/pkg/session"
	"github.com/rhydori/logs"
)

func (s *Server) handleMsg(c *session.Client, msg string) {
	parsed := protocol.ParseMessage(msg)
	if parsed == nil {
		return
	}
	switch parsed.Entity {
	case "auth":
		s.handleAuth(c, parsed)
	case "character":
		c.Char.HandleCharacter(parsed)
	default:
		logs.Warnf("handleMsg: %s - Invalid Entity '%s'", c.Conn.RemoteAddr(), parsed.Entity)
		return
	}
}

func (s *Server) handleAuth(c *session.Client, msg *protocol.Message) {
	switch msg.Action {
	case "login":
		token, err := s.authService.Login(msg.Params)
		if err != nil {
			logs.Errorf("handleAuth - Login: %s: %s", c.Conn.RemoteAddr(), err.Error())

			msgbyte := protocol.CreateMessageBytes(msg.Entity, "error", err.Error())
			writeConn(c, msgbyte)
			return
		}
		msgbyte := protocol.CreateMessageBytes(msg.Entity, "logged", token.Value)
		writeConn(c, msgbyte)
	case "register":
		if err := s.authService.Register(msg.Params); err != nil {
			logs.Errorf("handleAuth - Register: %s: %s", c.Conn.RemoteAddr(), err.Error())

			msgbyte := protocol.CreateMessageBytes(msg.Entity, "error", err.Error())
			writeConn(c, msgbyte)
			return
		}
		msgbyte := protocol.CreateMessageBytes(msg.Entity, "registered", "User registered successfully.")
		writeConn(c, msgbyte)
	case "logout":
		if err := s.authService.Logout(msg.Params); err != nil {
			msgbyte := protocol.CreateMessageBytes(msg.Entity, "error", "Logout error.")
			writeConn(c, msgbyte)
		}
		msgbyte := protocol.CreateMessageBytes(msg.Entity, "logout", "User logout successfully.")
		writeConn(c, msgbyte)
	default:
		logs.Warnf("handleAuth: Invalid action %s", msg.Action)
		msgbyte := []byte(fmt.Sprintf("Invalid action %s", msg.Action))
		writeConn(c, msgbyte)
		return
	}
}
