package storage

import (
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupTestDB creates a fresh Postgres DB connection for tests
func SetupTestDB(t *testing.T) *Storage {
	t.Helper()

	dsn := "host=localhost user=postgres password=secret dbname=wispr_test port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect test db: %v", err)
	}

	// Drop and recreate table for a clean slate
	if err := db.Migrator().DropTable(&Message{}); err != nil {
		t.Fatalf("could not drop table: %v", err)
	}
	if err := db.AutoMigrate(&Message{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return &Storage{DB: db}
}
