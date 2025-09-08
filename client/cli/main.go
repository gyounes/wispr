package main

import (
    "bufio"
    "context"
    "fmt"
    "log"
    "os"
    "strings"
    "time"

    "google.golang.org/grpc"
    pb "github.com/gyounes/wispr/backend/proto"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run main.go <username>")
        return
    }
    username := os.Args[1]

    // Connect to gRPC server
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()

    client := pb.NewChatServiceClient(conn)

    // Start receiving messages in a goroutine
    go func() {
        stream, err := client.ReceiveMessages(context.Background(), &pb.Message{Sender: username})
        if err != nil {
            log.Fatalf("Failed to receive messages: %v", err)
        }

        for {
            msg, err := stream.Recv()
            if err != nil {
                log.Printf("Receive error: %v", err)
                return
            }
            fmt.Printf("\n[%s] %s: %s\n> ", msg.Timestamp, msg.Sender, msg.Content)
        }
    }()

    // Read input from user and send messages
    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Print("> ")
        text, _ := reader.ReadString('\n')
        text = strings.TrimSpace(text)

        if text == "" {
            continue
        }

        // Expect format: recipient: message
        parts := strings.SplitN(text, ":", 2)
        if len(parts) != 2 {
            fmt.Println("Use format: recipient: message")
            continue
        }
        recipient := strings.TrimSpace(parts[0])
        content := strings.TrimSpace(parts[1])

        _, err := client.SendMessage(context.Background(), &pb.Message{
            Sender:    username,
            Recipient: recipient,
            Content:   content,
            Timestamp: time.Now().Format(time.RFC3339),
        })
        if err != nil {
            log.Printf("Send error: %v", err)
        }
    }
}
