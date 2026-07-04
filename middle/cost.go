package middle

import (
	"github.com/labstack/echo/v5"
)

// CostTime cost time output to response.out
func CostTime(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		//start := time.Now()
		//c.Logger().Info("<CostTime>", "cost=====", 11111)
		//err := next(c)
		//cost := time.Since(start).String()
		//c.Logger().Info("<CostTime>", "cost=====", cost)
		//c.Set("i_cost_time", cost)
		return nil
	}
}
