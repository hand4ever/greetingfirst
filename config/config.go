package config

import (
	"fmt"
	"os"
	"path/filepath"

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

// DatabaseConfig holds database connection settings.
type DatabaseConfig struct {
	MySQL  MySQLConfig  `toml:"mysql"`
	SQLite SQLiteConfig `toml:"sqlite"`
}

// MySQLConfig holds MySQL connection settings.
type MySQLConfig struct {
	DSN string `toml:"dsn"`
}

// SQLiteConfig holds SQLite connection settings (independent from MySQL).
type SQLiteConfig struct {
	DSN string `toml:"dsn"`
}

// ChangelogConfig represents a single changelog entry.
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
			{Date: "2026-07-18", Content: "Add mandatory changelog registration rule: all changes and new tasks must appear in GET /common/changelog"},
			{Date: "2026-07-17", Content: "Read server listen port from config.toml in main startup"},
			{Date: "2026-07-17", Content: "Add independent SQLite instance and /sqlite/testuser CRUD test interface"},
			{Date: "2026-07-16", Content: "Add project summary report for specs 001-009 with formatted PDF"},
			{Date: "2026-07-16", Content: "Add missing spec templates and fix bilingual headers in 009 artifacts"},
			{Date: "2026-07-16", Content: "Add DEPLOY_USR variable, runqa shortcut, sudo supervisorctl, and rm old binary before scp"},
			{Date: "2026-07-16", Content: "Amend constitution to v1.3.3 (clarify middleware abort rule + text localization)"},
			{Date: "2026-07-16", Content: "Localize speckit templates with Chinese(English) bilingual headings"},
			{Date: "2026-07-16", Content: "Optimize Makefile with 9 documented targets, variable-driven config, and deploy guard"},
			{Date: "2026-07-15", Content: "Migrate to MySQL-only database architecture"},
			{Date: "2026-07-15", Content: "Add MySQL CRUD endpoints (/demo/usr)"},
			{Date: "2026-07-14", Content: "Add common router with version, changelog, and setting endpoints"},
			{Date: "2026-07-13", Content: "Add /demo/sha256 endpoint for SHA256 hash computation"},
			{Date: "2026-07-04", Content: "Implement request cost tracking middleware"},
			{Date: "2026-06-30", Content: "Initialize project skeleton with layered architecture"},
		},
	}
}

// InitConfig loads config from the given TOML file path, plus the changelog
// from a sibling changelog.toml file. Falls back to defaults if missing/invalid.
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

	// changelog is kept in a separate file (changelog.toml) next to config.toml
	loadChangelog(filepath.Dir(configPath))
	return nil
}

// loadChangelog populates Cfg.Changelog from changelog.toml located in the
// same directory as the main config file. On any error it falls back to the
// default changelog defined in defaultConfig (fail-soft).
func loadChangelog(configDir string) {
	path := filepath.Join(configDir, "changelog.toml")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("[config] changelog file not found at %s, using default changelog\n", path)
			return
		}
		fmt.Printf("[config] failed to read %s: %v, using default changelog\n", path, err)
		return
	}

	var cl struct {
		Changelog []ChangelogConfig `toml:"changelog"`
	}
	if err := toml.Unmarshal(data, &cl); err != nil {
		fmt.Printf("[config] failed to parse %s: %v, using default changelog\n", path, err)
		return
	}

	Cfg.Changelog = cl.Changelog
	fmt.Printf("[config] loaded %d changelog entries from %s\n", len(Cfg.Changelog), path)
}
