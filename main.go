package main

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"greeting.first/config"
	"greeting.first/middle"
	"greeting.first/model"
	"greeting.first/router"
)

// corsConfig defines the CORS middleware configuration.
// Modify these values to restrict allowed origins, methods, or headers for production.
var corsConfig = middleware.CORSConfig{
	AllowOrigins:     []string{"*"},
	AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions, http.MethodPatch, http.MethodHead},
	AllowHeaders:     []string{"Content-Type", "Authorization", "X-Requested-With", "Accept", "Origin"},
	AllowCredentials: false,
	MaxAge:           86400,
}

func main() {
	// init config
	if err := config.InitConfig("config.toml"); err != nil {
		panic("failed to load config: " + err.Error())
	}

	// init database: connect both MySQL and SQLite
	if err := model.InitDB(config.Cfg.Database.MySQL.DSN, config.Cfg.Database.SQLite.DSN); err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// ensure MySQL users table exists (pause-and-continue if missing)
	if err := model.EnsureUserTable(model.DB, "mysql", 0); err != nil {
		panic("failed to ensure mysql users table: " + err.Error())
	}

	e := echo.New()

	//e.HTTPErrorHandler = middle.CustomHTTPErrorHandler
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(corsConfig))
	e.Use(middleware.RequestID())
	e.Use(middle.CostTime)

	router.Router(e)

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
