package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"blog-system/internal/model"
	"blog-system/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UploadHandler 上传处理器
type UploadHandler struct{}

// NewUploadHandler 创建上传处理器
func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

// UploadImage 上传图片
func (h *UploadHandler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未选择文件"})
		return
	}

	// 检查文件大小（最大 5MB）
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小不能超过 5MB"})
		return
	}

	// 检查文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	if !containsString(allowedExts, ext) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的图片格式"})
		return
	}

	// 生成唯一文件名
	filename := uuid.New().String() + ext
	savePath := "./uploads/images/" + filename

	// 保存文件
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传失败"})
		return
	}

	// 记录到数据库
	upload := &model.Upload{
		Filename:     filename,
		OriginalName: file.Filename,
		FilePath:     savePath,
		FileURL:      "/static/images/" + filename,
		FileType:     "image",
		FileSize:     file.Size,
		MimeType:     file.Header.Get("Content-Type"),
		UploaderID:   c.GetUint64("user_id"),
		Status:       1,
	}

	if err := repository.DB.Create(upload).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url":      upload.FileURL,
		"filename": filename,
	})
}

// UploadPDF 上传 PDF
func (h *UploadHandler) UploadPDF(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未选择文件"})
		return
	}

	// 检查文件大小（最大 10MB）
	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小不能超过 10MB"})
		return
	}

	// 检查文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".pdf" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只支持 PDF 格式"})
		return
	}

	filename := uuid.New().String() + ext
	savePath := "./uploads/pdf/" + filename

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传失败"})
		return
	}

	upload := &model.Upload{
		Filename:     filename,
		OriginalName: file.Filename,
		FilePath:     savePath,
		FileURL:      "/static/pdf/" + filename,
		FileType:     "pdf",
		FileSize:     file.Size,
		MimeType:     file.Header.Get("Content-Type"),
		UploaderID:   c.GetUint64("user_id"),
		Status:       1,
	}

	repository.DB.Create(upload)

	c.JSON(http.StatusOK, gin.H{
		"url":      upload.FileURL,
		"filename": filename,
	})
}

// UploadVideo 上传视频
func (h *UploadHandler) UploadVideo(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未选择文件"})
		return
	}

	// 检查文件大小（最大 100MB）
	if file.Size > 100*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小不能超过 100MB"})
		return
	}

	// 检查文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := []string{".mp4", ".webm", ".ogg", ".mov"}
	if !containsString(allowedExts, ext) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的视频格式"})
		return
	}

	filename := uuid.New().String() + ext
	savePath := "./uploads/video/" + filename

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传失败"})
		return
	}

	upload := &model.Upload{
		Filename:     filename,
		OriginalName: file.Filename,
		FilePath:     savePath,
		FileURL:      "/static/video/" + filename,
		FileType:     "video",
		FileSize:     file.Size,
		MimeType:     file.Header.Get("Content-Type"),
		UploaderID:   c.GetUint64("user_id"),
		Status:       1,
	}

	repository.DB.Create(upload)

	c.JSON(http.StatusOK, gin.H{
		"url":      upload.FileURL,
		"filename": filename,
	})
}

// GetUploadList 获取上传文件列表
func (h *UploadHandler) GetUploadList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	fileType := c.Query("type")

	var uploads []*model.Upload
	var total int64

	query := repository.DB.Model(&model.Upload{})
	if fileType != "" {
		query = query.Where("file_type = ?", fileType)
	}
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&uploads)

	c.JSON(http.StatusOK, gin.H{
		"list":  uploads,
		"total": total,
		"page":  page,
	})
}

// DeleteUpload 删除上传文件
func (h *UploadHandler) DeleteUpload(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := repository.DB.Delete(&model.Upload{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// containsString 检查字符串切片是否包含某字符串
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}
