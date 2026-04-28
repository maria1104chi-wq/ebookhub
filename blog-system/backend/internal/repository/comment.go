package repository

import (
	"blog-system/internal/model"
)

// CommentRepository 评论仓库
type CommentRepository struct{}

// CreateComment 创建评论
func (r *CommentRepository) CreateComment(comment *model.Comment) error {
	return DB.Create(comment).Error
}

// UpdateComment 更新评论
func (r *CommentRepository) UpdateComment(comment *model.Comment) error {
	return DB.Save(comment).Error
}

// DeleteComment 删除评论
func (r *CommentRepository) DeleteComment(id uint64) error {
	return DB.Delete(&model.Comment{}, id).Error
}

// GetCommentByID 根据 ID 获取评论
func (r *CommentRepository) GetCommentByID(id uint64) (*model.Comment, error) {
	var comment model.Comment
	err := DB.First(&comment, id).Error
	return &comment, err
}

// GetCommentsByPostID 根据文章 ID 获取评论列表
func (r *CommentRepository) GetCommentsByPostID(postID uint64, page, pageSize int) ([]*model.Comment, int64, error) {
	var comments []*model.Comment
	var total int64

	query := DB.Model(&model.Comment{}).Where("post_id = ? AND status = ?", postID, 2)
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("created_at ASC").Limit(pageSize).Offset(offset).Find(&comments).Error
	return comments, total, err
}

// GetChildComments 获取子评论
func (r *CommentRepository) GetChildComments(parentID uint64) ([]*model.Comment, error) {
	var comments []*model.Comment
	err := DB.Where("parent_id = ? AND status = ?", parentID, 2).Order("created_at ASC").Find(&comments).Error
	return comments, err
}

// GetCommentTree 获取评论树
func (r *CommentRepository) GetCommentTree(postID uint64) ([]map[string]interface{}, error) {
	var comments []*model.Comment
	DB.Where("post_id = ? AND parent_id IS NULL AND status = ?", postID, 2).
		Order("created_at ASC").Find(&comments)

	result := make([]map[string]interface{}, 0)
	for _, comment := range comments {
		childComments, _ := r.GetChildComments(comment.ID)
		
		commentMap := map[string]interface{}{
			"id":           comment.ID,
			"user_name":    comment.UserName,
			"content":      comment.Content,
			"user_location": comment.UserLocation,
			"like_count":   comment.LikeCount,
			"is_admin":     comment.IsAdmin,
			"created_at":   comment.CreatedAt,
			"children":     childComments,
		}
		result = append(result, commentMap)
	}

	return result, nil
}

// IncrementCommentLikeCount 增加评论点赞数
func (r *CommentRepository) IncrementCommentLikeCount(id uint64) error {
	return DB.Model(&model.Comment{}).Where("id = ?", id).Update("like_count", DB.Raw("like_count + 1")).Error
}

// CountCommentsByPostID 统计文章评论数
func (r *CommentRepository) CountCommentsByPostID(postID uint64) (int64, error) {
	var count int64
	err := DB.Model(&model.Comment{}).Where("post_id = ? AND status = ?", postID, 2).Count(&count).Error
	return count, err
}

// GetCommentsByUserID 根据用户 ID 获取评论
func (r *CommentRepository) GetCommentsByUserID(userID uint64, page, pageSize int) ([]*model.Comment, int64, error) {
	var comments []*model.Comment
	var total int64

	query := DB.Model(&model.Comment{}).Where("user_id = ?", userID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&comments).Error
	return comments, total, err
}

// GetRecentComments 获取最新评论
func (r *CommentRepository) GetRecentComments(limit int) ([]*model.Comment, error) {
	var comments []*model.Comment
	err := DB.Where("status = ?", 2).
		Order("created_at DESC").
		Limit(limit).
		Find(&comments).Error
	return comments, err
}

// UpdateCommentStatus 更新评论状态
func (r *CommentRepository) UpdateCommentStatus(id uint64, status int8) error {
	return DB.Model(&model.Comment{}).Where("id = ?", id).Update("status", status).Error
}

// GetPendingComments 获取待审核评论
func (r *CommentRepository) GetPendingComments(page, pageSize int) ([]*model.Comment, int64, error) {
	var comments []*model.Comment
	var total int64

	query := DB.Model(&model.Comment{}).Where("status = ?", 1)
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&comments).Error
	return comments, total, err
}
