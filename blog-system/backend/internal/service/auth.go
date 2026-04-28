package service

import (
	"blog-system/internal/model"
	"blog-system/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AuthService 认证服务
type AuthService struct {
	jwtSecret string
}

// NewAuthService 创建认证服务
func NewAuthService(jwtSecret string) *AuthService {
	return &AuthService{
		jwtSecret: jwtSecret,
	}
}

// Claims JWT 声明
type Claims struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	Role     int8   `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT 令牌
func (s *AuthService) GenerateToken(user *model.User) (string, error) {
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// ParseToken 解析 JWT 令牌
func (s *AuthService) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// ErrInvalidToken 无效令牌错误
var ErrInvalidToken = &jwt.ValidationError{Errors: jwt.ValidationErrorSignatureInvalid}

// UserRepository 用户仓库操作
type UserRepository struct{}

// GetUserByUsername 根据用户名获取用户
func (r *UserRepository) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := repository.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

// GetUserByID 根据 ID 获取用户
func (r *UserRepository) GetUserByID(id uint64) (*model.User, error) {
	var user model.User
	err := repository.DB.First(&user, id).Error
	return &user, err
}

// UpdateLoginInfo 更新登录信息
func (r *UserRepository) UpdateLoginInfo(id uint64, ip string) error {
	now := time.Now()
	return repository.DB.Model(&model.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"last_login_at": now,
		"last_login_ip": ip,
	}).Error
}
