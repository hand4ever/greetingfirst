package middle

import (
	"time"

	"github.com/labstack/echo/v5"
)

// CostTime cost time output to response.out
func CostTime(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		start := time.Now()
		c.Set("i_start_time", start)
		err := next(c)
		cost := time.Since(start).String()
		c.Logger().Info("<CostTime>", "cost", cost)
		return err
	}
}
