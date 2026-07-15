package model

import (
	_ "embed"
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//go:embed schema.sql
var sqliteSchemaSQL string

// DB is the global MySQL database instance.
var DB *gorm.DB

// SQLiteDB is the global SQLite database instance.
var SQLiteDB *gorm.DB

// schemaPollInterval is the interval between table existence checks.
const schemaPollInterval = 3 * time.Second

// InitDB opens both MySQL and SQLite connections.
// It uses fail-fast: either connection failure causes an immediate error.
func InitDB(mysqlDSN, sqliteDSN string) error {
	// open MySQL
	mysqlDB, err := gorm.Open(mysql.Open(mysqlDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("mysql (%s): %w", mysqlDSN, err)
	}
	sqlDB, err := mysqlDB.DB()
	if err != nil {
		return fmt.Errorf("mysql get underlying db: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("mysql ping (%s): %w", mysqlDSN, err)
	}
	DB = mysqlDB

	// open SQLite
	sqliteDB, err := gorm.Open(sqlite.Open(sqliteDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("sqlite (%s): %w", sqliteDSN, err)
	}
	sqlDB2, err := sqliteDB.DB()
	if err != nil {
		return fmt.Errorf("sqlite get underlying db: %w", err)
	}
	if err := sqlDB2.Ping(); err != nil {
		return fmt.Errorf("sqlite ping (%s): %w", sqliteDSN, err)
	}
	SQLiteDB = sqliteDB

	return nil
}

// ApplySchema executes the embedded SQLite schema SQL against the given DB.
// Used only by tests to initialize in-memory SQLite databases.
func ApplySchema(db *gorm.DB) error {
	// strip comment lines first, then join and split by statement
	lines := strings.Split(sqliteSchemaSQL, "\n")
	var cleanLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "--") {
			continue
		}
		cleanLines = append(cleanLines, line)
	}
	cleanSQL := strings.Join(cleanLines, "\n")

	statements := strings.Split(cleanSQL, ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if err := db.Exec(stmt).Error; err != nil {
			return fmt.Errorf("exec schema statement: %w\nstatement: %s", err, stmt)
		}
	}
	return nil
}

// EnsureUserTable checks if the users table exists, blocks and polls until it does.
// maxWait <= 0 means wait indefinitely.
// It MUST NOT auto-create the table.
func EnsureUserTable(db *gorm.DB, dialect string, maxWait time.Duration) error {
	tableName := "users"
	if dialect == "sqlite" {
		tableName = "sl_users"
	}

	// check immediately
	if db.Migrator().HasTable(tableName) {
		return nil
	}

	fmt.Printf("[db] %s table '%s' not found — please run: migrations/001_user.%s.sql\n", dialect, tableName, dialect)
	fmt.Printf("[db] waiting for table '%s' to be created (checking every %v)...\n", tableName, schemaPollInterval)

	deadline := time.Time{}
	if maxWait > 0 {
		deadline = time.Now().Add(maxWait)
	}

	for {
		time.Sleep(schemaPollInterval)

		if db.Migrator().HasTable(tableName) {
			fmt.Printf("[db] %s table '%s' detected, continuing...\n", dialect, tableName)
			return nil
		}

		if maxWait > 0 && time.Now().After(deadline) {
			return fmt.Errorf("%s table '%s' not created within %v", dialect, tableName, maxWait)
		}
	}
}
