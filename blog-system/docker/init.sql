-- 数据库初始化脚本

-- 创建数据库
CREATE DATABASE IF NOT EXISTS blog_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE blog_db;

-- 这里放置额外的初始化 SQL
-- 表结构由 GORM 自动迁移创建

-- 插入默认管理员（密码：admin123）
-- 注意：实际密码哈希应在应用中生成
INSERT INTO users (username, password_hash, phone, role, status, created_at, updated_at)
VALUES ('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '13800138000', 1, 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE username=username;

-- 插入默认分类
INSERT INTO categories (name, slug, description, sort_order, created_at, updated_at) VALUES
('技术', 'tech', '技术相关文章', 1, NOW(), NOW()),
('生活', 'life', '生活随笔', 2, NOW(), NOW()),
('随笔', 'notes', '随想笔记', 3, NOW(), NOW())
ON DUPLICATE KEY UPDATE name=name;

-- 插入默认配置
INSERT INTO settings (setting_key, setting_value, setting_type, description, created_at, updated_at) VALUES
('site_name', '我的博客', 'string', '网站名称', NOW(), NOW()),
('site_description', '个人博客系统', 'string', '网站描述', NOW(), NOW()),
('site_keywords', '博客，技术，生活', 'string', '网站关键词', NOW(), NOW()),
('posts_per_page', '10', 'string', '每页文章数', NOW(), NOW()),
('allow_comment', 'true', 'string', '允许评论', NOW(), NOW()),
('comment_audit', 'false', 'string', '评论需审核', NOW(), NOW())
ON DUPLICATE KEY UPDATE setting_key=setting_key;
