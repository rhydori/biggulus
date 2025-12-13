package session

import (
	"net"
	"sync"

	"github.com/rhydori/logs"
)

type ClientStore struct {
	Clients map[string]*Client
	Mu      sync.Mutex
}

type Client struct {
	ID   string
	Conn net.Conn

	Input *Input
	X, Y  float64

	LastX, LastY float64
}

type Input struct {
	Left  bool
	Right bool
	Up    bool
	Down  bool
}

type InputSetter func(input *Input, state bool)

var inputActions = map[string]InputSetter{
	"Left": func(i *Input, state bool) {
		i.Left = state
	},
	"Right": func(i *Input, state bool) {
		i.Right = state
	},
	"Up": func(i *Input, state bool) {
		i.Up = state
	},
	"Down": func(i *Input, state bool) {
		i.Down = state
	},
}

func NewClientStore() *ClientStore {
	return &ClientStore{
		Clients: make(map[string]*Client),
	}
}

func (cs *ClientStore) AddClient(c *Client) {
	cs.Mu.Lock()
	defer cs.Mu.Unlock()
	cs.Clients[c.ID] = c
}

func (cs *ClientStore) RemoveClient(id string) {
	cs.Mu.Lock()
	defer cs.Mu.Unlock()
	delete(cs.Clients, id)
}

func (cs *ClientStore) GetClient(id string) *Client {
	cs.Mu.Lock()
	defer cs.Mu.Unlock()
	client, ok := cs.Clients[id]
	if !ok {
		logs.Errorf("ChangeInput: Client not found in ClientStore %s", id)
	}
	return client

}

func (cs *ClientStore) UpdateClient(id string) {
	cs.Mu.Lock()
	defer cs.Mu.Unlock()

}

func (cs *ClientStore) UpdateClientInput(id string, key string, state bool) {
	client := cs.GetClient(id)
	if client == nil {
		return
	}

	cs.Mu.Lock()
	defer cs.Mu.Unlock()

	setter, found := inputActions[key]
	if !found {
		logs.Errorf("UpdateClientInput: Unknown input key %s for client %s", key, id)
		return
	}

	setter(client.Input, state)
}
