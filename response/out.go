package response

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

var Out = &_Out{}

type _Out struct{}

// Ok ...
func Ok(c *echo.Context, data any) error {
	return c.JSON(http.StatusOK, data)
}
