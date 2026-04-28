# 部署文档

## 环境要求
- Debian 11+ 云服务器
- Docker 24.0+
- Docker Compose 2.20+
- 域名已解析到服务器 IP
- 阿里云短信服务账号

## 一、服务器准备

### 1. 安装 Docker
```bash
# 更新系统
apt update && apt upgrade -y

# 安装 Docker
curl -fsSL https://get.docker.com | bash

# 启动 Docker
systemctl enable docker
systemctl start docker

# 安装 Docker Compose
apt install docker-compose-plugin -y
```

### 2. 开放防火墙端口
```bash
# 开放 80, 443, 22 端口
ufw allow 80/tcp
ufw allow 443/tcp
ufw allow 22/tcp
ufw enable
```

## 二、项目部署

### 1. 上传项目文件
```bash
# 创建目录
mkdir -p /opt/blog-system
cd /opt/blog-system

# 上传项目文件（使用 scp 或 git clone）
```

### 2. 配置环境变量
```bash
cd /opt/blog-system/docker
cp .env.example .env

# 编辑 .env 文件，填入实际配置
nano .env
```

需要配置的内容：
- MySQL 密码
- Redis 密码
- 阿里云短信 AccessKey/Secret
- 短信签名和模板代码
- 管理员手机号

### 3. 准备敏感词库
```bash
# 将 Sensitive.txt 放到指定目录
cp /path/to/Sensitive.txt /opt/blog-system/backend/data/Sensitive.txt
```

### 4. 启动服务
```bash
cd /opt/blog-system/docker
docker compose up -d --build
```

### 5. 初始化数据库
```bash
# 进入后端容器
docker compose exec backend sh

# 运行数据库迁移
./blog-system migrate

# 导入敏感词
./blog-system import-sensitive

# 退出容器
exit
```

### 6. 查看日志
```bash
docker compose logs -f
```

## 三、Caddy 自动 HTTPS

Caddy 会自动申请 Let's Encrypt 证书，只需确保：
1. 域名正确解析到服务器
2. 80 和 443 端口开放
3. Caddyfile 中配置了正确的域名

## 四、安全建议

### 1. 系统安全
- 定期更新系统：`apt update && apt upgrade`
- 禁用 root SSH 登录
- 使用密钥认证代替密码
- 配置 fail2ban 防止暴力破解

### 2. 应用安全
- 定期备份数据库
- 监控异常访问
- 限制上传文件大小
- 启用 CORS 策略
- 使用 HTTPS

### 3. 数据安全
- 定期备份 MySQL 数据
- 敏感信息加密存储
- 限制数据库远程访问

### 4. 监控告警
- 配置日志轮转
- 监控磁盘空间
- 监控内存使用
- 设置异常告警

## 五、日常维护

### 数据库备份
```bash
docker compose exec mysql mysqldump -u root -p blog_db > backup.sql
```

### 查看运行状态
```bash
docker compose ps
```

### 重启服务
```bash
docker compose restart
```

### 更新版本
```bash
git pull
docker compose down
docker compose up -d --build
```

## 六、访问网站

部署完成后，访问 `https://你的域名` 即可使用博客系统。

管理员登录地址：`/admin/login`
