package repository

import (
	"blog-system/internal/model"
	"time"
)

// PostRepository 博文仓库
type PostRepository struct{}

// CreatePost 创建博文
func (r *PostRepository) CreatePost(post *model.Post) error {
	return DB.Create(post).Error
}

// UpdatePost 更新博文
func (r *PostRepository) UpdatePost(post *model.Post) error {
	return DB.Save(post).Error
}

// DeletePost 删除博文
func (r *PostRepository) DeletePost(id uint64) error {
	return DB.Delete(&model.Post{}, id).Error
}

// GetPostByID 根据 ID 获取博文
func (r *PostRepository) GetPostByID(id uint64) (*model.Post, error) {
	var post model.Post
	err := DB.First(&post, id).Error
	return &post, err
}

// GetPostBySlug 根据 slug 获取博文
func (r *PostRepository) GetPostBySlug(slug string) (*model.Post, error) {
	var post model.Post
	err := DB.Where("slug = ?", slug).First(&post).Error
	return &post, err
}

// GetPostList 获取博文列表
func (r *PostRepository) GetPostList(page, pageSize int, status int8) ([]*model.Post, int64, error) {
	var posts []*model.Post
	var total int64

	query := DB.Model(&model.Post{}).Where("status = ?", status)
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("is_top DESC, published_at DESC").
		Limit(pageSize).Offset(offset).Find(&posts).Error

	return posts, total, err
}

// GetHotPosts 获取热门博文（点击率前十）
func (r *PostRepository) GetHotPosts(limit int) ([]*model.Post, error) {
	var posts []*model.Post
	err := DB.Where("status = ?", 1).
		Order("view_count DESC").
		Limit(limit).
		Find(&posts).Error
	return posts, err
}

// SearchPosts 搜索博文
func (r *PostRepository) SearchPosts(keyword string, page, pageSize int) ([]*model.Post, int64, error) {
	var posts []*model.Post
	var total int64

	searchTerm := "%" + keyword + "%"
	query := DB.Model(&model.Post{}).
		Where("status = ? AND (title LIKE ? OR summary LIKE ?)", 1, searchTerm, searchTerm)
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("published_at DESC").
		Limit(pageSize).Offset(offset).Find(&posts).Error

	return posts, total, err
}

// IncrementViewCount 增加浏览次数
func (r *PostRepository) IncrementViewCount(id uint64) error {
	return DB.Model(&model.Post{}).Where("id = ?", id).Update("view_count", DB.Raw("view_count + 1")).Error
}

// IncrementLikeCount 增加点赞数
func (r *PostRepository) IncrementLikeCount(id uint64) error {
	return DB.Model(&model.Post{}).Where("id = ?", id).Update("like_count", DB.Raw("like_count + 1")).Error
}

// IncrementShareCount 增加分享数
func (r *PostRepository) IncrementShareCount(id uint64) error {
	return DB.Model(&model.Post{}).Where("id = ?", id).Update("share_count", DB.Raw("share_count + 1")).Error
}

// UpdateCommentCount 更新评论数
func (r *PostRepository) UpdateCommentCount(id uint64, count int) error {
	return DB.Model(&model.Post{}).Where("id = ?", id).Update("comment_count", count).Error
}

// GetPostsByCategory 根据分类获取博文
func (r *PostRepository) GetPostsByCategory(categoryID uint64, page, pageSize int) ([]*model.Post, int64, error) {
	var posts []*model.Post
	var total int64

	query := DB.Model(&model.Post{}).Where("status = ? AND category_id = ?", 1, categoryID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("published_at DESC").
		Limit(pageSize).Offset(offset).Find(&posts).Error

	return posts, total, err
}

// GetPostsByTag 根据标签获取博文
func (r *PostRepository) GetPostsByTag(tag string, page, pageSize int) ([]*model.Post, int64, error) {
	var posts []*model.Post
	var total int64

	query := DB.Model(&model.Post{}).
		Where("status = ? AND tags LIKE ?", 1, "%"+tag+"%")
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("published_at DESC").
		Limit(pageSize).Offset(offset).Find(&posts).Error

	return posts, total, err
}

// GetRecentPosts 获取最新博文
func (r *PostRepository) GetRecentPosts(limit int) ([]*model.Post, error) {
	var posts []*model.Post
	err := DB.Where("status = ?", 1).
		Order("published_at DESC").
		Limit(limit).
		Find(&posts).Error
	return posts, err
}

// GetArchives 获取归档列表
func (r *PostRepository) GetArchives() ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := DB.Model(&model.Post{}).
		Select("DATE_FORMAT(published_at, '%Y-%m') as month, COUNT(*) as count").
		Where("status = ?", 1).
		Group("month").
		Order("month DESC").
		Scan(&results).Error
	return results, err
}

// CountPostsByDate 统计指定日期的文章数
func (r *PostRepository) CountPostsByDate(date time.Time) (int64, error) {
	var count int64
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	err := DB.Model(&model.Post{}).
		Where("status = ? AND created_at BETWEEN ? AND ?", 1, startOfDay, endOfDay).
		Count(&count).Error
	return count, err
}
