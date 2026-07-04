package response

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
)

// getCost reads start time from context and returns elapsed time string.
func getCost(c *echo.Context) string {
	start, ok := c.Get("i_start_time").(time.Time)
	if !ok {
		return "-"
	}
	return time.Since(start).String()
}

// Ok ...
func Ok(c *echo.Context, data any) error {
	respData := &ErrMsg{
		Code:    ErrCodeOk,
		Message: "",
		Data:    data,
	}
	respData.TraceID = c.Response().Header().Get(echo.HeaderXRequestID)
	respData.Cost = getCost(c)
	return c.JSON(http.StatusOK, respData)
}

// NotOk normal error return response
func NotOk(c *echo.Context, message string) error {
	respData := &ErrMsg{
		Code:    ErrCodeCustom,
		Message: message,
		Data:    "",
	}
	respData.TraceID = c.Response().Header().Get(echo.HeaderXRequestID)
	respData.Cost = getCost(c)
	return c.JSON(http.StatusOK, respData)
}

// NotOkWithCode with special error code
func NotOkWithCode(c *echo.Context, message string, code Code) error {
	respData := &ErrMsg{
		Code:    code,
		Message: message,
		Data:    "",
	}
	respData.TraceID = c.Response().Header().Get(echo.HeaderXRequestID)
	respData.Cost = getCost(c)
	return c.JSON(http.StatusOK, respData)
}
