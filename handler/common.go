package handler

import (
	"runtime"

	"github.com/labstack/echo/v5"
	"greeting.first/entity/common"
	"greeting.first/response"
)

// Common is the handler instance for public/common endpoints.
var Common = &_Common{}

type _Common struct{}

// default values, override via ldflags at build time
var (
	AppVersion = "0.1.0"
	BuildTime  = "unknown"
)

// Version returns application version information.
func (*_Common) Version(c *echo.Context) error {
	return response.Ok(c, common.VersionResponse{
		Version:   AppVersion,
		BuildTime: BuildTime,
		GoVersion: runtime.Version(),
	})
}

// changelog data (static, update as needed).
var changelogData = []common.ChangelogEntry{
	{Date: "2026-07-14", Content: "Add common router with version, changelog, and setting endpoints"},
	{Date: "2026-07-13", Content: "Introduce GORM + SQLite: global model.DB, auto-init on startup"},
	{Date: "2026-07-13", Content: "Add /demo/sha256 endpoint for SHA256 hash computation"},
	{Date: "2026-07-04", Content: "Implement request cost tracking middleware"},
	{Date: "2026-06-30", Content: "Initialize project skeleton with layered architecture"},
}

// Changelog returns the application changelog.
func (*_Common) Changelog(c *echo.Context) error {
	return response.Ok(c, changelogData)
}

// setting data (static, update as needed).
var settingData = []common.SettingItem{
	{Key: "app_name", Value: "Greeting", Description: "Application name"},
	{Key: "port", Value: ":1323", Description: "Server listen port"},
	{Key: "db_type", Value: "sqlite", Description: "Database type"},
	{Key: "db_dsn", Value: "greeting.db", Description: "Database DSN"},
}

// Setting returns the application settings.
func (*_Common) Setting(c *echo.Context) error {
	return response.Ok(c, settingData)
}
