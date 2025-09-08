package utils

import (
	"testing"
)

func TestFormatTimestamp(t *testing.T) {
	ts := FormatTimestamp()
	if ts == "" {
		t.Fatal("Timestamp should not be empty")
	}
}

func TestColors(t *testing.T) {
	if ColorReset == "" || ColorGreen == "" || ColorBlue == "" {
		t.Fatal("Color constants should be defined")
	}
}
