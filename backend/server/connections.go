package server

import (
	"sync"
	"time"

	pb "github.com/gyounes/wispr/backend/proto"
	"github.com/gyounes/wispr/backend/storage"
)

type Connections struct {
	mu      sync.Mutex
	clients map[string]chan *Message
	Storage *storage.Storage
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

func (c *Connections) Broadcast(msg *pb.Message) {
	// save to DB
	if c.Storage != nil {
		timestamp, _ := time.Parse(time.RFC3339, msg.Timestamp)
		_ = c.Storage.SaveMessage(msg.Sender, msg.Recipient, msg.Content, timestamp)
	}

	// send to recipient channel
	c.mu.Lock()
	ch, ok := c.clients[msg.Recipient]
	c.mu.Unlock()
	if ok {
		ch <- msg
	}
}
