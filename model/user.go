package model

import (
	"database/sql/driver"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// LocalTime 自定义时间类型，JSON 序列化为 "2006-01-02 15:04:05" 格式
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

// Value 实现 driver.Valuer，GORM 写入数据库时使用
func (t LocalTime) Value() (driver.Value, error) {
	return time.Time(t), nil
}

// Scan 实现 sql.Scanner，GORM 从数据库读取时使用
func (t *LocalTime) Scan(v interface{}) error {
	if tv, ok := v.(time.Time); ok {
		*t = LocalTime(tv)
		return nil
	}
	return fmt.Errorf("cannot scan %T into LocalTime", v)
}

// User table
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Phone     string         `gorm:"uniqueIndex;type:varchar(32);not null" json:"phone"`
	Name      string         `gorm:"type:varchar(64);not null" json:"name"`
	Age       int            `gorm:"default:0" json:"age"`
	CreatedAt LocalTime      `json:"created_at"`
	UpdatedAt LocalTime      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// CreateUser insert a new user
func CreateUser(user *User) error {
	return DB.Create(user).Error
}

// GetUserByID query user by id
func GetUserByID(id uint) (*User, error) {
	var user User
	err := DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByPhone query user by phone
func GetUserByPhone(phone string) (*User, error) {
	var user User
	err := DB.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser update user fields
func UpdateUser(user *User) error {
	return DB.Save(user).Error
}

// DeleteUser soft-delete user by id
func DeleteUser(id uint) error {
	return DB.Delete(&User{}, id).Error
}

// RestoreUserByPhone restore a soft-deleted user by phone
func RestoreUserByPhone(phone string) (*User, error) {
	var user User
	err := DB.Unscoped().Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	user.DeletedAt = gorm.DeletedAt{}
	if err := DB.Unscoped().Model(&user).Update("deleted_at", nil).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
