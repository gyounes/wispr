package server

import (
	"context"
	"os"
	"testing"
	"time"

	pb "github.com/gyounes/wispr/backend/proto"
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

func TestConnections(t *testing.T) {
	conn := NewConnections()
	ch := make(chan *pb.Message, 1)
	conn.Add("Alice", ch)

	if _, ok := conn.Get("Alice"); !ok {
		t.Fatal("Alice should exist")
	}

	conn.Remove("Alice")
	if _, ok := conn.Get("Alice"); ok {
		t.Fatal("Alice should have been removed")
	}
}

func TestSendAndReceiveMessageWithDB(t *testing.T) {
	// Connect to test DB
	store := storage.NewStorage("postgres", "secret", "wispr_test", "localhost", 5432)
	connections := NewConnections()
	connections.Storage = store

	s := &Server{Connections: connections}

	ch := make(chan *pb.Message, 1)
	s.Connections.Add("Bob", ch)

	msg := &pb.Message{
		Sender:    "Alice",
		Recipient: "Bob",
		Content:   "Hello Bob!",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	ack, err := s.SendMessage(context.Background(), msg)
	if err != nil {
		t.Fatal(err)
	}
	if !ack.Success {
		t.Fatal("SendMessage failed")
	}

	select {
	case m := <-ch:
		if m.Content != "Hello Bob!" {
			t.Fatalf("Expected 'Hello Bob!', got '%s'", m.Content)
		}
	default:
		t.Fatal("Message not received by Bob")
	}
}
