package cli

import (
	"fmt"
	"time"
)

const (
	ColorReset = "\033[0m"
	ColorGreen = "\033[32m"
	ColorBlue  = "\033[34m"
)

// FormatIncoming prints an incoming message
func FormatIncoming(sender, content, timestamp string) string {
	return fmt.Sprintf("%s[Incoming][%s] %s: %s%s", ColorGreen, timestamp, sender, content, ColorReset)
}

// FormatOutgoing prints an outgoing message
func FormatOutgoing(recipient, content, timestamp string) string {
	return fmt.Sprintf("%s[Outgoing][%s] To %s: %s%s", ColorBlue, timestamp, recipient, content, ColorReset)
}

// GetTimestamp returns formatted current timestamp
func GetTimestamp() string {
	return time.Now().Format(time.RFC3339)
}
