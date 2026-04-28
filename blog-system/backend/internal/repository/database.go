package repository

import (
	"fmt"
	"log"

	"blog-system/internal/config"
	"blog-system/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(cfg *config.DatabaseConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("连接数据库失败：%w", err)
	}

	// 自动迁移表结构
	err = DB.AutoMigrate(
		&model.User{},
		&model.SMSCode{},
		&model.Post{},
		&model.Category{},
		&model.SensitiveWord{},
		&model.Comment{},
		&model.Like{},
		&model.Upload{},
		&model.Statistics{},
		&model.Setting{},
		&model.OperationLog{},
	)
	if err != nil {
		return fmt.Errorf("数据库迁移失败：%w", err)
	}

	log.Println("数据库初始化成功")
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
