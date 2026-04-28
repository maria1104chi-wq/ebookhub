package config

import (
	"github.com/spf13/viper"
)

// Config 配置结构体
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	SMS      SMSConfig
	Upload   UploadConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string
	Mode string
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// SMSConfig 短信配置
type SMSConfig struct {
	AccessKeyID     string
	AccessKeySecret string
	SignName        string
	TemplateCode    string
}

// UploadConfig 上传配置
type UploadConfig struct {
	MaxImageSize int64
	MaxPDFSize   int64
	MaxVideoSize int64
	ImagePath    string
	PDFPath      string
	VideoPath    string
}

// Load 加载配置
func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := &Config{
		Server: ServerConfig{
			Port: viper.GetString("SERVER_PORT"),
			Mode: viper.GetString("SERVER_MODE"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			DBName:   viper.GetString("DB_NAME"),
		},
		Redis: RedisConfig{
			Host:     viper.GetString("REDIS_HOST"),
			Port:     viper.GetString("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
		SMS: SMSConfig{
			AccessKeyID:     viper.GetString("ALIYUN_ACCESS_KEY_ID"),
			AccessKeySecret: viper.GetString("ALIYUN_ACCESS_KEY_SECRET"),
			SignName:        viper.GetString("ALIYUN_SMS_SIGN"),
			TemplateCode:    viper.GetString("ALIYUN_SMS_TEMPLATE"),
		},
		Upload: UploadConfig{
			MaxImageSize: viper.GetInt64("UPLOAD_MAX_IMAGE_SIZE"),
			MaxPDFSize:   viper.GetInt64("UPLOAD_MAX_PDF_SIZE"),
			MaxVideoSize: viper.GetInt64("UPLOAD_MAX_VIDEO_SIZE"),
			ImagePath:    viper.GetString("UPLOAD_IMAGE_PATH"),
			PDFPath:      viper.GetString("UPLOAD_PDF_PATH"),
			VideoPath:    viper.GetString("UPLOAD_VIDEO_PATH"),
		},
	}

	return config, nil
}
