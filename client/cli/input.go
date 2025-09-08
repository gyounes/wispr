package cli

import "strings"

// ParseInput parses raw input into recipient & message
func ParseInput(input string) (recipient, message string, isCommand bool) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", "", false
	}

	if strings.HasPrefix(input, "/") {
		return "", input, true
	}

	parts := strings.SplitN(input, ":", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), false
}
