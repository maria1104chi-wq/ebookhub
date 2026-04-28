package handler

import (
	"net/http"
	"strconv"

	"blog-system/internal/model"
	"blog-system/internal/repository"
	"blog-system/internal/service"

	"github.com/gin-gonic/gin"
)

// CommentHandler 评论处理器
type CommentHandler struct {
	commentRepo  *repository.CommentRepository
	sensitiveSvc *service.SensitiveService
}

// NewCommentHandler 创建评论处理器
func NewCommentHandler() *CommentHandler {
	return &CommentHandler{
		commentRepo:  &repository.CommentRepository{},
		sensitiveSvc: service.NewSensitiveService(),
	}
}

// CreateComment 创建评论
func (h *CommentHandler) CreateComment(c *gin.Context) {
	type Request struct {
		PostID    uint64 `json:"post_id" binding:"required"`
		ParentID  uint64 `json:"parent_id"`
		UserName  string `json:"user_name" binding:"required"`
		UserEmail string `json:"user_email"`
		Content   string `json:"content" binding:"required"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查敏感词
	hasSensitive, word := h.sensitiveSvc.CheckContent(req.Content)
	if hasSensitive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "包含敏感词：" + word})
		return
	}

	// 过滤敏感词
	cleanContent := h.sensitiveSvc.FilterContent(req.Content)

	// 获取 IP 和地理位置
	ip := c.ClientIP()
	location := getIPLocation(ip)

	comment := &model.Comment{
		PostID:       req.PostID,
		ParentID:     &req.ParentID,
		UserName:     req.UserName,
		UserEmail:    req.UserEmail,
		UserIP:       ip,
		UserLocation: location,
		UserAgent:    c.Request.UserAgent(),
		Content:      cleanContent,
		Status:       2, // 默认通过
		IsAdmin:      0,
	}

	if err := h.commentRepo.CreateComment(comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "评论失败"})
		return
	}

	// 更新文章评论数
	count, _ := h.commentRepo.CountCommentsByPostID(req.PostID)
	h.updatePostCommentCount(req.PostID, int(count))

	c.JSON(http.StatusOK, gin.H{"id": comment.ID, "message": "评论成功"})
}

// GetComments 获取评论列表
func (h *CommentHandler) GetComments(c *gin.Context) {
	postID, _ := strconv.ParseUint(c.Param("post_id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	tree, err := h.commentRepo.GetCommentTree(postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"list": tree})
}

// LikeComment 点赞评论
func (h *CommentHandler) LikeComment(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.commentRepo.IncrementCommentLikeCount(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "点赞失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "点赞成功"})
}

// DeleteComment 删除评论
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.commentRepo.DeleteComment(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// AuditComment 审核评论
func (h *CommentHandler) AuditComment(c *gin.Context) {
	type Request struct {
		ID     uint64 `json:"id" binding:"required"`
		Status int8   `json:"status" binding:"required,oneof=2 3"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.commentRepo.UpdateCommentStatus(req.ID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "审核失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "审核成功"})
}

// GetPendingComments 获取待审核评论
func (h *CommentHandler) GetPendingComments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	comments, total, err := h.commentRepo.GetPendingComments(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"list":  comments,
		"total": total,
		"page":  page,
	})
}

// updatePostCommentCount 更新文章评论数
func (h *CommentHandler) updatePostCommentCount(postID uint64, count int) {
	postRepo := &repository.PostRepository{}
	postRepo.UpdateCommentCount(postID, count)
}

// getIPLocation 获取 IP 地理位置（简化实现）
func getIPLocation(ip string) string {
	// 实际可接入 IP 地理位置 API
	return "未知地区"
}
