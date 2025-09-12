package transport

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	pb "github.com/gyounes/wispr/backend/proto"
	"github.com/gyounes/wispr/backend/server"
	"github.com/gyounes/wispr/backend/storage"
)

// Helper: dial test WebSocket server
func dialWS(tsURL, username string, t *testing.T) *websocket.Conn {
	wsURL := "ws" + tsURL[4:] + "/ws?username=" + username
	c, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Dial failed: %v", err)
	}
	return c
}

func TestWebSocketBroadcastWithDB(t *testing.T) {
	// in-memory DB
	db := storage.New(":memory:")

	connections := server.NewConnections()
	connections.Storage = db
	wss := NewWebSocketServer(connections)

	// Start test HTTP server
	ts := httptest.NewServer(http.HandlerFunc(wss.HandleWS))
	defer ts.Close()

	// Connect two clients
	alice := dialWS(ts.URL, "Alice", t)
	defer alice.Close()
	bob := dialWS(ts.URL, "Bob", t)
	defer bob.Close()

	// Give channels a moment to register
	time.Sleep(50 * time.Millisecond)

	// Alice sends message to Bob
	msg := &pb.Message{
		Sender:    "Alice",
		Recipient: "Bob",
		Content:   "Hello Bob!",
		Timestamp: time.Now().Format(time.RFC3339),
	}
	data, _ := json.Marshal(msg)
	if err := alice.WriteMessage(websocket.TextMessage, data); err != nil {
		t.Fatalf("Alice send failed: %v", err)
	}

	// Bob should receive it
	_, bobData, err := bob.ReadMessage()
	if err != nil {
		t.Fatalf("Bob read failed: %v", err)
	}
	var received pb.Message
	if err := json.Unmarshal(bobData, &received); err != nil {
		t.Fatalf("Bob unmarshal failed: %v", err)
	}

	if received.Content != "Hello Bob!" || received.Sender != "Alice" {
		t.Fatalf("Unexpected message received: %+v", received)
	}

	// Check DB persistence
	msgs, err := db.GetLastMessages("Bob", 10)
	if err != nil {
		t.Fatalf("DB GetLastMessages failed: %v", err)
	}
	if len(msgs) == 0 || msgs[0].Content != "Hello Bob!" {
		t.Fatal("Message not saved in DB")
	}
}
