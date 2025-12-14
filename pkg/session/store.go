package session

import (
	"sync"
)

type ClientStore struct {
	Clients map[string]*Client
	Mu      sync.Mutex
}

func NewClientStore() *ClientStore {
	return &ClientStore{
		Clients: make(map[string]*Client),
	}
}

func (cs *ClientStore) AddClientToStore(c *Client) {
	cs.Mu.Lock()
	defer cs.Mu.Unlock()
	cs.Clients[c.ID] = c
}

func (cs *ClientStore) RemoveClientFromStore(id string) {
	cs.Mu.Lock()
	defer cs.Mu.Unlock()
	delete(cs.Clients, id)
}

func (cs *ClientStore) ClientStoreSnapshot() []*Client {
	cs.Mu.Lock()
	defer cs.Mu.Unlock()

	out := make([]*Client, 0, len(cs.Clients))
	for _, c := range cs.Clients {
		out = append(out, c)
	}
	return out
}
