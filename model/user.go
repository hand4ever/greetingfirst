package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// LocalTime custom time type for JSON serialization in "2006-01-02 15:04:05" format.
type LocalTime time.Time

const localTimeFormat = "2006-01-02 15:04:05"

func (t LocalTime) MarshalJSON() ([]byte, error) {
	s := time.Time(t).Format(localTimeFormat)
	return []byte(`"` + s + `"`), nil
}

func (t *LocalTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		return nil
	}
	parsed, err := time.ParseInLocation(`"`+localTimeFormat+`"`, string(data), time.Local)
	if err != nil {
		return err
	}
	*t = LocalTime(parsed)
	return nil
}

// Value implements driver.Valuer for GORM writes.
func (t LocalTime) Value() (driver.Value, error) {
	return time.Time(t), nil
}

// Scan implements sql.Scanner for GORM reads.
func (t *LocalTime) Scan(v interface{}) error {
	if tv, ok := v.(time.Time); ok {
		*t = LocalTime(tv)
		return nil
	}
	return fmt.Errorf("cannot scan %T into LocalTime", v)
}

// ============================================================================
// User Entity
// ============================================================================

// User maps to the users table.
type User struct {
	ID           int        `gorm:"primaryKey" json:"id"`
	Phone        string     `gorm:"type:varchar(20);not null" json:"phone"`
	Realname     string     `gorm:"type:varchar(100)" json:"realname"`
	Username     string     `gorm:"type:varchar(20)" json:"username"`
	Age          int        `gorm:"default:0" json:"age"`
	PasswordHash string     `gorm:"type:varchar(200)" json:"-"`
	CreatedAt    LocalTime  `json:"created_at"`
	UpdatedAt    LocalTime  `json:"updated_at"`
	DeletedAt    *time.Time `json:"-"`
}

// CreateUser inserts a new user.
func CreateUser(user *User) error {
	return DB.Create(user).Error
}

// GetUserByID queries a non-deleted user by id.
func GetUserByID(id int) (*User, error) {
	var user User
	err := DB.Where("id = ? AND deleted_at IS NULL", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByPhone queries a non-deleted user by phone.
func GetUserByPhone(phone string) (*User, error) {
	var user User
	err := DB.Where("phone = ? AND deleted_at IS NULL", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates user fields.
func UpdateUser(user *User) error {
	return DB.Save(user).Error
}

// DeleteUser soft-deletes a user by id (sets deleted_at = NOW()).
func DeleteUser(id int) error {
	now := time.Now()
	return DB.Model(&User{}).Where("id = ?", id).Update("deleted_at", now).Error
}
