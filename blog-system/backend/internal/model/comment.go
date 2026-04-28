package model

import (
	"time"
)

// Comment 评论
type Comment struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	PostID      uint64    `gorm:"index;not null" json:"post_id"`
	ParentID    *uint64   `gorm:"index" json:"parent_id"`
	UserName    string    `gorm:"size:50;not null" json:"user_name"`
	UserEmail   string    `gorm:"size:100" json:"user_email"`
	UserIP      string    `gorm:"size:45;not null" json:"user_ip"`
	UserLocation string   `gorm:"size:100" json:"user_location"`
	UserAgent   string    `gorm:"size:500" json:"user_agent"`
	Content     string    `gorm:"type:text;not null" json:"content"`
	Status      int8      `gorm:"default:1;index" json:"status"`
	LikeCount   int       `gorm:"default:0" json:"like_count"`
	IsAdmin     int8      `gorm:"default:0" json:"is_admin"`
	CreatedAt   time.Time `gorm:"index" json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 表名
func (Comment) TableName() string {
	return "comments"
}

// Like 点赞记录
type Like struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	UserID    *uint64   `json:"user_id"`
	SessionID string    `gorm:"size:64" json:"session_id"`
	PostID    *uint64   `gorm:"index" json:"post_id"`
	CommentID *uint64   `gorm:"index" json:"comment_id"`
	IPAddress string    `gorm:"size:45" json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 表名
func (Like) TableName() string {
	return "likes"
}

// Upload 上传文件
type Upload struct {
	ID           uint64    `gorm:"primaryKey" json:"id"`
	Filename     string    `gorm:"size:255;not null" json:"filename"`
	OriginalName string    `gorm:"size:255" json:"original_name"`
	FilePath     string    `gorm:"size:500;not null" json:"file_path"`
	FileURL      string    `gorm:"size:500;not null" json:"file_url"`
	FileType     string    `gorm:"size:50;not null" json:"file_type"`
	FileSize     int64     `gorm:"not null" json:"file_size"`
	MimeType     string    `gorm:"size:100" json:"mime_type"`
	UploaderID   *uint64   `gorm:"index" json:"uploader_id"`
	PostID       *uint64   `gorm:"index" json:"post_id"`
	Status       int8      `gorm:"default:1" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

// TableName 表名
func (Upload) TableName() string {
	return "uploads"
}
