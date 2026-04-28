package repository

import (
	"blog-system/internal/model"
	"bufio"
	"os"
)

// SensitiveRepository 敏感词仓库
type SensitiveRepository struct{}

// ImportSensitiveWords 从文件导入敏感词
func (r *SensitiveRepository) ImportSensitiveWords(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var count int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := scanner.Text()
		if word == "" {
			continue
		}

		sensitiveWord := &model.SensitiveWord{
			Word:     word,
			Category: "default",
			Level:    1,
			Status:   1,
		}

		// 忽略已存在的敏感词
		DB.Where("word = ?", word).FirstOrCreate(sensitiveWord)
		count++
	}

	return count, scanner.Err()
}

// GetSensitiveWords 获取所有启用的敏感词
func (r *SensitiveRepository) GetSensitiveWords() ([]string, error) {
	var words []model.SensitiveWord
	err := DB.Where("status = ?", 1).Find(&words).Error
	if err != nil {
		return nil, err
	}

	result := make([]string, len(words))
	for i, w := range words {
		result[i] = w.Word
	}
	return result, nil
}

// AddSensitiveWord 添加敏感词
func (r *SensitiveRepository) AddSensitiveWord(word, category string, level int8) error {
	sensitiveWord := &model.SensitiveWord{
		Word:     word,
		Category: category,
		Level:    level,
		Status:   1,
	}
	return DB.Create(sensitiveWord).Error
}

// DeleteSensitiveWord 删除敏感词
func (r *SensitiveRepository) DeleteSensitiveWord(id uint64) error {
	return DB.Delete(&model.SensitiveWord{}, id).Error
}

// UpdateSensitiveWordStatus 更新敏感词状态
func (r *SensitiveRepository) UpdateSensitiveWordStatus(id uint64, status int8) error {
	return DB.Model(&model.SensitiveWord{}).Where("id = ?", id).Update("status", status).Error
}

// GetSensitiveWordList 获取敏感词列表
func (r *SensitiveRepository) GetSensitiveWordList(page, pageSize int) ([]*model.SensitiveWord, int64, error) {
	var words []*model.SensitiveWord
	var total int64

	query := DB.Model(&model.SensitiveWord{})
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("id DESC").Limit(pageSize).Offset(offset).Find(&words).Error
	return words, total, err
}

// CheckSensitiveWord 检查是否包含敏感词
func (r *SensitiveRepository) CheckSensitiveWord(content string) (bool, string) {
	var words []model.SensitiveWord
	DB.Where("status = ?", 1).Find(&words)

	for _, w := range words {
		if contains(content, w.Word) {
			return true, w.Word
		}
	}
	return false, ""
}

// contains 检查字符串是否包含子串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		findSubstring(s, substr))
}

// findSubstring 查找子串
func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// FilterSensitiveWords 过滤敏感词
func (r *SensitiveRepository) FilterSensitiveWords(content string) string {
	var words []model.SensitiveWord
	DB.Where("status = ?", 1).Find(&words)

	result := content
	for _, w := range words {
		result = replaceAll(result, w.Word, "***")
	}
	return result
}

// replaceAll 替换所有出现
func replaceAll(s, old, new string) string {
	result := ""
	for i := 0; i <= len(s)-len(old); {
		if s[i:i+len(old)] == old {
			result += new
			i += len(old)
		} else {
			result += string(s[i])
			i++
		}
	}
	result += s[len(s)-(len(s)-len(result)+len(new)*(len(s)/len(old))):]
	return result
}
