package model

import (
	"time"
)

// Statistics 访问统计
type Statistics struct {
	ID             uint64    `gorm:"primaryKey" json:"id"`
	StatDate       time.Time `gorm:"type:date;uniqueIndex" json:"stat_date"`
	PageViews      int       `gorm:"default:0" json:"page_views"`
	UniqueVisitors int       `gorm:"default:0" json:"unique_visitors"`
	NewPosts       int       `gorm:"default:0" json:"new_posts"`
	NewComments    int       `gorm:"default:0" json:"new_comments"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// TableName 表名
func (Statistics) TableName() string {
	return "statistics"
}

// Setting 系统配置
type Setting struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	SettingKey  string    `gorm:"size:100;uniqueIndex" json:"setting_key"`
	SettingValue string   `gorm:"type:text" json:"setting_value"`
	SettingType string    `gorm:"size:20;default:string" json:"setting_type"`
	Description string    `gorm:"size:200" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 表名
func (Setting) TableName() string {
	return "settings"
}

// OperationLog 操作日志
type OperationLog struct {
	ID           uint64     `gorm:"primaryKey" json:"id"`
	UserID       *uint64    `gorm:"index" json:"user_id"`
	Action       string     `gorm:"size:50;not null" json:"action"`
	Module       string     `gorm:"size:50" json:"module"`
	Description  string     `gorm:"size:500" json:"description"`
	IPAddress    string     `gorm:"size:45" json:"ip_address"`
	UserAgent    string     `gorm:"size:500" json:"user_agent"`
	RequestData  string     `gorm:"type:json" json:"request_data"`
	ResponseCode int        `json:"response_code"`
	CreatedAt    time.Time  `gorm:"index" json:"created_at"`
}

// TableName 表名
func (OperationLog) TableName() string {
	return "operation_logs"
}
