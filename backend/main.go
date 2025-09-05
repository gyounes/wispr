
package main

import (
    "log"
    "net"

    "google.golang.org/grpc"
    "github.com/gyounes/wispr/backend/proto"
)

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    s := grpc.NewServer()
    proto.RegisterChatServiceServer(s, &Server{})

    log.Println("gRPC server running on :50051")
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
