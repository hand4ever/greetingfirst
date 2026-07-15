package model

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global database instance.
var DB *gorm.DB

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
