package router

import (
	"github.com/labstack/echo/v5"
	"greeting.first/handler"
)

func demo(e *echo.Echo) {
	d := e.Group("/demo")
	d.GET("/search", handler.Demo.Search)
	d.GET("/err/debug/:str", handler.Demo.ErrDebug)
}
