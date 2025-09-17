package main

import (
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"

	pb "github.com/gyounes/wispr/backend/proto"
	"github.com/gyounes/wispr/backend/server"
	"github.com/gyounes/wispr/backend/storage"
	"github.com/gyounes/wispr/backend/transport"
)

func main() {
	// Connect to Postgres (running in Docker)
	store := storage.NewStorage("postgres", "secret", "wispr_dev", "localhost", 5432)

	// Shared connection manager with DB
	connections := server.NewConnections()
	connections.Storage = store

	// Start gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterChatServiceServer(grpcServer, &server.Server{Connections: connections})

	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		log.Println("gRPC server running on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// Start WebSocket server
	wss := transport.NewWebSocketServer(connections)
	http.HandleFunc("/ws", wss.HandleWS)

	log.Println("WebSocket server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("failed to serve WebSocket: %v", err)
	}
}
