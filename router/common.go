package router

import (
	"github.com/labstack/echo/v5"
	"greeting.first/handler"
)

func common(e *echo.Echo) {
	c := e.Group("/common")
	c.GET("/version", handler.Common.Version)
	c.GET("/changelog", handler.Common.Changelog)
	c.GET("/setting", handler.Common.Setting)
}
