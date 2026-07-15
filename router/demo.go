package router

import (
	"github.com/labstack/echo/v5"
	"greeting.first/handler"
)

func demo(e *echo.Echo) {
	d := e.Group("/demo")
	d.GET("/search", handler.Demo.Search)
	d.GET("/err/debug/:str", handler.Demo.ErrDebug)
	d.GET("/user/phone", handler.Demo.GetUserByPhoneTest)
	d.GET("/sha256", handler.Sha256.Compute)

	// MySQL CRUD endpoints under /demo/usr
	u := e.Group("/demo/usr")
	u.POST("", handler.User.Create)
	u.GET("/:id", handler.User.Get)
	u.PUT("/:id", handler.User.Update)
	u.DELETE("/:id", handler.User.Delete)
	u.GET("s", handler.User.List) // /demo/usrs
}
