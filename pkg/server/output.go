package server

import (
	"github.com/rhydori/biggulus/pkg/session"
	"github.com/rhydori/logs"
)

func (s *Server) broadcast() {
	for msg := range s.bcastCh {
		clients := s.clientStore.ClientStoreSnapshot()

		for _, c := range clients {
			select {
			case c.OutCh <- append(msg, '\n'):
			default:
				logs.Warnf("broadcast: Client %s OutCh full, dropping messages", c.Conn.RemoteAddr())
			}
		}
	}
}

func (s *Server) forwardEngineUpdates() {
	for msg := range s.engine.UpdateCh {
		select {
		case s.bcastCh <- msg:
		default:
			logs.Warnf("fowardEngineUpdates: Broadcast buffer full, dropping")
		}
	}
}

func writeConn(c *session.Client, msg []byte) {
	logs.Warn(string(msg))
	c.OutCh <- []byte(append(msg, '\n'))
}
