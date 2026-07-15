package handler

import (
	"errors"
	"strconv"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
	"greeting.first/entity/user"
	"greeting.first/model"
	"greeting.first/response"
)

// User is the handler instance for MySQL user CRUD under /demo/usr.
var User = &_User{}

type _User struct{}

// Create handles POST /demo/usr
func (*_User) Create(c *echo.Context) error {
	var req user.UserCreateReq
	if err := c.Bind(&req); err != nil {
		return response.NotOk(c, "invalid request body")
	}

	// validate required fields
	if req.Name == "" || req.Phone == "" {
		return response.NotOk(c, "name and phone are required")
	}

	u := &model.User{
		Name:  req.Name,
		Phone: req.Phone,
	}
	if req.Age != nil {
		u.Age = *req.Age
	}

	if err := model.CreateUser(u); err != nil {
		return response.NotOk(c, "phone already exists")
	}

	return response.Ok(c, u)
}

// extractUserID extracts the :id path parameter as uint.
func extractUserID(c *echo.Context) (uint, error) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// Get handles GET /demo/usr/:id
func (*_User) Get(c *echo.Context) error {
	id, err := extractUserID(c)
	if err != nil {
		return response.NotOk(c, "invalid path parameter")
	}

	u, err := model.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NotOk(c, "user not found")
		}
		return response.NotOk(c, "query user failed: "+err.Error())
	}

	return response.Ok(c, u)
}

// Update handles PUT /demo/usr/:id
func (*_User) Update(c *echo.Context) error {
	id, err := extractUserID(c)
	if err != nil {
		return response.NotOk(c, "invalid path parameter")
	}

	var req user.UserUpdateReq
	if err := c.Bind(&req); err != nil {
		return response.NotOk(c, "invalid request body")
	}

	u, err := model.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NotOk(c, "user not found")
		}
		return response.NotOk(c, "query user failed: "+err.Error())
	}

	// partial update: only update provided fields
	if req.Name != nil {
		u.Name = *req.Name
	}
	if req.Age != nil {
		u.Age = *req.Age
	}

	if err := model.UpdateUser(u); err != nil {
		return response.NotOk(c, "update user failed: "+err.Error())
	}

	return response.Ok(c, u)
}

// Delete handles DELETE /demo/usr/:id
func (*_User) Delete(c *echo.Context) error {
	id, err := extractUserID(c)
	if err != nil {
		return response.NotOk(c, "invalid path parameter")
	}

	if err := model.DeleteUser(id); err != nil {
		return response.NotOk(c, "delete user failed: "+err.Error())
	}

	return response.Ok(c, "")
}

// List handles GET /demo/usrs
func (*_User) List(c *echo.Context) error {
	var users []model.User
	if err := model.DB.Where("deleted_at IS NULL").Order("created_at DESC").Find(&users).Error; err != nil {
		return response.NotOk(c, "query users failed: "+err.Error())
	}

	return response.Ok(c, users)
}
