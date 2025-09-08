package server

import (
	"context"
	"testing"
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

func TestSendAndReceiveMessage(t *testing.T) {
	s := NewServer()
	ch := make(chan *Message, 1)
	s.connections.Add("Bob", ch)

	msg := NewMessage("Alice", "Bob", "Hello Bob!")
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
