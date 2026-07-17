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
	e := echo.New()
	// init config
	if err := config.InitConfig("config.toml"); err != nil {
		panic("failed to load config: " + err.Error())
	}

	// init database: connect to MySQL
	if err := model.InitDB(config.Cfg.Database.MySQL.DSN); err != nil {
		// panic("failed to connect database: " + err.Error())
		e.Logger.Error("failed to connect database", "error", err)
	}

	// init SQLite instance (coexists with MySQL), fail-fast on error
	if err := model.InitSQLite(config.Cfg.Database.SQLite.DSN); err != nil {
		panic("failed to connect sqlite: " + err.Error())
	}

	//e.HTTPErrorHandler = middle.CustomHTTPErrorHandler
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(corsConfig))
	e.Use(middleware.RequestID())
	e.Use(middle.CostTime)

	router.Router(e)

	// start HTTP server on the port configured in [server].port (config.toml)
	if err := e.Start(config.Cfg.Server.Port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
