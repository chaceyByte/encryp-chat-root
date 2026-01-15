# 🚀 部署指南

本项目包含前端 (Vue 3) 和后端 (Go) 两部分。

## 📋 目录

- [快速开始](#快速开始)
- [Docker 部署](#docker-部署)
- [本地开发](#本地开发)
- [生产环境](#生产环境)

---

## 快速开始

### 使用 Docker Compose (推荐)

```bash
# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

服务将在以下端口运行：
- **前端**: http://localhost:3000
- **后端**: http://localhost:8082 (WebSocket: ws://localhost:8082/ws)
- **Nginx**: http://localhost:80 (可选)

---

## Docker 部署

### 1. 构建镜像

```bash
# 构建后端镜像
docker build -t chat-server ./backend

# 构建前端镜像
docker build -t chat-frontend ./frontend
```

### 2. 运行容器

```bash
# 运行后端
docker run -d -p 8082:8082 --name chat-server chat-server

# 运行前端
docker run -d -p 3000:3000 --name chat-frontend chat-frontend
```

### 3. 使用 Docker Compose

```bash
# 开发环境
docker-compose up -d

# 生产环境
docker-compose -f docker-compose.prod.yml up -d
```

---

## 本地开发

### 后端 (Go)

```bash
cd backend

# 安装依赖
go mod download

# 运行开发服务器
go run cmd/server/main.go

# 或构建后运行
go build -o server cmd/server/main.go
./server
```

服务器将在 `http://localhost:8082` 运行。

### 前端 (Vue 3)

```bash
cd frontend

# 安装依赖
pnpm install

# 运行开发服务器
pnpm dev

# 构建生产版本
pnpm build
```

前端开发服务器将在 `http://localhost:5173` 运行。

---

## 生产环境

### 使用 Nginx 反向代理

1. 配置 Nginx (`nginx.conf`)
2. 启用 Nginx 服务

```bash
# 启动所有服务（包含 Nginx）
docker-compose up -d
```

### 环境变量配置

#### 后端环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| PORT | 服务端口 | 8082 |
| NODE_ENV | 运行环境 | development |

#### 前端环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| NODE_ENV | 运行环境 | development |
| PORT | 前端端口 | 3000 |

---

## 🔍 健康检查

### 后端健康检查

```bash
curl http://localhost:8082/health
```

响应：
```json
{
  "status": "ok"
}
```

### Docker 容器状态

```bash
docker ps
docker-compose ps
```

---

## 📊 监控与日志

### 查看日志

```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f chat-server
docker-compose logs -f chat-frontend
```

### 资源监控

```bash
# 查看容器资源使用
docker stats

# 查看特定容器
docker stats chat-server
```

---

## 🔧 故障排查

### 后端无法启动

```bash
# 查看后端日志
docker-compose logs chat-server

# 检查端口是否被占用
lsof -i :8082
```

### 前端无法连接后端

1. 确认后端服务正常运行
2. 检查 WebSocket 连接：`ws://localhost:8082/ws`
3. 检查防火墙设置

### 容器重启

```bash
# 重启所有服务
docker-compose restart

# 重启特定服务
docker-compose restart chat-server
```

---

## 🎯 性能指标

### 后端 (Go)

- **内存占用**: ~30-50MB
- **启动时间**: <1秒
- **并发连接**: 1000+
- **消息延迟**: <10ms

### 前端 (Vue 3)

- **构建时间**: ~30秒
- **静态资源大小**: ~500KB (gzipped)

---

## 📝 注意事项

1. **生产环境**：
   - 使用 HTTPS/WSS
   - 配置 SSL 证书
   - 启用防火墙
   - 设置日志轮转

2. **安全**：
   - 限制 WebSocket 连接频率
   - 使用 Rate Limiting
   - 验证用户输入
   - 定期更新依赖

3. **备份**：
   - 定期备份数据
   - 保存 Docker 镜像
   - 记录配置变更

---

## 📚 更多文档

- [前端开发文档](./frontend/README.md)
- [后端开发文档](./backend/README.md)
- [Docker 部署文档](./DOCKER_DEPLOYMENT.md)

---

## 🆘 获取帮助

如遇问题，请：
1. 检查日志文件
2. 参考故障排查章节
3. 查看相关文档

---

## 📄 许可证

MIT License
