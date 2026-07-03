package response

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

// Ok ...
func Ok(c *echo.Context, data any) error {
	respData := &ErrMsg{
		Code:    ErrCodeOk,
		Message: "",
		Data:    data,
		TraceID: "dummy",
		Cost:    "-0.111s",
	}
	return c.JSON(http.StatusOK, respData)
}

// NotOk normal error return response
func NotOk(c *echo.Context, message string) error {
	respData := &ErrMsg{
		Code:    ErrCodeCustom,
		Message: message,
		Data:    "",
		TraceID: "dummy",
		Cost:    "-0.111s",
	}
	return c.JSON(http.StatusOK, respData)
}

// NotOkWithCode with special error code
func NotOkWithCode(c *echo.Context, message string, code Code) error {
	respData := &ErrMsg{
		Code:    ErrCodeCustom,
		Message: message,
		Data:    "",
		TraceID: "dummy111",
		Cost:    "-0.111222s",
	}
	return c.JSON(http.StatusOK, respData)
}
