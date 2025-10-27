# 🐳 加密聊天室 - Docker 部署指南

## 📋 部署前准备

### 系统要求
- Docker 20.10+
- Docker Compose 2.0+
- 至少 1GB 可用内存

### 端口要求
- **3000**: 前端应用端口
- **8082**: WebSocket 服务器端口
- **80/443**: Nginx 代理端口（可选）

## 🚀 快速部署

### 方式一：使用部署脚本（推荐）

```bash
# 1. 给脚本执行权限
chmod +x deploy.sh

# 2. 运行部署脚本
./deploy.sh
```

### 方式二：手动部署

```bash
# 1. 构建前端应用
pnpm build

# 2. 构建并启动所有服务
docker-compose up -d

# 3. 查看服务状态
docker-compose ps

# 4. 查看实时日志
docker-compose logs -f
```

## 📊 服务架构

```
用户浏览器
    ↓
Nginx 反向代理 (可选)
    ↓
前端静态服务器 (chat-frontend:3000)
    ↓
WebSocket 聊天服务器 (chat-server:8082)
```

## 🔧 常用命令

### 服务管理
```bash
# 启动服务
docker-compose up -d

# 停止服务
docker-compose down

# 重启服务
docker-compose restart

# 查看服务状态
docker-compose ps

# 查看实时日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f chat-server
```

### 镜像管理
```bash
# 重新构建镜像
docker-compose build --no-cache

# 清理无用镜像
docker image prune

# 查看容器资源使用
docker stats
```

## 🌐 访问应用

### 开发环境
- **前端应用**: http://localhost:3000
- **WebSocket服务器**: localhost:8082

### 生产环境（使用Nginx）
- **主域名**: http://your-domain.com
- **WebSocket**: wss://your-domain.com/ws

## ⚙️ 环境配置

### 环境变量
可以在 `docker-compose.yml` 中配置：

```yaml
environment:
  - NODE_ENV=production
  - PORT=8082
  - LOG_LEVEL=info
```

### 自定义端口
修改 `docker-compose.yml` 中的端口映射：

```yaml
ports:
  - "8080:3000"           # 前端端口
  - "8081:8082"          # WebSocket端口
```

## 🔒 安全配置

### 1. 生产环境配置
```bash
# 创建生产环境配置
cp docker-compose.yml docker-compose.prod.yml

# 修改为生产环境配置
vim docker-compose.prod.yml
```

### 2. HTTPS 配置（需要SSL证书）
```nginx
# 在 nginx.conf 中取消注释 HTTPS 配置
server {
    listen 443 ssl http2;
    server_name your-domain.com;
    
    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    
    # 其他配置...
}
```

## 📈 监控和维护

### 健康检查
所有服务都配置了健康检查：
```bash
# 手动检查健康状态
docker-compose ps

# 查看健康检查日志
docker inspect --format='{{json .State.Health}}' chat-server
```

### 日志管理
```bash
# 查看所有日志
docker-compose logs

# 实时跟踪日志
docker-compose logs -f --tail=100

# 导出日志到文件
docker-compose logs > chat-logs-$(date +%Y%m%d).log
```

### 备份和恢复
```bash
# 备份数据（如果有数据库）
docker-compose exec chat-server tar -czf backup.tar.gz /app/data

# 恢复数据
docker cp backup.tar.gz chat-server:/app/
docker-compose exec chat-server tar -xzf backup.tar.gz
```

## 🐛 故障排除

### 常见问题

1. **端口被占用**
```bash
# 检查端口占用
netstat -tulpn | grep :3000

# 停止占用端口的进程
sudo kill -9 $(lsof -ti:3000)
```

2. **容器启动失败**
```bash
# 查看详细错误信息
docker-compose logs chat-server

# 重新构建镜像
docker-compose build --no-cache
```

3. **WebSocket连接失败**
```bash
# 检查WebSocket服务状态
curl -I http://localhost:8082

# 检查网络连接
docker network ls
docker network inspect chat-room_chat-network
```

4. **内存不足**
```bash
# 查看系统资源
docker stats

# 清理无用资源
docker system prune
```

### 性能优化

1. **调整资源限制**
```yaml
services:
  chat-server:
    deploy:
      resources:
        limits:
          memory: 512M
        reservations:
          memory: 256M
```

2. **启用缓存**
```nginx
# 在 nginx.conf 中配置缓存
location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
    expires 1y;
    add_header Cache-Control "public, immutable";
}
```

## 🔄 更新部署

```bash
# 1. 拉取最新代码
git pull

# 2. 重新构建前端
pnpm build

# 3. 重新构建镜像
docker-compose build

# 4. 重启服务
docker-compose up -d

# 5. 清理旧镜像
docker image prune -f
```

---

## 📞 技术支持

如果遇到问题，请检查：
1. Docker 和 Docker Compose 版本
2. 系统资源是否充足
3. 端口是否被占用
4. 查看容器日志获取详细错误信息

**部署成功！** 🎉 现在可以通过浏览器访问你的加密聊天室了。