package repository

import (
	"blog-system/internal/model"
	"context"
	"time"
)

// SMSRepository 短信仓库
type SMSRepository struct{}

// SaveCode 保存验证码
func (r *SMSRepository) SaveCode(phone, code, purpose string, expiresAt time.Time) error {
	smsCode := &model.SMSCode{
		Phone:     phone,
		Code:      code,
		Purpose:   purpose,
		ExpiresAt: expiresAt,
		Used:      0,
	}
	return DB.Create(smsCode).Error
}

// VerifyCode 验证验证码
func (r *SMSRepository) VerifyCode(phone, code string) (bool, error) {
	var smsCode model.SMSCode
	err := DB.Where("phone = ? AND code = ? AND used = ? AND expires_at > ?", 
		phone, code, 0, time.Now()).First(&smsCode).Error
	if err != nil {
		return false, err
	}

	// 标记为已使用
	DB.Model(&smsCode).Update("used", 1)
	return true, nil
}

// CacheSet 设置缓存
func CacheSet(key string, value interface{}, expiration time.Duration) error {
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

// CacheGet 获取缓存
func CacheGet(key string) (string, error) {
	return RedisClient.Get(ctx, key).Result()
}

// CacheDel 删除缓存
func CacheDel(key string) error {
	return RedisClient.Del(ctx, key).Err()
}

// CacheExists 检查缓存是否存在
func CacheExists(key string) (bool, error) {
	result, err := RedisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}
