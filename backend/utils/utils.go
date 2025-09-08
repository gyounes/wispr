package utils

import "time"

const (
	ColorReset = "\033[0m"
	ColorGreen = "\033[32m"
	ColorBlue  = "\033[34m"
)

func FormatTimestamp() string {
	return time.Now().Format(time.RFC3339)
}
