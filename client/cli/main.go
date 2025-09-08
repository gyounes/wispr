package main

import (
    "bufio"
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "strings"
    "syscall"
    "time"

    "google.golang.org/grpc"
    pb "github.com/gyounes/wispr/backend/proto"
)

const (
    ColorReset  = "\033[0m"
    ColorGreen  = "\033[32m"
    ColorBlue   = "\033[34m"
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

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Handle Ctrl+C for graceful exit
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        fmt.Println("\nGoodbye!")
        cancel()
        os.Exit(0)
    }()

    // Start receiving messages
    go func() {
        stream, err := client.ReceiveMessages(ctx, &pb.Message{Sender: username})
        if err != nil {
            log.Fatalf("Failed to receive messages: %v", err)
        }

        for {
            msg, err := stream.Recv()
            if err != nil {
                return
            }
            fmt.Printf("\n%s[Incoming][%s] %s: %s%s\n> ", ColorGreen, msg.Timestamp, msg.Sender, msg.Content, ColorReset)
        }
    }()

    // Read input from user
    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Print("> ")
        text, _ := reader.ReadString('\n')
        text = strings.TrimSpace(text)
        if text == "" {
            continue
        }

        // Commands
        if strings.HasPrefix(text, "/") {
            switch text {
            case "/quit":
                fmt.Println("Goodbye!")
                return
            case "/list":
                fmt.Println("Connected users: Alice, Bob") // placeholder; server-driven later
                continue
            default:
                fmt.Println("Unknown command")
                continue
            }
        }

        // Send message: expect format recipient: message
        parts := strings.SplitN(text, ":", 2)
        if len(parts) != 2 {
            fmt.Println("Use format: recipient: message")
            continue
        }
        recipient := strings.TrimSpace(parts[0])
        content := strings.TrimSpace(parts[1])

        _, err := client.SendMessage(ctx, &pb.Message{
            Sender:    username,
            Recipient: recipient,
            Content:   content,
            Timestamp: time.Now().Format(time.RFC3339),
        })
        if err != nil {
            log.Printf("Send error: %v", err)
        } else {
            fmt.Printf("%s[Outgoing][%s] To %s: %s%s\n", ColorBlue, time.Now().Format(time.RFC3339), recipient, content, ColorReset)
        }
    }
}
