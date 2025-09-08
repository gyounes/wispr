package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/gyounes/wispr/backend/proto"
	"github.com/gyounes/wispr/client/cli"
	"google.golang.org/grpc"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <username>")
		return
	}
	username := os.Args[1]

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewChatServiceClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Ctrl+C handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nGoodbye!")
		cancel()
		os.Exit(0)
	}()

	// Receive messages
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
			fmt.Printf("\n%s\n> ", cli.FormatIncoming(msg.Sender, msg.Content, msg.Timestamp))
		}
	}()

	// Read input
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')

		recipient, message, isCommand := cli.ParseInput(text)
		if isCommand {
			if cli.ExecuteCommand(message) {
				return
			}
			continue
		}
		if recipient == "" || message == "" {
			fmt.Println("Use format: recipient: message")
			continue
		}

		_, err := client.SendMessage(ctx, &pb.Message{
			Sender:    username,
			Recipient: recipient,
			Content:   message,
			Timestamp: cli.GetTimestamp(),
		})
		if err != nil {
			log.Printf("Send error: %v", err)
		} else {
			fmt.Println(cli.FormatOutgoing(recipient, message, cli.GetTimestamp()))
		}
	}
}
