package server

import (
	"context"
	"testing"

	"github.com/gyounes/wispr/backend/storage"
)

func TestConnections(t *testing.T) {
	conn := NewConnections()
	ch := make(chan *Message, 1)
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
	// in-memory DB
	db := storage.New(":memory:")

	s := NewServer()
	s.Connections.Storage = db
	s.Storage = db

	ch := make(chan *Message, 1)
	s.Connections.Add("Bob", ch)

	msg := NewMessage("Alice", "Bob", "Hello Bob!")
	ack, err := s.SendMessage(context.Background(), msg)
	if err != nil {
		t.Fatal(err)
	}
	if !ack.Success {
		t.Fatal("SendMessage failed")
	}

	// Check channel delivery
	select {
	case m := <-ch:
		if m.Content != "Hello Bob!" {
			t.Fatalf("Expected 'Hello Bob!', got '%s'", m.Content)
		}
	default:
		t.Fatal("Message not received by Bob")
	}

	// Check DB persistence
	msgs, err := db.GetLastMessages("Bob", 10)
	if err != nil {
		t.Fatalf("DB GetLastMessages failed: %v", err)
	}
	if len(msgs) == 0 {
		t.Fatal("No messages saved in DB")
	}
	if msgs[0].Content != "Hello Bob!" || msgs[0].Sender != "Alice" {
		t.Fatalf("Unexpected DB message: %+v", msgs[0])
	}
}
