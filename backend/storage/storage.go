package storage

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	DB *gorm.DB
}

type Message struct {
	ID        uint `gorm:"primaryKey"`
	Sender    string
	Recipient string
	Content   string
	Timestamp time.Time
}

// NewStorage initializes a Postgres-backed DB
func NewStorage(user, password, dbname, host string, port int) *Storage {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		host, user, password, dbname, port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// AutoMigrate creates the messages table if it doesn’t exist
	if err := db.AutoMigrate(&Message{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	return &Storage{DB: db}
}

func (s *Storage) SaveMessage(sender, recipient, content string, timestamp time.Time) error {
	return s.DB.Create(&Message{
		Sender:    sender,
		Recipient: recipient,
		Content:   content,
		Timestamp: timestamp,
	}).Error
}

func (s *Storage) GetLastMessages(user string, limit int) ([]Message, error) {
	var msgs []Message
	err := s.DB.Where("sender = ? OR recipient = ?", user, user).
		Order("timestamp desc").
		Limit(limit).
		Find(&msgs).Error
	return msgs, err
}
