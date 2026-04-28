package handler

import (
	"net/http"
	"strconv"
	"time"

	"blog-system/internal/model"
	"blog-system/internal/repository"
	"blog-system/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

// PostHandler 博文处理器
type PostHandler struct {
	postRepo     *repository.PostRepository
	sensitiveSvc *service.SensitiveService
}

// NewPostHandler 创建博文处理器
func NewPostHandler() *PostHandler {
	return &PostHandler{
		postRepo:     &repository.PostRepository{},
		sensitiveSvc: service.NewSensitiveService(),
	}
}

// CreatePost 创建博文
func (h *PostHandler) CreatePost(c *gin.Context) {
	type Request struct {
		Title      string   `json:"title" binding:"required"`
		Summary    string   `json:"summary"`
		Content    string   `json:"content" binding:"required"`
		CoverImage string   `json:"cover_image"`
		CategoryID uint64   `json:"category_id"`
		Tags       []string `json:"tags"`
		Status     int8     `json:"status"`
		IsTop      int8     `json:"is_top"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查敏感词
	hasSensitive, word := h.sensitiveSvc.CheckContent(req.Title + req.Content)
	if hasSensitive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "包含敏感词：" + word})
		return
	}

	// Markdown 转 HTML
	md := goldmark.New(goldmark.WithExtensions(extension.GFM))
	var htmlBuf []byte
	err := md.Convert([]byte(req.Content), &htmlBuf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Markdown 转换失败"})
		return
	}

	// SEO 优化
	seoTitle, seoKeywords, seoDescription := h.sensitiveSvc.AutoOptimizeSEO(req.Title, req.Content)

	// 生成 slug
	slug := generateSlug(req.Title)

	now := time.Now()
	post := &model.Post{
		Title:          req.Title,
		Slug:           slug,
		Summary:        req.Summary,
		Content:        req.Content,
		ContentHTML:    string(htmlBuf),
		CoverImage:     req.CoverImage,
		AuthorID:       c.GetUint64("user_id"),
		CategoryID:     req.CategoryID,
		Tags:           arrayToJSON(req.Tags),
		Status:         req.Status,
		IsTop:          req.IsTop,
		SeoTitle:       seoTitle,
		SeoKeywords:    seoKeywords,
		SeoDescription: seoDescription,
		PublishedAt:    &now,
	}

	if err := h.postRepo.CreatePost(post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": post.ID, "message": "创建成功"})
}

// UpdatePost 更新博文
func (h *PostHandler) UpdatePost(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	post, err := h.postRepo.GetPostByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	type Request struct {
		Title      string   `json:"title"`
		Summary    string   `json:"summary"`
		Content    string   `json:"content"`
		CoverImage string   `json:"cover_image"`
		CategoryID uint64   `json:"category_id"`
		Tags       []string `json:"tags"`
		Status     int8     `json:"status"`
		IsTop      int8     `json:"is_top"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新字段
	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Summary != "" {
		post.Summary = req.Summary
	}
	if req.Content != "" {
		// 检查敏感词
		hasSensitive, word := h.sensitiveSvc.CheckContent(req.Content)
		if hasSensitive {
			c.JSON(http.StatusBadRequest, gin.H{"error": "包含敏感词：" + word})
			return
		}

		post.Content = req.Content
		md := goldmark.New(goldmark.WithExtensions(extension.GFM))
		var htmlBuf []byte
		md.Convert([]byte(req.Content), &htmlBuf)
		post.ContentHTML = string(htmlBuf)
	}
	if req.CoverImage != "" {
		post.CoverImage = req.CoverImage
	}
	if req.CategoryID > 0 {
		post.CategoryID = req.CategoryID
	}
	if len(req.Tags) > 0 {
		post.Tags = arrayToJSON(req.Tags)
	}
	post.Status = req.Status
	post.IsTop = req.IsTop

	if err := h.postRepo.UpdatePost(post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// DeletePost 删除博文
func (h *PostHandler) DeletePost(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.postRepo.DeletePost(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetPost 获取博文详情
func (h *PostHandler) GetPost(c *gin.Context) {
	slug := c.Param("slug")

	post, err := h.postRepo.GetPostBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 增加浏览次数
	h.postRepo.IncrementViewCount(post.ID)

	c.JSON(http.StatusOK, post)
}

// GetPostList 获取博文列表
func (h *PostHandler) GetPostList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	posts, total, err := h.postRepo.GetPostList(page, pageSize, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"list":  posts,
		"total": total,
		"page":  page,
	})
}

// SearchPosts 搜索博文
func (h *PostHandler) SearchPosts(c *gin.Context) {
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请输入搜索关键词"})
		return
	}

	posts, total, err := h.postRepo.SearchPosts(keyword, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"list":  posts,
		"total": total,
		"page":  page,
	})
}

// GetHotPosts 获取热门博文
func (h *PostHandler) GetHotPosts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	posts, err := h.postRepo.GetHotPosts(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"list": posts})
}

// LikePost 点赞博文
func (h *PostHandler) LikePost(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.postRepo.IncrementLikeCount(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "点赞失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "点赞成功"})
}

// SharePost 分享博文
func (h *PostHandler) SharePost(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.postRepo.IncrementShareCount(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "分享失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "分享成功"})
}

// generateSlug 生成 URL 别名
func generateSlug(title string) string {
	// 简化实现，实际可用拼音库
	return title
}

// arrayToJSON 数组转 JSON 字符串
func arrayToJSON(arr []string) string {
	if len(arr) == 0 {
		return "[]"
	}
	result := "["
	for i, s := range arr {
		if i > 0 {
			result += ","
		}
		result += "\"" + s + "\""
	}
	result += "]"
	return result
}
