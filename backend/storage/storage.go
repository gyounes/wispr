package storage

import (
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Message represents a chat message in the DB
type Message struct {
	ID        uint `gorm:"primaryKey"`
	Sender    string
	Recipient string
	Content   string
	Timestamp time.Time
}

// Storage wraps the DB connection
type Storage struct {
	DB *gorm.DB
}

// New initializes the DB (SQLite)
func New(dbPath string) *Storage {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	// Migrate the schema
	if err := db.AutoMigrate(&Message{}); err != nil {
		log.Fatalf("failed to migrate DB: %v", err)
	}

	return &Storage{DB: db}
}

// SaveMessage persists a message
func (s *Storage) SaveMessage(sender, recipient, content string, timestamp time.Time) error {
	msg := &Message{
		Sender:    sender,
		Recipient: recipient,
		Content:   content,
		Timestamp: timestamp,
	}
	return s.DB.Create(msg).Error
}

// GetLastMessages retrieves last N messages (for a user)
func (s *Storage) GetLastMessages(username string, limit int) ([]Message, error) {
	var messages []Message
	err := s.DB.
		Where("sender = ? OR recipient = ?", username, username).
		Order("timestamp desc").
		Limit(limit).
		Find(&messages).Error
	if err != nil {
		return nil, err
	}

	// reverse so oldest first
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}
