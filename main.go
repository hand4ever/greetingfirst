package main

import (
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"greeting.first/router"
)

func main() {
	e := echo.New()

	//e.HTTPErrorHandler = middle.CustomHTTPErrorHandler
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	//e.Use(middle.CostTime)

	router.Router(e)

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
