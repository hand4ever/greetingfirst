package model

import (
	"database/sql/driver"
	"fmt"
	"time"

	"gorm.io/gorm"
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
// MySQL Entity: User
// ============================================================================

// User is the MySQL user entity (table: users).
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Phone     string         `gorm:"type:varchar(32);not null" json:"phone"`
	Name      string         `gorm:"type:varchar(64);not null" json:"name"`
	Age       int            `gorm:"default:0" json:"age"`
	CreatedAt LocalTime      `json:"created_at"`
	UpdatedAt LocalTime      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// CreateUser inserts a new user into MySQL.
func CreateUser(user *User) error {
	return DB.Create(user).Error
}

// GetUserByID queries a user by id from MySQL.
func GetUserByID(id uint) (*User, error) {
	var user User
	err := DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByPhone queries a user by phone from MySQL (non-deleted only).
func GetUserByPhone(phone string) (*User, error) {
	var user User
	err := DB.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates user fields in MySQL.
func UpdateUser(user *User) error {
	return DB.Save(user).Error
}

// DeleteUser soft-deletes a user by id in MySQL.
func DeleteUser(id uint) error {
	return DB.Delete(&User{}, id).Error
}

// ============================================================================
// SQLite Entity: SQLiteUser
// ============================================================================

// SQLiteUser is the SQLite user entity (table: sl_users, renamed from User).
type SQLiteUser struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Phone     string         `gorm:"type:varchar(32);not null" json:"phone"`
	Name      string         `gorm:"type:varchar(64);not null" json:"name"`
	Age       int            `gorm:"default:0" json:"age"`
	CreatedAt LocalTime      `json:"created_at"`
	UpdatedAt LocalTime      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName overrides the default table name.
func (SQLiteUser) TableName() string {
	return "sl_users"
}

// CreateSQLiteUser inserts a new user into SQLite.
func CreateSQLiteUser(user *SQLiteUser) error {
	return SQLiteDB.Create(user).Error
}

// GetSQLiteUserByID queries a user by id from SQLite.
func GetSQLiteUserByID(id uint) (*SQLiteUser, error) {
	var user SQLiteUser
	err := SQLiteDB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetSQLiteUserByPhone queries a user by phone from SQLite (non-deleted only).
func GetSQLiteUserByPhone(phone string) (*SQLiteUser, error) {
	var user SQLiteUser
	err := SQLiteDB.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateSQLiteUser updates user fields in SQLite.
func UpdateSQLiteUser(user *SQLiteUser) error {
	return SQLiteDB.Save(user).Error
}

// DeleteSQLiteUser soft-deletes a user by id in SQLite.
func DeleteSQLiteUser(id uint) error {
	return SQLiteDB.Delete(&SQLiteUser{}, id).Error
}
