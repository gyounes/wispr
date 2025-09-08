package server

import "sync"

type Connections struct {
	mu      sync.Mutex
	clients map[string]chan *Message
}

func NewConnections() *Connections {
	return &Connections{
		clients: make(map[string]chan *Message),
	}
}

func (c *Connections) Add(username string, ch chan *Message) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.clients[username] = ch
}

func (c *Connections) Remove(username string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.clients, username)
}

func (c *Connections) Get(username string) (chan *Message, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	ch, ok := c.clients[username]
	return ch, ok
}

func (c *Connections) List() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	users := make([]string, 0, len(c.clients))
	for k := range c.clients {
		users = append(users, k)
	}
	return users
}

func (c *Connections) Broadcast(msg *Message) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ch, ok := c.clients[msg.Recipient]; ok {
		ch <- msg
	}
}
