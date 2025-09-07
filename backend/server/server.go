// backend/server/server.go
package server

import (
    "sync"
    "time"

    pb "github.com/gyounes/wispr/backend/proto"
)

type Server struct {
    pb.UnimplementedChatServiceServer
    mu       sync.Mutex
    clients  map[string]chan *pb.Message
}

func NewServer() *Server {
    return &Server{
        clients: make(map[string]chan *pb.Message),
    }
}

func (s *Server) SendMessage(msg *pb.Message) (*pb.Ack, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    // Broadcast to recipient if connected
    if ch, ok := s.clients[msg.Recipient]; ok {
        ch <- msg
    }

    return &pb.Ack{Success: true}, nil
}

func (s *Server) ReceiveMessages(msg *pb.Message, stream pb.ChatService_ReceiveMessagesServer) error {
    ch := make(chan *pb.Message, 10)

    s.mu.Lock()
    s.clients[msg.Sender] = ch
    s.mu.Unlock()

    // Clean up when function exits
    defer func() {
        s.mu.Lock()
        delete(s.clients, msg.Sender)
        s.mu.Unlock()
        close(ch)
    }()

    // Stream messages as they arrive
    for m := range ch {
        if err := stream.Send(m); err != nil {
            return err
        }
    }

    return nil
}
