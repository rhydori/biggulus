package server

import (
	"bufio"
	"errors"
	"net"

	"github.com/rhydori/biggulus/pkg/auth"
	"github.com/rhydori/biggulus/pkg/engine"
	"github.com/rhydori/biggulus/pkg/session"
	"github.com/rhydori/logs"
)

type Server struct {
	svrAddr     string
	engine      *engine.Engine
	clientStore *session.ClientStore
	authService *auth.AuthService

	bcastCh chan []byte
}

func NewServer(svrAddr string, engine *engine.Engine, cs *session.ClientStore, authService *auth.AuthService) *Server {
	return &Server{
		svrAddr:     svrAddr,
		engine:      engine,
		clientStore: cs,
		authService: authService,

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
		var netErr *net.OpError
		if errors.As(err, &netErr) {
			return
		}
		logs.Errorf("readConn: ScannerError from %s - '%v'", c.Conn.RemoteAddr(), err)
	}
}
