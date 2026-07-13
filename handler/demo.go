package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"greeting.first/entity/demo"
	"greeting.first/model"
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

// ErrDebug ...
func (*_Demo) ErrDebug(c *echo.Context) error {
	var s demo.Echo
	tid := c.Response().Header().Get(echo.HeaderXRequestID)

	if err := c.Bind(&s); err != nil {
		return response.NotOk(c, "参数有误")
	}
	time.Sleep(time.Millisecond * 300)
	c.Logger().Info("<ErrDebug>", "request", s, "tid", tid)
	return response.Ok(c, s)
}

const testPhone = "13636311005"

// GetUserByPhoneTest test-only: query user by phone, insert if not exists
func (*_Demo) GetUserByPhoneTest(c *echo.Context) error {
	user, err := model.GetUserByPhone(testPhone)
	if err != nil {
		// user not found (maybe soft-deleted), create a test user
		user = &model.User{
			Phone: testPhone,
			Name:  "test_user",
			Age:   25,
		}
		if err := model.CreateUser(user); err != nil {
			// duplicate phone (soft-deleted record exists), restore it
			if restored, restoreErr := model.RestoreUserByPhone(testPhone); restoreErr == nil {
				c.Logger().Info("<GetUserByPhoneTest> restored test user", "id", restored.ID, "phone", testPhone)
				return response.Ok(c, restored)
			}
			return response.NotOk(c, "create test user failed: "+err.Error())
		}
		c.Logger().Info("<GetUserByPhoneTest> created test user", "id", user.ID, "phone", testPhone)
		return response.Ok(c, user)
	}

	c.Logger().Info("<GetUserByPhoneTest> found user by phone", "id", user.ID, "phone", testPhone)
	return response.Ok(c, user)
}
