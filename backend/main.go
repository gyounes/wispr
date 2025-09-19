package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"google.golang.org/grpc"

	pb "github.com/gyounes/wispr/backend/proto"
	"github.com/gyounes/wispr/backend/server"
	"github.com/gyounes/wispr/backend/storage"
	"github.com/gyounes/wispr/backend/transport"
)

func main() {
	// Load DB settings from environment, fallback to defaults
	user := getEnv("DB_USER", "postgres")
	pass := getEnv("DB_PASS", "secret")
	name := getEnv("DB_NAME", "wispr_dev")
	host := getEnv("DB_HOST", "localhost")
	port := getEnvAsInt("DB_PORT", 5432)

	// Connect to Postgres
	store := storage.NewStorage(user, pass, name, host, port)

	// Shared connection manager with DB
	connections := server.NewConnections()
	connections.Storage = store

	// Start gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterChatServiceServer(grpcServer, &server.Server{Connections: connections})

	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("❌ failed to listen: %v", err)
		}
		log.Println("🚀 gRPC server running on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("❌ failed to serve gRPC: %v", err)
		}
	}()

	// Start WebSocket server
	wss := transport.NewWebSocketServer(connections)
	http.HandleFunc("/ws", wss.HandleWS)

	log.Println("🚀 WebSocket server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("❌ failed to serve WebSocket: %v", err)
	}
}

// Helpers for env variables
func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if valStr, ok := os.LookupEnv(key); ok {
		if val, err := strconv.Atoi(valStr); err == nil {
			return val
		}
	}
	return fallback
}
