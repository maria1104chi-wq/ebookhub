package handler

import (
	"net/http"
	"strconv"

	"blog-system/internal/model"
	"blog-system/internal/repository"
	"blog-system/internal/service"

	"github.com/gin-gonic/gin"
)

// AdminHandler 后台管理处理器
type AdminHandler struct {
	sensitiveSvc *service.SensitiveService
}

// NewAdminHandler 创建后台管理处理器
func NewAdminHandler() *AdminHandler {
	return &AdminHandler{
		sensitiveSvc: service.NewSensitiveService(),
	}
}

// GetDashboard 获取仪表盘数据
func (h *AdminHandler) GetDashboard(c *gin.Context) {
	var postCount, commentCount, userCount int64
	repository.DB.Model(&model.Post{}).Where("status = ?", 1).Count(&postCount)
	repository.DB.Model(&model.Comment{}).Where("status = ?", 2).Count(&commentCount)
	repository.DB.Model(&model.User{}).Count(&userCount)

	// 获取今日数据
	todayStats := h.getTodayStats()

	c.JSON(http.StatusOK, gin.H{
		"post_count":    postCount,
		"comment_count": commentCount,
		"user_count":    userCount,
		"today":         todayStats,
	})
}

// getTodayStats 获取今日统计
func (h *AdminHandler) getTodayStats() map[string]interface{} {
	// 简化实现
	return map[string]interface{}{
		"views":   0,
		"posts":   0,
		"comments": 0,
	}
}

// GetSensitiveList 获取敏感词列表
func (h *AdminHandler) GetSensitiveList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	list, total, err := h.sensitiveSvc.GetList(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
		"page":  page,
	})
}

// AddSensitiveWord 添加敏感词
func (h *AdminHandler) AddSensitiveWord(c *gin.Context) {
	type Request struct {
		Word     string `json:"word" binding:"required"`
		Category string `json:"category"`
		Level    int8   `json:"level"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.sensitiveSvc.AddWord(req.Word, req.Category, req.Level); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "添加成功"})
}

// DeleteSensitiveWord 删除敏感词
func (h *AdminHandler) DeleteSensitiveWord(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.sensitiveSvc.DeleteWord(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// UpdateSensitiveStatus 更新敏感词状态
func (h *AdminHandler) UpdateSensitiveStatus(c *gin.Context) {
	type Request struct {
		ID     uint64 `json:"id" binding:"required"`
		Status int8   `json:"status" binding:"required"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.sensitiveSvc.UpdateStatus(req.ID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// ImportSensitiveWords 导入敏感词
func (h *AdminHandler) ImportSensitiveWords(c *gin.Context) {
	count, err := h.sensitiveSvc.ImportFromFile("./data/Sensitive.txt")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "导入失败：" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "导入成功", "count": count})
}

// GetStatistics 获取统计数据
func (h *AdminHandler) GetStatistics(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	var stats []*model.Statistics
	query := repository.DB.Model(&model.Statistics{})
	if startDate != "" && endDate != "" {
		query = query.Where("stat_date BETWEEN ? AND ?", startDate, endDate)
	}
	query.Order("stat_date DESC").Find(&stats)

	c.JSON(http.StatusOK, gin.H{"list": stats})
}

// BackupDatabase 备份数据库
func (h *AdminHandler) BackupDatabase(c *gin.Context) {
	// 实际应调用 mysqldump 命令
	c.JSON(http.StatusOK, gin.H{"message": "备份功能待实现"})
}
