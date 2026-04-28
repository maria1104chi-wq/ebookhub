# 中文个人博客系统

## 项目简介
基于 Golang + Gin + Vue3 开发的个人博客系统，支持 Markdown 编辑、敏感词过滤、SEO 优化、评论互动等功能。

## 功能特性
1. 伪静态页面，支持滚动加载和站内搜索
2. Markdown 富文本编辑器，支持图片、PDF、视频上传
3. 自动 SEO 优化和敏感词屏蔽
4. 点赞、分享、评论功能，显示访问统计
5. 管理员短信验证码登录（阿里云短信）
6. 后台管理：日常管理、统计、备份

## 技术栈
- 后端：Golang 1.21 + Gin + GORM + MySQL + Redis
- 前端：Vue 3 + Vite + Element Plus + VDitor
- 部署：Docker + Docker Compose + Caddy

## 目录结构
```
blog-system/
├── backend/          # 后端代码
├── frontend/         # 前端代码
├── docker/           # Docker 配置
└── docs/             # 文档
```

## 快速开始
详见 `deployment.md` 部署文档
