package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"greeting.first/entity/demo"
	"greeting.first/response"
)

var Demo = &_Demo{}

type _Demo struct{}

// Search ...
func (*_Demo) Search(c *echo.Context) error {
	var f demo.Filter
	if err := c.Bind(&f); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	return response.Ok(c, f)
}
