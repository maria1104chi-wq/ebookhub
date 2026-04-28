# 数据库设计文档

## 数据库结构 SQL

```sql
-- 创建数据库
CREATE DATABASE IF NOT EXISTS blog_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE blog_db;

-- 用户表（管理员）
CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希',
    phone VARCHAR(20) NOT NULL UNIQUE COMMENT '手机号',
    email VARCHAR(100) COMMENT '邮箱',
    avatar VARCHAR(255) COMMENT '头像 URL',
    role TINYINT DEFAULT 1 COMMENT '角色：1 管理员',
    status TINYINT DEFAULT 1 COMMENT '状态：1 正常 0 禁用',
    last_login_at DATETIME COMMENT '最后登录时间',
    last_login_ip VARCHAR(45) COMMENT '最后登录 IP',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_username (username),
    INDEX idx_phone (phone),
    INDEX idx_status (status)
) ENGINE=InnoDB COMMENT='管理员用户表';

-- 短信验证码表
CREATE TABLE sms_codes (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    phone VARCHAR(20) NOT NULL COMMENT '手机号',
    code VARCHAR(6) NOT NULL COMMENT '验证码',
    purpose VARCHAR(20) NOT NULL COMMENT '用途：login/register',
    expires_at DATETIME NOT NULL COMMENT '过期时间',
    used TINYINT DEFAULT 0 COMMENT '是否使用：0 未用 1 已用',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_phone (phone),
    INDEX idx_expires (expires_at),
    INDEX idx_used (used)
) ENGINE=InnoDB COMMENT='短信验证码表';

-- 博文表
CREATE TABLE posts (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(200) NOT NULL COMMENT '标题',
    slug VARCHAR(200) UNIQUE COMMENT 'URL 别名',
    summary VARCHAR(500) COMMENT '摘要',
    content LONGTEXT NOT NULL COMMENT '内容 Markdown',
    content_html LONGTEXT COMMENT 'HTML 内容',
    cover_image VARCHAR(255) COMMENT '封面图',
    author_id BIGINT UNSIGNED NOT NULL COMMENT '作者 ID',
    category_id BIGINT UNSIGNED COMMENT '分类 ID',
    tags JSON COMMENT '标签数组',
    view_count INT DEFAULT 0 COMMENT '浏览次数',
    like_count INT DEFAULT 0 COMMENT '点赞数',
    share_count INT DEFAULT 0 COMMENT '分享数',
    comment_count INT DEFAULT 0 COMMENT '评论数',
    status TINYINT DEFAULT 1 COMMENT '状态：1 发布 0 草稿 -1 删除',
    is_top TINYINT DEFAULT 0 COMMENT '是否置顶',
    seo_title VARCHAR(200) COMMENT 'SEO 标题',
    seo_keywords VARCHAR(500) COMMENT 'SEO 关键词',
    seo_description VARCHAR(500) COMMENT 'SEO 描述',
    published_at DATETIME COMMENT '发布时间',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (author_id) REFERENCES users(id),
    INDEX idx_slug (slug),
    INDEX idx_author (author_id),
    INDEX idx_category (category_id),
    INDEX idx_status (status),
    INDEX idx_published (published_at),
    INDEX idx_view_count (view_count),
    INDEX idx_is_top (is_top),
    FULLTEXT INDEX idx_search (title, summary, content)
) ENGINE=InnoDB COMMENT='博文表';

-- 分类表
CREATE TABLE categories (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE COMMENT '分类名',
    slug VARCHAR(50) UNIQUE COMMENT 'URL 别名',
    description VARCHAR(200) COMMENT '描述',
    parent_id BIGINT UNSIGNED COMMENT '父分类 ID',
    sort_order INT DEFAULT 0 COMMENT '排序',
    post_count INT DEFAULT 0 COMMENT '文章数',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (parent_id) REFERENCES categories(id),
    INDEX idx_slug (slug),
    INDEX idx_parent (parent_id)
) ENGINE=InnoDB COMMENT='分类表';

-- 敏感词表
CREATE TABLE sensitive_words (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    word VARCHAR(100) NOT NULL UNIQUE COMMENT '敏感词',
    category VARCHAR(50) DEFAULT 'default' COMMENT '分类',
    level TINYINT DEFAULT 1 COMMENT '级别：1 一般 2 严重',
    status TINYINT DEFAULT 1 COMMENT '状态：1 启用 0 禁用',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_word (word),
    INDEX idx_status (status),
    INDEX idx_level (level)
) ENGINE=InnoDB COMMENT='敏感词表';

-- 评论表
CREATE TABLE comments (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    post_id BIGINT UNSIGNED NOT NULL COMMENT '博文 ID',
    parent_id BIGINT UNSIGNED COMMENT '父评论 ID',
    user_name VARCHAR(50) NOT NULL COMMENT '用户名/昵称',
    user_email VARCHAR(100) COMMENT '邮箱',
    user_ip VARCHAR(45) NOT NULL COMMENT 'IP 地址',
    user_location VARCHAR(100) COMMENT 'IP 地理位置',
    user_agent VARCHAR(500) COMMENT 'User-Agent',
    content TEXT NOT NULL COMMENT '评论内容',
    status TINYINT DEFAULT 1 COMMENT '状态：1 待审核 2 通过 3 拒绝',
    like_count INT DEFAULT 0 COMMENT '点赞数',
    is_admin TINYINT DEFAULT 0 COMMENT '是否管理员',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES comments(id) ON DELETE CASCADE,
    INDEX idx_post (post_id),
    INDEX idx_parent (parent_id),
    INDEX idx_status (status),
    INDEX idx_created (created_at)
) ENGINE=InnoDB COMMENT='评论表';

-- 点赞记录表
CREATE TABLE likes (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED COMMENT '用户 ID（匿名则为空）',
    session_id VARCHAR(64) COMMENT '会话 ID（匿名用户）',
    post_id BIGINT UNSIGNED COMMENT '博文 ID',
    comment_id BIGINT UNSIGNED COMMENT '评论 ID',
    ip_address VARCHAR(45) COMMENT 'IP 地址',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    UNIQUE KEY uk_user_post (user_id, post_id),
    UNIQUE KEY uk_session_post (session_id, post_id),
    UNIQUE KEY uk_user_comment (user_id, comment_id),
    UNIQUE KEY uk_session_comment (session_id, comment_id),
    INDEX idx_post (post_id),
    INDEX idx_comment (comment_id)
) ENGINE=InnoDB COMMENT='点赞记录表';

-- 上传文件表
CREATE TABLE uploads (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    filename VARCHAR(255) NOT NULL COMMENT '文件名',
    original_name VARCHAR(255) COMMENT '原始文件名',
    file_path VARCHAR(500) NOT NULL COMMENT '文件路径',
    file_url VARCHAR(500) NOT NULL COMMENT '访问 URL',
    file_type VARCHAR(50) NOT NULL COMMENT '文件类型',
    file_size BIGINT NOT NULL COMMENT '文件大小',
    mime_type VARCHAR(100) COMMENT 'MIME 类型',
    uploader_id BIGINT UNSIGNED COMMENT '上传者 ID',
    post_id BIGINT UNSIGNED COMMENT '关联博文 ID',
    status TINYINT DEFAULT 1 COMMENT '状态：1 正常 0 禁用',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (uploader_id) REFERENCES users(id),
    FOREIGN KEY (post_id) REFERENCES posts(id),
    INDEX idx_uploader (uploader_id),
    INDEX idx_post (post_id),
    INDEX idx_type (file_type)
) ENGINE=InnoDB COMMENT='上传文件表';

-- 访问统计表
CREATE TABLE statistics (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    stat_date DATE NOT NULL COMMENT '统计日期',
    page_views INT DEFAULT 0 COMMENT '页面浏览量',
    unique_visitors INT DEFAULT 0 COMMENT '独立访客',
    new_posts INT DEFAULT 0 COMMENT '新增文章',
    new_comments INT DEFAULT 0 COMMENT '新增评论',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_stat_date (stat_date),
    INDEX idx_date (stat_date)
) ENGINE=InnoDB COMMENT='访问统计表';

-- 系统配置表
CREATE TABLE settings (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    setting_key VARCHAR(100) NOT NULL UNIQUE COMMENT '配置键',
    setting_value TEXT COMMENT '配置值',
    setting_type VARCHAR(20) DEFAULT 'string' COMMENT '类型',
    description VARCHAR(200) COMMENT '描述',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_key (setting_key)
) ENGINE=InnoDB COMMENT='系统配置表';

-- 操作日志表
CREATE TABLE operation_logs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED COMMENT '操作用户 ID',
    action VARCHAR(50) NOT NULL COMMENT '操作类型',
    module VARCHAR(50) COMMENT '模块',
    description VARCHAR(500) COMMENT '描述',
    ip_address VARCHAR(45) COMMENT 'IP 地址',
    user_agent VARCHAR(500) COMMENT 'User-Agent',
    request_data JSON COMMENT '请求数据',
    response_code INT COMMENT '响应码',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    INDEX idx_user (user_id),
    INDEX idx_action (action),
    INDEX idx_created (created_at)
) ENGINE=InnoDB COMMENT='操作日志表';

-- 初始化默认数据
INSERT INTO users (username, password_hash, phone, role, status) 
VALUES ('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '13800138000', 1, 1);

INSERT INTO categories (name, slug, description, sort_order) VALUES
('技术', 'tech', '技术相关文章', 1),
('生活', 'life', '生活随笔', 2),
('随笔', 'notes', '随想笔记', 3);

INSERT INTO settings (setting_key, setting_value, description) VALUES
('site_name', '我的博客', '网站名称'),
('site_description', '个人博客系统', '网站描述'),
('site_keywords', '博客，技术，生活', '网站关键词'),
('posts_per_page', '10', '每页文章数'),
('allow_comment', 'true', '允许评论'),
('comment_audit', 'false', '评论需审核');
```

## 索引优化说明

1. **主键索引**：所有表使用 BIGINT UNSIGNED 自增主键
2. **唯一索引**：用户名、手机号、文章 slug 等唯一字段
3. **外键索引**：所有外键字段建立索引
4. **查询索引**：常用查询条件字段建立索引
5. **全文索引**：文章标题、摘要、内容建立全文索引用于搜索
6. **复合索引**：根据实际查询场景可添加复合索引

## 性能优化建议

1. 定期清理过期短信验证码
2. 点赞记录定期归档
3. 操作日志定期清理
4. 统计数据按月分区
5. 大字段（content）考虑分离存储
