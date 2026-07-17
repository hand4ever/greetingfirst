package router

import (
	"github.com/labstack/echo/v5"
)

func Router(e *echo.Echo) {
	demo(e)   // demo 相关
	common(e) // common public components
	sqlite(e) // SQLite 测试接口
}
