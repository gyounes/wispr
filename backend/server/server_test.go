package server

import (
    "context"
    "testing"
    "time"

    pb "github.com/gyounes/wispr/backend/proto"
)

func TestSendMessage(t *testing.T) {
    s := NewServer()

    // Create fake client channel
    s.clients["Bob"] = make(chan *pb.Message, 1)

    msg := &pb.Message{
        Sender:    "Alice",
        Recipient: "Bob",
        Content:   "Hello Bob!",
        Timestamp: time.Now().Format(time.RFC3339),
    }

    ack, err := s.SendMessage(context.Background(), msg) // ✅ add context
    if err != nil {
        t.Fatalf("SendMessage error: %v", err)
    }
    if !ack.Success {
        t.Fatalf("SendMessage failed")
    }

    // Check if Bob received the message
    select {
    case m := <-s.clients["Bob"]:
        if m.Content != "Hello Bob!" {
            t.Fatalf("Expected 'Hello Bob!', got '%s'", m.Content)
        }
    case <-time.After(time.Second):
        t.Fatalf("Message not received by Bob")
    }
}

func TestReceiveMessages(t *testing.T) {
    s := NewServer()

    // fake stream
    msgs := make(chan *pb.Message, 1)
    s.clients["Alice"] = msgs

    msg := &pb.Message{
        Sender:    "Bob",
        Recipient: "Alice",
        Content:   "Hi Alice!",
        Timestamp: time.Now().Format(time.RFC3339),
    }

    _, err := s.SendMessage(context.Background(), msg) // ✅ add context
    if err != nil {
        t.Fatalf("SendMessage error: %v", err)
    }

    select {
    case m := <-msgs:
        if m.Content != "Hi Alice!" {
            t.Fatalf("Expected 'Hi Alice!', got '%s'", m.Content)
        }
    case <-time.After(time.Second):
        t.Fatalf("Message not received")
    }
}
