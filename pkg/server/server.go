package server

import (
	"bufio"
	"net"

	"github.com/google/uuid"
	"github.com/rhydori/biggulus/pkg/engine"
	"github.com/rhydori/biggulus/pkg/session"
	"github.com/rhydori/logs"
)

type Server struct {
	svrAddr string
	engine  *engine.Engine
	session *session.ClientStore

	bcastCh chan []byte
}

func NewServer(svrAddr string, engine *engine.Engine, session *session.ClientStore) *Server {
	return &Server{
		svrAddr: svrAddr,
		engine:  engine,
		session: session,

		bcastCh: make(chan []byte, 1024),
	}
}

func (s *Server) StartServer() {
	ln, err := net.Listen("tcp", s.svrAddr)
	if err != nil {
		logs.Fatalf("StartServer: %v", err)
	}
	logs.Infof("Server started at %s", ln.Addr())

	go func() {
		for msg := range s.engine.UpdateCh {
			select {
			case s.bcastCh <- msg:
			default:
				logs.Warnf("Server broadcast channel full, dropping update")
			}
		}
	}()
	go s.broadcast()
	go s.acceptConn(ln)
}

func (s *Server) acceptConn(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			logs.Errorf("acceptConn: %v", err)
			continue
		}
		c := &session.Client{
			ID:    uuid.NewString(),
			Conn:  conn,
			Input: &session.Input{},
		}
		s.session.AddClient(c)

		logs.Debugf("Connected: %s", c.Conn.RemoteAddr())
		s.writeConn(c, []byte("client_id|"+c.ID+"\n"))

		go s.readConn(c)
	}
}

func (s *Server) readConn(c *session.Client) {
	defer func() {
		c.Conn.Close()
		s.session.RemoveClient(c.ID)

		logs.Debugf("Disconnected: %s", c.Conn.RemoteAddr())
	}()
	scanner := bufio.NewScanner(c.Conn)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		msg := scanner.Bytes()
		s.handleMsg(c, msg)

		logs.Debugf("%s: %s", c.Conn.RemoteAddr(), msg)
	}
	if err := scanner.Err(); err != nil {
		logs.Errorf("Scanner error from %s: %v", c.Conn.RemoteAddr(), err)
	}
}

func (s *Server) writeConn(c *session.Client, msg []byte) {
	_, err := c.Conn.Write(msg)
	if err != nil {
		logs.Errorf("Write error to %s: %v", c.Conn.RemoteAddr(), err)
	}
}

func (s *Server) handleMsg(c *session.Client, msg []byte) {
	switch string(msg) {
	case "Left_PRESS":
		s.session.UpdateClientInput(c.ID, "Left", true)
	case "Left_RELEASE":
		s.session.UpdateClientInput(c.ID, "Left", false)
	case "Right_PRESS":
		s.session.UpdateClientInput(c.ID, "Right", true)
	case "Right_RELEASE":
		s.session.UpdateClientInput(c.ID, "Right", false)
	case "Up_PRESS":
		s.session.UpdateClientInput(c.ID, "Up", true)
	case "Up_RELEASE":
		s.session.UpdateClientInput(c.ID, "Up", false)
	case "Down_PRESS":
		s.session.UpdateClientInput(c.ID, "Down", true)
	case "Down_RELEASE":
		s.session.UpdateClientInput(c.ID, "Down", false)
	}
}

func (s *Server) broadcast() {
	go func() {
		for msg := range s.bcastCh {
			s.session.Mu.Lock()
			clients := make([]*session.Client, 0, len(s.session.Clients))
			for _, c := range s.session.Clients {
				clients = append(clients, c)
			}
			s.session.Mu.Unlock()

			for _, c := range clients {
				s.writeConn(c, append(msg, '\n'))
			}
		}
	}()
}
