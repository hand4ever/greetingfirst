package model

import (
	"gorm.io/gorm"
)

// TestUser maps to the user-owned test_user table (SQLite only, test interface).
// It is fully isolated from the MySQL User entity and serves the test CRUD API.
type TestUser struct {
	ID        int            `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Phone     string         `gorm:"type:varchar(20);not null" json:"phone"`
	Age       int            `gorm:"default:0" json:"age"`
	CreatedAt LocalTime      `json:"created_at"`
	UpdatedAt LocalTime      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// TableName returns the explicit table name for the TestUser model.
func (TestUser) TableName() string {
	return "test_user"
}

// CreateTestUser inserts a new test user.
func CreateTestUser(t *TestUser) error {
	return SQLiteDB.Create(t).Error
}

// GetTestUserByID queries a non-deleted test user by id.
// Returns gorm.ErrRecordNotFound when missing or soft-deleted.
func GetTestUserByID(id int) (*TestUser, error) {
	var t TestUser
	if err := SQLiteDB.First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

// UpdateTestUser persists changes to a test user.
func UpdateTestUser(t *TestUser) error {
	return SQLiteDB.Save(t).Error
}

// DeleteTestUser soft-deletes a test user by id (sets deleted_at).
func DeleteTestUser(id int) error {
	return SQLiteDB.Delete(&TestUser{}, id).Error
}

// ListTestUsers returns all non-deleted test users ordered by created_at DESC.
func ListTestUsers() ([]TestUser, error) {
	var ts []TestUser
	if err := SQLiteDB.Order("created_at DESC").Find(&ts).Error; err != nil {
		return nil, err
	}
	return ts, nil
}

// PhoneActiveExists reports whether an active (non-deleted) test user with the
// given phone already exists. Used before creation to enforce uniqueness while
// still allowing reuse after a soft delete.
func PhoneActiveExists(phone string) (bool, error) {
	var count int64
	err := SQLiteDB.Model(&TestUser{}).
		Where("phone = ? AND deleted_at IS NULL", phone).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
