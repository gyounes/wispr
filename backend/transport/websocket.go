package transport

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	pb "github.com/gyounes/wispr/backend/proto"
	"github.com/gyounes/wispr/backend/server"
)

type WebSocketServer struct {
	upgrader    websocket.Upgrader
	connections *server.Connections
}

func NewWebSocketServer(connections *server.Connections) *WebSocketServer {
	return &WebSocketServer{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true }, // allow all origins for now
		},
		connections: connections,
	}
}

func (wss *WebSocketServer) HandleWS(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := wss.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	// First message (or query param) defines the username
	username := r.URL.Query().Get("username")
	if username == "" {
		log.Println("Missing username in query params")
		return
	}
	log.Printf("User %s connected via WebSocket", username)

	// Channel for messages destined for this user
	ch := make(chan *pb.Message, 10)
	wss.connections.Add(username, ch)
	defer wss.connections.Remove(username)

	// Send last 50 messages from DB
	if wss.connections.Storage != nil {
		lastMsgs, _ := wss.connections.Storage.GetLastMessages(username, 50)
		for _, m := range lastMsgs {
			ch <- &pb.Message{
				Sender:    m.Sender,
				Recipient: m.Recipient,
				Content:   m.Content,
				Timestamp: m.Timestamp.Format(time.RFC3339),
			}
		}
	}

	// Writer goroutine (send messages to WebSocket)
	go func() {
		for msg := range ch {
			data, _ := json.Marshal(msg)
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("Write error for %s: %v", username, err)
				return
			}
		}
	}()

	// Reader loop (receive messages from WebSocket)
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error for %s: %v", username, err)
			break
		}

		var msg pb.Message
		if err := json.Unmarshal(data, &msg); err != nil {
			log.Printf("Invalid message from %s: %v", username, err)
			continue
		}

		// Ensure sender matches connection username
		msg.Sender = username
		log.Printf("WS Message from %s to %s: %s", msg.Sender, msg.Recipient, msg.Content)

		// Broadcast through existing connection manager
		wss.connections.Broadcast(&msg)
	}
}
