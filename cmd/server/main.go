package main

import (
	"bufio"
	"fmt"
	"math"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rhydori/logs"
)

type Server struct {
	svrAddr string
	clients map[string]*Client
	mu      sync.Mutex
	bcastCh chan []byte
	engine  *Engine
}

type Client struct {
	id    string
	conn  net.Conn
	input *Input
	X, Y  float64

	lastX, lastY float64
}

type Engine struct {
	tick         time.Duration
	speedperTick float64
	inputCh      chan *Input
}

type Input struct {
	left  bool
	right bool
	up    bool
	down  bool
}

type Vec2 struct {
	X, Y float64
}

func (v Vec2) Normalize() Vec2 {
	l := math.Hypot(v.X, v.Y)
	if l == 0 {
		return Vec2{0, 0}
	}
	return Vec2{v.X / l, v.Y / l}
}

func main() {
	eng := newEngine(1, 300.0)
	svr := newServer("[::]:8080", eng)
	svr.startServer(svr.svrAddr)
	eng.startEngine(svr)
	select {}
}

func newServer(svrAddr string, eng *Engine) *Server {
	return &Server{
		svrAddr: svrAddr,
		clients: make(map[string]*Client),
		bcastCh: make(chan []byte, 1024),
		engine:  eng,
	}
}

func newEngine(tickMs int8, speed float64) *Engine {
	return &Engine{
		tick:         time.Duration(tickMs) * time.Millisecond,
		inputCh:      make(chan *Input, 1024),
		speedperTick: speed * float64(tickMs) / 1000.0,
	}
}

func (s *Server) startServer(svrAddr string) {
	ln, err := net.Listen("tcp", svrAddr)
	if err != nil {
		logs.Fatal(err)
	}
	logs.Infof("Server started at %s", ln.Addr())

	go s.broadcast()
	go s.acceptConn(ln)
}

func (e *Engine) startEngine(s *Server) {
	logs.Info("Engine started at ", e.tick)
	ticker := time.NewTicker(e.tick)
	go func() {
		for range ticker.C {
			e.updateClients(s)
		}
	}()
}

func (s *Server) acceptConn(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			logs.Error(err)
			continue
		}
		c := &Client{id: uuid.NewString(), conn: conn, input: &Input{}}
		s.mu.Lock()
		s.clients[c.id] = c
		s.mu.Unlock()
		logs.Debugf("Connected: %s", c.conn.RemoteAddr())

		_, err = c.conn.Write([]byte("client_id|" + c.id + "\n"))
		if err != nil {

		}
		go s.readConn(c)
	}
}

func (s *Server) readConn(c *Client) {
	defer func() {
		c.conn.Close()
		s.mu.Lock()
		delete(s.clients, c.id)
		s.mu.Unlock()
		logs.Debugf("Disconnected: %s", c.conn.RemoteAddr())
	}()

	scanner := bufio.NewScanner(c.conn)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		msg := scanner.Text()
		logs.Debugf("%s: %s", c.conn.RemoteAddr(), msg)
		s.handleMsg(c, []byte(msg))
	}
	if err := scanner.Err(); err != nil {
		logs.Error(err)
	}
}

func (s *Server) handleMsg(c *Client, msg []byte) {
	switch string(msg) {
	case "Left_PRESS":
		s.mu.Lock()
		c.input.left = true
		s.mu.Unlock()
		s.engine.inputCh <- c.input
	case "Left_RELEASE":
		s.mu.Lock()
		c.input.left = false
		s.mu.Unlock()
		s.engine.inputCh <- c.input
	case "Right_PRESS":
		s.mu.Lock()
		c.input.right = true
		s.mu.Unlock()
		s.engine.inputCh <- c.input
	case "Right_RELEASE":
		s.mu.Lock()
		c.input.right = false
		s.mu.Unlock()
		s.engine.inputCh <- c.input
	case "Up_PRESS":
		s.mu.Lock()
		c.input.up = true
		s.mu.Unlock()
		s.engine.inputCh <- c.input
	case "Up_RELEASE":
		s.mu.Lock()
		c.input.up = false
		s.mu.Unlock()
		s.engine.inputCh <- c.input
	case "Down_PRESS":
		s.mu.Lock()
		c.input.down = true
		s.mu.Unlock()
		s.engine.inputCh <- c.input
	case "Down_RELEASE":
		s.mu.Lock()
		c.input.down = false
		s.mu.Unlock()
		s.engine.inputCh <- c.input
	}
}

func (s *Server) broadcast() {
	go func() {
		for msg := range s.bcastCh {
			s.mu.Lock()
			clients := make([]*Client, 0, len(s.clients))
			for _, c := range s.clients {
				clients = append(clients, c)
			}
			s.mu.Unlock()

			for _, c := range clients {
				out := make([]byte, len(msg)+1)
				copy(out, msg)
				out[len(msg)] = '\n'
				_, err := c.conn.Write(out)
				if err != nil {
					logs.Errorf("Write error to %s: %v", c.conn.RemoteAddr(), err)
				}
			}
		}
	}()
}

func (e *Engine) updateClients(s *Server) {
	s.mu.Lock()
	clients := make([]*Client, 0, len(s.clients))
	for _, c := range s.clients {
		clients = append(clients, c)
	}
	s.mu.Unlock()

	for _, c := range clients {
		dir := Vec2{}
		s.mu.Lock()
		if c.input.left {
			dir.X -= 1
		}
		if c.input.right {
			dir.X += 1
		}
		if c.input.up {
			dir.Y -= 1
		}
		if c.input.down {
			dir.Y += 1
		}
		s.mu.Unlock()

		dir = dir.Normalize()
		c.X += dir.X * e.speedperTick
		c.Y += dir.Y * e.speedperTick

		if c.X != c.lastX || c.Y != c.lastY {
			msg := []byte(fmt.Sprintf("move|%s|%f|%f", c.id, c.X, c.Y))
			select {
			case s.bcastCh <- msg:
			default:
				// se estiver cheio, drop para não travar o servidor
				logs.Warnf("bcast channel full, dropping update for %s", c.id)
			}

			c.lastX = c.X
			c.lastY = c.Y
		}
	}

}
