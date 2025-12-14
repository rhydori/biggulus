package session

import (
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/rhydori/logs"
)

type Client struct {
	ID   string
	Conn net.Conn

	Char *Character

	OutCh chan []byte
}

func NewClient(conn net.Conn) *Client {
	c := &Client{
		ID:   uuid.NewString(),
		Conn: conn,
		Char: NewCharacter(),

		OutCh: make(chan []byte, 256),
	}

	go c.writePump()
	return c
}

func (c *Client) writePump() {
	for msg := range c.OutCh {
		c.Conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
		if _, err := c.Conn.Write(msg); err != nil {
			logs.Errorf("writePump error: %s: %v", c.Conn.RemoteAddr(), err)
			return
		}
	}
}
