package server

import (
	"bufio"
	"net"
	"strings"

	"github.com/rhydori/biggulus/pkg/engine"
	"github.com/rhydori/biggulus/pkg/session"
	"github.com/rhydori/logs"
)

type Server struct {
	svrAddr     string
	engine      *engine.Engine
	clientStore *session.ClientStore

	bcastCh chan []byte
}

func NewServer(svrAddr string, engine *engine.Engine, cs *session.ClientStore) *Server {
	return &Server{
		svrAddr:     svrAddr,
		engine:      engine,
		clientStore: cs,

		bcastCh: make(chan []byte, 256),
	}
}

func (s *Server) StartServer() {
	ln, err := net.Listen("tcp", s.svrAddr)
	if err != nil {
		logs.Fatalf("StartServer: %v", err)
	}
	logs.Infof("Server started at %s", ln.Addr())

	go s.acceptConn(ln)
	go s.broadcast()
	go s.forwardEngineUpdates()
}

func (s *Server) acceptConn(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			logs.Errorf("acceptConn: %v", err)
			continue
		}
		logs.Debugf("Connected: %s", conn.RemoteAddr())

		c := session.NewClient(conn)
		s.clientStore.AddClientToStore(c)

		c.OutCh <- []byte("client_id|" + c.ID + "\n")

		go s.readConn(c)
	}
}

func (s *Server) readConn(c *session.Client) {
	defer func() {
		c.Conn.Close()
		close(c.OutCh)
		s.clientStore.RemoveClientFromStore(c.ID)

		logs.Debugf("Disconnected: %s", c.Conn.RemoteAddr())
	}()
	scanner := bufio.NewScanner(c.Conn)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		msg := scanner.Text()
		s.handleMsg(c, msg)
	}
	if err := scanner.Err(); err != nil {
		logs.Errorf("Scanner error from %s: %v", c.Conn.RemoteAddr(), err)
	}
}

func (s *Server) writeConn(c *session.Client, msg string) {
	_, err := c.Conn.Write([]byte(msg))
	if err != nil {
		logs.Errorf("Write error to %s: %v", c.Conn.RemoteAddr(), err)
	}
}

func (s *Server) handleMsg(c *session.Client, msg string) {
	// parts example: entity|action|obj|state
	parts := strings.Split(msg, "|")
	if parts[0] != "character" {
		return
	}
	entity := parts[0]

	logs.Debug(parts)
	switch entity {
	case "character":
		c.Char.HandleCharacter(parts)
	case "inventory":
	default:
		logs.Warnf("handleMsg: %s - Entity '%s' not found", c.Conn.RemoteAddr(), entity)
	}
}

func (s *Server) broadcast() {
	for msg := range s.bcastCh {
		s.clientStore.Mu.Lock()
		clients := make([]*session.Client, 0, len(s.clientStore.Clients))
		for _, c := range s.clientStore.Clients {
			clients = append(clients, c)
		}
		s.clientStore.Mu.Unlock()

		for _, c := range clients {
			select {
			case c.OutCh <- append(msg, '\n'):
			default:
				logs.Warnf("Client %s OutCh full, dropping messages", c.Conn)
			}
		}
	}
}

func (s *Server) forwardEngineUpdates() {
	for msg := range s.engine.UpdateCh {
		select {
		case s.bcastCh <- msg:
		default:
			logs.Warnf("Broadcast buffer full, dropping")
		}
	}
}
