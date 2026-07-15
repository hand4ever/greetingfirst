package model

import (
	"fmt"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestApplySchema(t *testing.T) {
	// open a fresh in-memory SQLite DB
	sqliteDB, err := initSoloDB(":memory:")
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	if err := ApplySchema(sqliteDB); err != nil {
		t.Fatalf("ApplySchema failed: %v", err)
	}

	// verify table exists
	if !sqliteDB.Migrator().HasTable("sl_users") {
		t.Error("expected sl_users table to exist after ApplySchema")
	}
	logOK(t, "ApplySchema PASS: sl_users table created")
}

func TestEnsureUserTable_AlreadyExists(t *testing.T) {
	sqliteDB, err := initSoloDB(":memory:")
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	// create table manually for testing (SQLite dialect uses sl_users)
	if err := sqliteDB.Exec(`
		CREATE TABLE sl_users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			phone VARCHAR(32) NOT NULL,
			name VARCHAR(64) NOT NULL,
			age INT DEFAULT 0,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		)
	`).Error; err != nil {
		t.Fatalf("failed to create test table: %v", err)
	}

	// should return immediately since table exists
	if err := EnsureUserTable(sqliteDB, "sqlite", 5*time.Second); err != nil {
		t.Fatalf("EnsureUserTable should succeed when table exists: %v", err)
	}
	logOK(t, "EnsureUserTable PASS: table exists, returned immediately")
}

func TestEnsureUserTable_Timeout(t *testing.T) {
	sqliteDB, err := initSoloDB(":memory:")
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	start := time.Now()
	err = EnsureUserTable(sqliteDB, "sqlite", 2*time.Second)
	elapsed := time.Since(start)

	if err == nil {
		t.Error("expected timeout error when table does not exist")
	}
	if elapsed < 2*time.Second {
		t.Errorf("expected at least 2s wait, got %v", elapsed)
	}
	if elapsed > 5*time.Second {
		t.Errorf("expected around 3s wait (poll interval), got %v", elapsed)
	}
	logOK(t, "EnsureUserTable timeout PASS: err=%v, elapsed=%v", err, elapsed)
}

// initSoloDB opens a standalone in-memory SQLite DB for isolated tests.
func initSoloDB(sqliteDSN string) (*gorm.DB, error) {
	sqliteDB, err := gorm.Open(sqlite.Open(sqliteDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("sqlite (%s): %w", sqliteDSN, err)
	}
	return sqliteDB, nil
}
