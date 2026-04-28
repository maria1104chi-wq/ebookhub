package service

import (
	"blog-system/internal/repository"
	"strings"
)

// SensitiveService 敏感词服务
type SensitiveService struct {
	repo *repository.SensitiveRepository
}

// NewSensitiveService 创建敏感词服务
func NewSensitiveService() *SensitiveService {
	return &SensitiveService{
		repo: &repository.SensitiveRepository{},
	}
}

// ImportFromFile 从文件导入敏感词
func (s *SensitiveService) ImportFromFile(filePath string) (int, error) {
	return s.repo.ImportSensitiveWords(filePath)
}

// FilterContent 过滤内容中的敏感词
func (s *SensitiveService) FilterContent(content string) string {
	var words []string
	words, _ = s.repo.GetSensitiveWords()

	result := content
	for _, word := range words {
		result = strings.ReplaceAll(result, word, "***")
	}
	return result
}

// CheckContent 检查内容是否包含敏感词
func (s *SensitiveService) CheckContent(content string) (bool, string) {
	words, _ := s.repo.GetSensitiveWords()

	for _, word := range words {
		if strings.Contains(content, word) {
			return true, word
		}
	}
	return false, ""
}

// AddWord 添加敏感词
func (s *SensitiveService) AddWord(word, category string, level int8) error {
	return s.repo.AddSensitiveWord(word, category, level)
}

// DeleteWord 删除敏感词
func (s *SensitiveService) DeleteWord(id uint64) error {
	return s.repo.DeleteSensitiveWord(id)
}

// UpdateStatus 更新敏感词状态
func (s *SensitiveService) UpdateStatus(id uint64, status int8) error {
	return s.repo.UpdateSensitiveWordStatus(id, status)
}

// GetList 获取敏感词列表
func (s *SensitiveService) GetList(page, pageSize int) ([]interface{}, int64, error) {
	words, total, err := s.repo.GetSensitiveWordList(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	result := make([]interface{}, len(words))
	for i, w := range words {
		result[i] = w
	}
	return result, total, nil
}

// AutoOptimizeSEO SEO 自动优化
func (s *SensitiveService) AutoOptimizeSEO(title, content string) (seoTitle, seoKeywords, seoDescription string) {
	// 清理敏感词
	cleanTitle := s.FilterContent(title)
	cleanContent := s.FilterContent(content)

	// 生成 SEO 标题（限制 60 字符）
	if len(cleanTitle) > 60 {
		seoTitle = cleanTitle[:60]
	} else {
		seoTitle = cleanTitle
	}

	// 提取关键词（简化版：取前 10 个名词）
	seoKeywords = extractKeywords(cleanContent, 10)

	// 生成描述（取前 200 字符）
	if len(cleanContent) > 200 {
		seoDescription = cleanContent[:200] + "..."
	} else {
		seoDescription = cleanContent
	}

	return
}

// extractKeywords 提取关键词（简化实现）
func extractKeywords(content string, limit int) string {
	// 这里可以使用更复杂的 NLP 算法
	// 简化版：返回固定关键词
	return "博客，技术，分享"
}
