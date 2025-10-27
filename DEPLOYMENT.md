# 🔒 加密聊天室 - 部署指南

## 📦 项目结构

```
chat-room/
├── dist/                    # 构建产物（生产环境）
│   ├── index.html
│   └── assets/
├── server/                 # WebSocket 服务器源码
│   └── server.ts
├── package.json
└── README.md
```

## 🚀 部署方式

### 方式一：本地部署（开发环境）

```bash
# 1. 安装依赖
pnpm install

# 2. 启动开发环境（同时启动前端和WebSocket服务器）
pnpm dev

# 或者分别启动：
# 终端1：启动WebSocket服务器
pnpm dev:server

# 终端2：启动前端开发服务器
pnpm dev:client
```

### 方式二：生产环境部署

#### 1. 构建前端应用
```bash
# 构建生产版本
pnpm build

# 构建产物位于 dist/ 目录
```

#### 2. 部署WebSocket服务器
```bash
# 启动生产环境的WebSocket服务器
pnpm dev:server

# 或者使用PM2等进程管理器
pnpm add -g pm2
pm2 start server/server.ts --name "chat-server" --interpreter tsx
```

#### 3. 部署静态文件
将 `dist/` 目录部署到任意静态文件服务器：

**使用Node.js静态服务器：**
```bash
# 安装serve
pnpm add -g serve

# 启动静态文件服务器
serve -s dist -l 3000
```

**使用Nginx：**
```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    # 前端静态文件
    location / {
        root /path/to/chat-room/dist;
        try_files $uri $uri/ /index.html;
    }
    
    # WebSocket代理
    location /ws {
        proxy_pass http://localhost:8082;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
    }
}
```

**使用Docker：**
```dockerfile
# Dockerfile
FROM node:18-alpine

WORKDIR /app

# 复制项目文件
COPY package.json ./
COPY server/ ./server/
COPY dist/ ./dist/

# 安装依赖
RUN npm install -g tsx
RUN npm install

# 暴露端口
EXPOSE 8082 3000

# 启动应用
CMD ["tsx", "server/server.ts"]
```

## 🔧 环境配置

### 端口配置
- **前端开发服务器**: 5173 (可自动切换)
- **WebSocket服务器**: 8082
- **生产环境前端**: 可配置任意端口

### 环境变量（可选）
```bash
# 可设置的环境变量
export CHAT_SERVER_PORT=8082
export NODE_ENV=production
```

## 📋 部署检查清单

- [ ] 确保Node.js版本 >= 16
- [ ] 安装所有依赖：`pnpm install`
- [ ] 构建前端：`pnpm build`
- [ ] 启动WebSocket服务器：`pnpm dev:server`
- [ ] 部署静态文件到Web服务器
- [ ] 配置WebSocket代理（如使用Nginx）
- [ ] 测试聊天功能
- [ ] 配置HTTPS（生产环境推荐）

## 🔒 安全建议

1. **生产环境使用HTTPS**
2. **配置WebSocket安全策略**
3. **限制消息大小和频率**
4. **实现用户认证（可选）**
5. **定期备份聊天记录（如需要）**

## 🐛 故障排除

### 常见问题

1. **WebSocket连接失败**
   - 检查端口是否被占用
   - 验证防火墙设置
   - 检查代理配置

2. **静态文件404错误**
   - 确认dist目录文件完整
   - 检查服务器路由配置

3. **消息发送失败**
   - 检查WebSocket服务器状态
   - 验证房间密钥匹配

### 日志查看
```bash
# 查看WebSocket服务器日志
tail -f server.log

# 查看PM2日志
pm2 logs chat-server
```

## 📞 技术支持

如有部署问题，请检查：
1. 端口配置是否正确
2. 防火墙设置
3. 网络代理配置
4. 浏览器控制台错误信息

---

**部署完成！** 访问你的服务器地址即可使用加密聊天室。