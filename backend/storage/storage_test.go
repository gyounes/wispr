package storage

import (
	"os"
	"testing"
	"time"
)

// Setup test DB before running tests
func init() {
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASS", "secret")
	os.Setenv("DB_NAME", "wispr_test")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
}

func setupTestStorage(t *testing.T) *Storage {
	store := NewStorage("postgres", "secret", "wispr_test", "localhost", 5432)

	// Clean the table before each test
	if err := store.DB.Exec("TRUNCATE messages RESTART IDENTITY").Error; err != nil {
		t.Fatalf("failed to truncate messages: %v", err)
	}
	return store
}

func TestSaveMessage(t *testing.T) {
	store := setupTestStorage(t)

	msgTime := time.Now().UTC()
	err := store.SaveMessage("Alice", "Bob", "Hello Bob!", msgTime)
	if err != nil {
		t.Fatalf("failed to save message: %v", err)
	}

	var count int64
	if err := store.DB.Model(&Message{}).Count(&count).Error; err != nil {
		t.Fatalf("failed to count messages: %v", err)
	}

	if count != 1 {
		t.Fatalf("expected 1 message, got %d", count)
	}
}

func TestGetLastMessages(t *testing.T) {
	store := setupTestStorage(t)

	now := time.Now().UTC()
	msgs := []Message{
		{Sender: "Alice", Recipient: "Bob", Content: "Msg1", Timestamp: now.Add(-2 * time.Minute)},
		{Sender: "Bob", Recipient: "Alice", Content: "Msg2", Timestamp: now.Add(-1 * time.Minute)},
		{Sender: "Alice", Recipient: "Charlie", Content: "Msg3", Timestamp: now},
	}

	for _, m := range msgs {
		if err := store.SaveMessage(m.Sender, m.Recipient, m.Content, m.Timestamp); err != nil {
			t.Fatalf("failed to save message: %v", err)
		}
	}

	results, err := store.GetLastMessages("Alice", 2)
	if err != nil {
		t.Fatalf("GetLastMessages failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(results))
	}

	// Ensure messages are ordered by timestamp desc
	if results[0].Content != "Msg3" || results[1].Content != "Msg2" {
		t.Fatalf("unexpected message order: %+v", results)
	}
}
