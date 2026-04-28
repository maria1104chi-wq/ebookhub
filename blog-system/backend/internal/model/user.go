package model

import (
	"time"
)

// User 管理员用户
type User struct {
	ID           uint64     `gorm:"primaryKey" json:"id"`
	Username     string     `gorm:"size:50;uniqueIndex" json:"username"`
	PasswordHash string     `gorm:"size:255" json:"-"`
	Phone        string     `gorm:"size:20;uniqueIndex" json:"phone"`
	Email        string     `gorm:"size:100" json:"email"`
	Avatar       string     `gorm:"size:255" json:"avatar"`
	Role         int8       `gorm:"default:1" json:"role"`
	Status       int8       `gorm:"default:1" json:"status"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	LastLoginIP  string     `gorm:"size:45" json:"last_login_ip"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// TableName 表名
func (User) TableName() string {
	return "users"
}

// SMSCode 短信验证码
type SMSCode struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Phone     string    `gorm:"size:20;index" json:"phone"`
	Code      string    `gorm:"size:6" json:"-"`
	Purpose   string    `gorm:"size:20" json:"purpose"`
	ExpiresAt time.Time `gorm:"index" json:"expires_at"`
	Used      int8      `gorm:"default:0;index" json:"used"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 表名
func (SMSCode) TableName() string {
	return "sms_codes"
}
