package handler

import (
	"runtime"

	"github.com/labstack/echo/v5"
	"greeting.first/config"
	"greeting.first/entity/common"
	"greeting.first/response"
)

// Common is the handler instance for public/common endpoints.
var Common = &_Common{}

type _Common struct{}

// Version returns application version information from config.
func (*_Common) Version(c *echo.Context) error {
	cfg := config.Cfg
	return response.Ok(c, common.VersionResponse{
		Version:   cfg.App.Version,
		BuildTime: cfg.App.BuildTime,
		GoVersion: runtime.Version(),
	})
}

// Changelog returns the application changelog from config.
func (*_Common) Changelog(c *echo.Context) error {
	entries := make([]common.ChangelogEntry, 0, len(config.Cfg.Changelog))
	for _, e := range config.Cfg.Changelog {
		entries = append(entries, common.ChangelogEntry{
			Date:    e.Date,
			Content: e.Content,
		})
	}
	return response.Ok(c, entries)
}

// Setting returns the application settings from config.
func (*_Common) Setting(c *echo.Context) error {
	cfg := config.Cfg
	items := []common.SettingItem{
		{Key: "app_name", Value: cfg.App.Name, Description: "Application name"},
		{Key: "app_version", Value: cfg.App.Version, Description: "Application version"},
		{Key: "server_port", Value: cfg.Server.Port, Description: "Server listen port"},
		{Key: "db_type", Value: cfg.Database.Type, Description: "Database type"},
		{Key: "db_dsn", Value: cfg.Database.DSN, Description: "Database DSN"},
	}
	return response.Ok(c, items)
}
