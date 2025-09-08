package server

import (
	pb "github.com/gyounes/wispr/backend/proto"
	"time"
)

type Message = pb.Message
type Ack = pb.Ack

func NewMessage(sender, recipient, content string) *Message {
	return &Message{
		Sender:    sender,
		Recipient: recipient,
		Content:   content,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

func NewAck(success bool) *Ack {
	return &Ack{Success: success}
}
