package server

import (
	"context"

	pb "github.com/gyounes/wispr/backend/proto"
)

type Server struct {
	pb.UnimplementedChatServiceServer
	connections *Connections
}

func NewServer() *Server {
	return &Server{
		connections: NewConnections(),
	}
}

// SendMessage sends message to recipient
func (s *Server) SendMessage(ctx context.Context, msg *pb.Message) (*pb.Ack, error) {
	s.connections.Broadcast(msg)
	return NewAck(true), nil
}

// ReceiveMessages streams messages for a client
func (s *Server) ReceiveMessages(msg *pb.Message, stream pb.ChatService_ReceiveMessagesServer) error {
	ch := make(chan *pb.Message, 10)
	s.connections.Add(msg.Sender, ch)
	defer func() {
		s.connections.Remove(msg.Sender)
		close(ch)
	}()

	for m := range ch {
		if err := stream.Send(m); err != nil {
			return err
		}
	}
	return nil
}
