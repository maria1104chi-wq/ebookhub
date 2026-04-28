package handler

import (
	"net/http"
	"time"

	"blog-system/internal/model"
	"blog-system/internal/repository"
	"blog-system/internal/service"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService  *service.AuthService
	smsService   *service.SMSService
	userRepo     *service.UserRepository
	sensitiveSvc *service.SensitiveService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService:  service.NewAuthService("your-jwt-secret-key-change-in-production"),
		smsService:   nil, // 初始化时设置
		userRepo:     &service.UserRepository{},
		sensitiveSvc: service.NewSensitiveService(),
	}
}

// SetSMSService 设置短信服务
func (h *AuthHandler) SetSMSService(svc *service.SMSService) {
	h.smsService = svc
}

// SendSMSCode 发送短信验证码
func (h *AuthHandler) SendSMSCode(c *gin.Context) {
	type Request struct {
		Phone string `json:"phone" binding:"required,len=11"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查手机号格式
	if len(req.Phone) != 11 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "手机号格式错误"})
		return
	}

	// 发送验证码
	code, err := h.smsService.SendVerifyCode(req.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发送失败：" + err.Error()})
		return
	}

	// 生产环境应该只返回成功，不返回验证码
	_ = code

	c.JSON(http.StatusOK, gin.H{"message": "验证码已发送"})
}

// Login 管理员登录（需要短信验证码）
func (h *AuthHandler) Login(c *gin.Context) {
	type Request struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Phone    string `json:"phone" binding:"required"`
		Code     string `json:"code" binding:"required,len=6"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证短信验证码
	valid, err := h.smsService.VerifyCode(req.Phone, req.Code)
	if err != nil || !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码错误"})
		return
	}

	// 查询用户
	user, err := h.userRepo.GetUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 检查状态
	if user.Status != 1 {
		c.JSON(http.StatusForbidden, gin.H{"error": "账号已被禁用"})
		return
	}

	// 更新登录信息
	ip := c.ClientIP()
	h.userRepo.UpdateLoginInfo(user.ID, ip)

	// 生成 Token
	token, err := h.authService.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

// GetUserInfo 获取当前用户信息
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	user, err := h.userRepo.GetUserByID(userID.(uint64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"avatar":   user.Avatar,
		"role":     user.Role,
	})
}

// Logout 登出
func (h *AuthHandler) Logout(c *gin.Context) {
	// 客户端删除 token 即可，服务端可加入黑名单
	c.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}

// RegisterAdmin 注册管理员（仅首次初始化使用）
func (h *AuthHandler) RegisterAdmin(c *gin.Context) {
	type Request struct {
		Username string `json:"username" binding:"required,min=3,max=20"`
		Password string `json:"password" binding:"required,min=6"`
		Phone    string `json:"phone" binding:"required,len=11"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查是否已有管理员
	var count int64
	repository.DB.Model(&model.User{}).Where("role = ?", 1).Count(&count)
	if count > 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "管理员已存在"})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	user := &model.User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		Phone:        req.Phone,
		Role:         1,
		Status:       1,
	}

	if err := repository.DB.Create(user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}
