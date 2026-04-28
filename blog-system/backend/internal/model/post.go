package model

import (
	"time"
)

// Post 博文
type Post struct {
	ID             uint64    `gorm:"primaryKey" json:"id"`
	Title          string    `gorm:"size:200;not null" json:"title"`
	Slug           string    `gorm:"size:200;uniqueIndex" json:"slug"`
	Summary        string    `gorm:"size:500" json:"summary"`
	Content        string    `gorm:"type:longtext;not null" json:"content"`
	ContentHTML    string    `gorm:"type:longtext" json:"content_html"`
	CoverImage     string    `gorm:"size:255" json:"cover_image"`
	AuthorID       uint64    `gorm:"index" json:"author_id"`
	CategoryID     uint64    `gorm:"index" json:"category_id"`
	Tags           string    `gorm:"type:json" json:"tags"`
	ViewCount      int       `gorm:"default:0;index" json:"view_count"`
	LikeCount      int       `gorm:"default:0" json:"like_count"`
	ShareCount     int       `gorm:"default:0" json:"share_count"`
	CommentCount   int       `gorm:"default:0" json:"comment_count"`
	Status         int8      `gorm:"default:1;index" json:"status"`
	IsTop          int8      `gorm:"default:0;index" json:"is_top"`
	SeoTitle       string    `gorm:"size:200" json:"seo_title"`
	SeoKeywords    string    `gorm:"size:500" json:"seo_keywords"`
	SeoDescription string    `gorm:"size:500" json:"seo_description"`
	PublishedAt    *time.Time `gorm:"index" json:"published_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// TableName 表名
func (Post) TableName() string {
	return "posts"
}

// Category 分类
type Category struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:50;uniqueIndex" json:"name"`
	Slug      string    `gorm:"size:50;uniqueIndex" json:"slug"`
	Description string  `gorm:"size:200" json:"description"`
	ParentID  *uint64   `gorm:"index" json:"parent_id"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	PostCount int       `gorm:"default:0" json:"post_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 表名
func (Category) TableName() string {
	return "categories"
}

// SensitiveWord 敏感词
type SensitiveWord struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Word      string    `gorm:"size:100;uniqueIndex" json:"word"`
	Category  string    `gorm:"size:50;default:default" json:"category"`
	Level     int8      `gorm:"default:1;index" json:"level"`
	Status    int8      `gorm:"default:1;index" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 表名
func (SensitiveWord) TableName() string {
	return "sensitive_words"
}
