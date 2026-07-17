package model

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global MySQL database instance.
var DB *gorm.DB

// SQLiteDB is the global SQLite database instance, coexisting with DB (MySQL).
// It backs the test-only test_user CRUD interface and is fully isolated from MySQL.
var SQLiteDB *gorm.DB

// InitDB opens a MySQL connection.
// It uses fail-fast: connection failure causes an immediate error.
// No table existence check is performed — schema is user-managed.
func InitDB(mysqlDSN string) error {
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

	return nil
}

// InitSQLite opens a SQLite connection and verifies the user-owned test_user
// table exists. It uses fail-fast: connection or table-check failure causes an
// immediate error. No table creation or migration is performed (principle VII).
func InitSQLite(dsn string) error {
	sqliteDB, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("sqlite (%s): %w", dsn, err)
	}
	sqlDB, err := sqliteDB.DB()
	if err != nil {
		return fmt.Errorf("sqlite get underlying db: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("sqlite ping (%s): %w", dsn, err)
	}
	// Verify the user-owned test_user table exists. The application MUST NOT
	// create or migrate it (principle VII); missing table is a fail-fast error.
	if err := sqliteDB.Exec("SELECT 1 FROM test_user LIMIT 1").Error; err != nil {
		return fmt.Errorf("sqlite table test_user missing (run migrations/002_test_user.sql): %w", err)
	}
	SQLiteDB = sqliteDB

	return nil
}
