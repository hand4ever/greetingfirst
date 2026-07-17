package handler

import (
	"errors"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
	"greeting.first/entity/sqliteusr"
	"greeting.first/model"
	"greeting.first/response"
)

// SqliteUser is the handler instance for the SQLite test_user CRUD under /sqlite/testuser.
var SqliteUser = &_SqliteUser{}

type _SqliteUser struct{}

// Create handles POST /sqlite/testuser
func (*_SqliteUser) Create(c *echo.Context) error {
	var req sqliteusr.TestUserCreateReq
	if err := c.Bind(&req); err != nil {
		return response.NotOk(c, "invalid request body")
	}
	if req.Name == "" {
		return response.NotOk(c, "name is required")
	}
	if req.Phone == "" {
		return response.NotOk(c, "phone is required")
	}

	exists, err := model.PhoneActiveExists(req.Phone)
	if err != nil {
		return response.NotOk(c, "check phone failed: "+err.Error())
	}
	if exists {
		return response.NotOk(c, "phone already exists")
	}

	t := &model.TestUser{
		Name:  req.Name,
		Phone: req.Phone,
	}
	if req.Age != nil {
		t.Age = *req.Age
	}
	if err := model.CreateTestUser(t); err != nil {
		return response.NotOk(c, "create test user failed: "+err.Error())
	}
	return response.Ok(c, t)
}

// Get handles GET /sqlite/testuser/:id
func (*_SqliteUser) Get(c *echo.Context) error {
	id, err := extractUserID(c)
	if err != nil {
		return response.NotOk(c, "invalid path parameter")
	}
	t, err := model.GetTestUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NotOk(c, "user not found")
		}
		return response.NotOk(c, "query test user failed: "+err.Error())
	}
	return response.Ok(c, t)
}

// Update handles PUT /sqlite/testuser/:id (partial update)
func (*_SqliteUser) Update(c *echo.Context) error {
	id, err := extractUserID(c)
	if err != nil {
		return response.NotOk(c, "invalid path parameter")
	}
	var req sqliteusr.TestUserUpdateReq
	if err := c.Bind(&req); err != nil {
		return response.NotOk(c, "invalid request body")
	}
	t, err := model.GetTestUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NotOk(c, "user not found")
		}
		return response.NotOk(c, "query test user failed: "+err.Error())
	}
	if req.Name != nil {
		t.Name = *req.Name
	}
	if req.Phone != nil {
		t.Phone = *req.Phone
	}
	if req.Age != nil {
		t.Age = *req.Age
	}
	if err := model.UpdateTestUser(t); err != nil {
		return response.NotOk(c, "update test user failed: "+err.Error())
	}
	return response.Ok(c, t)
}

// Delete handles DELETE /sqlite/testuser/:id (soft delete)
func (*_SqliteUser) Delete(c *echo.Context) error {
	id, err := extractUserID(c)
	if err != nil {
		return response.NotOk(c, "invalid path parameter")
	}
	if err := model.DeleteTestUser(id); err != nil {
		return response.NotOk(c, "delete test user failed: "+err.Error())
	}
	return response.Ok(c, "")
}

// List handles GET /sqlite/testusers
func (*_SqliteUser) List(c *echo.Context) error {
	ts, err := model.ListTestUsers()
	if err != nil {
		return response.NotOk(c, "query users failed: "+err.Error())
	}
	return response.Ok(c, ts)
}
