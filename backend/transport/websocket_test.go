package transport

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	pb "github.com/gyounes/wispr/backend/proto"
	"github.com/gyounes/wispr/backend/server"
	"github.com/gyounes/wispr/backend/storage"
)

// Setup test DB before running tests
func init() {
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASS", "secret")
	os.Setenv("DB_NAME", "wispr_test")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
}

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
	store := storage.NewStorage("postgres", "secret", "wispr_test", "localhost", 5432)
	connections := server.NewConnections()
	connections.Storage = store
	wss := NewWebSocketServer(connections)

	// Start test HTTP server
	ts := httptest.NewServer(http.HandlerFunc(wss.HandleWS))
	defer ts.Close()

	// Connect two clients
	alice := dialWS(ts.URL, "Alice", t)
	defer alice.Close()
	bob := dialWS(ts.URL, "Bob", t)
	defer bob.Close()

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
		t.Fatalf("Unexpected message received: Sender=%s, Recipient=%s, Content=%s, Timestamp=%s", received.Sender, received.Recipient, received.Content, received.Timestamp)
	}
}
