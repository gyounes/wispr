package cli

import (
	"strings"
	"testing"
)

func TestParseInput(t *testing.T) {
	recipient, msg, isCmd := ParseInput("Bob: hello")
	if recipient != "Bob" || msg != "hello" || isCmd {
		t.Fatal("ParseInput failed for normal message")
	}

	_, cmd, isCmd := ParseInput("/quit")
	if cmd != "/quit" || !isCmd {
		t.Fatal("ParseInput failed for command")
	}

	recipient, msg, isCmd = ParseInput("invalidinput")
	if recipient != "" || msg != "" || isCmd {
		t.Fatal("ParseInput failed for invalid input")
	}
}

func TestFormatIncomingOutgoing(t *testing.T) {
	in := FormatIncoming("Alice", "Hello", "2025-09-08T12:00:00")
	out := FormatOutgoing("Bob", "Hi", "2025-09-08T12:00:00")

	if !strings.Contains(in, "Alice") || !strings.Contains(in, "Hello") {
		t.Fatal("FormatIncoming failed")
	}
	if !strings.Contains(out, "Bob") || !strings.Contains(out, "Hi") {
		t.Fatal("FormatOutgoing failed")
	}
}
