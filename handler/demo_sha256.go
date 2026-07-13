package handler

import (
	"crypto/sha256"
	"fmt"

	"github.com/labstack/echo/v5"
	"greeting.first/entity/demo"
	"greeting.first/response"
)

// Sha256 is the handler instance for SHA256 computation endpoints.
var Sha256 = &_Sha256{}

type _Sha256 struct{}

// Compute calculates SHA256 hash of the `text` query parameter and returns
// both the original input and the lowercase hexadecimal hash.
func (*_Sha256) Compute(c *echo.Context) error {
	// check if text parameter exists in query string
	if !c.Request().URL.Query().Has("text") {
		return response.NotOk(c, "text parameter is required")
	}

	var req demo.Sha256Request
	if err := c.Bind(&req); err != nil {
		return response.NotOk(c, "text parameter is required")
	}

	sum := sha256.Sum256([]byte(req.Text))
	hash := fmt.Sprintf("%x", sum)

	return response.Ok(c, demo.Sha256Response{
		Input: req.Text,
		Hash:  hash,
	})
}
