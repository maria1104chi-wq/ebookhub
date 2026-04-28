package service

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"blog-system/internal/config"
	"blog-system/internal/repository"

	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2"
)

// SMSService 短信服务
type SMSService struct {
	client *dysmsapi20170525.Client
	config *config.SMSConfig
}

// NewSMSService 创建短信服务
func NewSMSService(cfg *config.SMSConfig) (*SMSService, error) {
	conf := &openapi.Config{
		AccessKeyId:     &cfg.AccessKeyID,
		AccessKeySecret: &cfg.AccessKeySecret,
	}
	conf.Endpoint = StringPtr("dysmsapi.aliyuncs.com")

	client, err := dysmsapi20170525.NewClient(conf)
	if err != nil {
		return nil, err
	}

	return &SMSService{
		client: client,
		config: cfg,
	}, nil
}

// StringPtr 字符串指针辅助函数
func StringPtr(s string) *string {
	return &s
}

// SendVerifyCode 发送验证码
func (s *SMSService) SendVerifyCode(phone string) (string, error) {
	// 生成 6 位验证码
	code, err := generateCode()
	if err != nil {
		return "", err
	}

	// 调用阿里云短信 API
	request := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  StringPtr(phone),
		SignName:      StringPtr(s.config.SignName),
		TemplateCode:  StringPtr(s.config.TemplateCode),
		TemplateParam: StringPtr(fmt.Sprintf(`{"code":"%s"}`, code)),
	}

	_, err = s.client.SendSms(request)
	if err != nil {
		return "", fmt.Errorf("发送短信失败：%w", err)
	}

	// 保存验证码到数据库
	expiresAt := time.Now().Add(5 * time.Minute)
	smsRepo := &repository.SMSRepository{}
	err = smsRepo.SaveCode(phone, code, "login", expiresAt)
	if err != nil {
		return "", fmt.Errorf("保存验证码失败：%w", err)
	}

	return code, nil
}

// VerifyCode 验证验证码
func (s *SMSService) VerifyCode(phone, code string) (bool, error) {
	smsRepo := &repository.SMSRepository{}
	return smsRepo.VerifyCode(phone, code)
}

// generateCode 生成 6 位随机验证码
func generateCode() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}
