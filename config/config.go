package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// Config holds all configuration sections.
type Config struct {
	App       AppConfig         `toml:"app"`
	Server    ServerConfig      `toml:"server"`
	Database  DatabaseConfig    `toml:"database"`
	Changelog []ChangelogConfig `toml:"changelog"`
}

// AppConfig holds application metadata.
type AppConfig struct {
	Name      string `toml:"name"`
	Version   string `toml:"version"`
	BuildTime string `toml:"build_time"`
}

// ServerConfig holds server settings.
type ServerConfig struct {
	Port string `toml:"port"`
}

// DatabaseConfig holds database connection settings for MySQL and SQLite.
type DatabaseConfig struct {
	MySQL  MySQLConfig  `toml:"mysql"`
	SQLite SQLiteConfig `toml:"sqlite"`
}

// MySQLConfig holds MySQL connection settings.
type MySQLConfig struct {
	DSN string `toml:"dsn"`
}

// SQLiteConfig holds SQLite connection settings.
type SQLiteConfig struct {
	DSN string `toml:"dsn"`
}

// ChangelogConfig represents a single changelog entry in config.
type ChangelogConfig struct {
	Date    string `toml:"date"`
	Content string `toml:"content"`
}

// Cfg is the global config instance, initialized at startup.
var Cfg = defaultConfig()

// defaultConfig returns safe defaults when config file is missing or invalid.
func defaultConfig() *Config {
	return &Config{
		App: AppConfig{
			Name:      "Greeting",
			Version:   "0.1.0",
			BuildTime: "unknown",
		},
		Server: ServerConfig{
			Port: ":1323",
		},
		Database: DatabaseConfig{
			MySQL: MySQLConfig{
				DSN: "root:password@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local",
			},
			SQLite: SQLiteConfig{
				DSN: "greeting.db",
			},
		},
		Changelog: []ChangelogConfig{
			{Date: "2026-07-14", Content: "Add common router with version, changelog, and setting endpoints"},
			{Date: "2026-07-13", Content: "Introduce GORM + SQLite: global model.DB, auto-init on startup"},
			{Date: "2026-07-13", Content: "Add /demo/sha256 endpoint for SHA256 hash computation"},
			{Date: "2026-07-04", Content: "Implement request cost tracking middleware"},
			{Date: "2026-06-30", Content: "Initialize project skeleton with layered architecture"},
		},
	}
}

// InitConfig loads config from the given TOML file path.
// Falls back to defaults if the file is missing or invalid.
func InitConfig(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("[config] config file not found at %s, using defaults\n", configPath)
			return nil
		}
		return fmt.Errorf("read config file: %w", err)
	}

	if err := toml.Unmarshal(data, Cfg); err != nil {
		fmt.Printf("[config] failed to parse %s: %v, using defaults\n", configPath, err)
		return nil
	}

	fmt.Printf("[config] loaded config from %s\n", configPath)
	return nil
}
