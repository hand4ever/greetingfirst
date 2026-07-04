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
		Cost:    "-0.111s",
	}
	respData.TraceID = c.Response().Header().Get(echo.HeaderXRequestID)
	//respData.Cost = c.Get("i_cost_time").(string)
	return c.JSON(http.StatusOK, respData)
}

// NotOk normal error return response
func NotOk(c *echo.Context, message string) error {
	respData := &ErrMsg{
		Code:    ErrCodeCustom,
		Message: message,
		Data:    "",
		Cost:    "-0.111s",
	}
	respData.TraceID = c.Response().Header().Get(echo.HeaderXRequestID)
	return c.JSON(http.StatusOK, respData)
}

// NotOkWithCode with special error code
func NotOkWithCode(c *echo.Context, message string, code Code) error {
	respData := &ErrMsg{
		Code:    code,
		Message: message,
		Data:    "",
		Cost:    "-0.111222s",
	}

	respData.TraceID = c.Response().Header().Get(echo.HeaderXRequestID)
	return c.JSON(http.StatusOK, respData)
}
