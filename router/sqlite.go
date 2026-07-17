package router

import (
	"github.com/labstack/echo/v5"
	"greeting.first/handler"
)

// sqlite registers the SQLite test_user CRUD routes, isolated from MySQL.
// Create/Get/Update/Delete live under /sqlite/testuser; the list endpoint is
// the sibling /sqlite/testusers (principle I: layered, registered centrally).
func sqlite(e *echo.Echo) {
	g := e.Group("/sqlite/testuser")
	g.POST("", handler.SqliteUser.Create)
	g.GET("/:id", handler.SqliteUser.Get)
	g.PUT("/:id", handler.SqliteUser.Update)
	g.DELETE("/:id", handler.SqliteUser.Delete)
	e.GET("/sqlite/testusers", handler.SqliteUser.List)
}
