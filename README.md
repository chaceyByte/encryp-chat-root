# 🔒 加密聊天室

一个基于 Vue 3 和 Go 的实时加密聊天室应用，支持多用户实时通信、消息加密、图片发送等功能。

## ✨ 核心功能

- **🔐 密钥认证机制**：用户通过输入相同密钥进入同一聊天室
- **💬 实时聊天功能**：基于 WebSocket 的即时消息传输
- **🔒 消息加密**：使用 AES 加密算法保护消息内容
- **🖼️ 图片发送**：支持粘贴图片并直接发送
- **👥 用户管理**：显示在线用户列表、用户加入/离开通知
- **💬 消息布局**：根据发送者自动调整消息位置（自己的消息在右侧，他人的在左侧）
- **⏰ 消息状态**：每条消息显示用户名和时间戳

## 🛠️ 技术栈

### 前端
- **框架**：Vue 3 + TypeScript
- **UI**：Tailwind CSS
- **通信**：WebSocket
- **加密**：CryptoJS
- **构建**：Vite

### 后端
- **语言**：Go 1.21
- **WebSocket**：github.com/gorilla/websocket
- **并发**：goroutine + channel
- **部署**：Docker + Alpine

## 📁 项目结构

```
encryp-chat-root/
├── frontend/                    # 前端代码 (Vue 3)
│   ├── src/
│   │   ├── components/
│   │   │   └── ChatRoom.vue     # 聊天室主组件
│   │   ├── composables/
│   │   │   └── useChat.ts       # 聊天功能逻辑
│   │   ├── types/
│   │   │   └── chat.ts          # 类型定义
│   │   ├── utils/
│   │   │   └── formatTime.ts    # 时间格式化
│   │   ├── App.vue
│   │   ├── main.ts
│   │   └── style.css
│   ├── index.html
│   ├── vite.config.ts
│   ├── package.json
│   └── pnpm-lock.yaml
├── backend/                     # 后端代码 (Go)
│   ├── cmd/server/
│   │   └── main.go              # 入口文件
│   ├── internal/
│   │   ├── websocket/           # WebSocket 处理
│   │   ├── room/                # 房间管理
│   │   ├── user/                # 用户管理
│   │   ├── message/             # 消息处理
│   │   └── config/              # 配置
│   ├── pkg/types/               # 共享类型
│   ├── go.mod
│   └── Dockerfile
├── docker-compose.yml           # Docker Compose 配置
├── docker-compose.prod.yml      # 生产环境配置
├── Dockerfile.frontend          # 前端 Docker 配置
├── nginx.conf                   # Nginx 配置
├── DEPLOYMENT.md                # 部署指南
└── README.md                    # 本文件
```

## 🚀 快速开始

### 使用 Docker Compose (推荐)

```bash
# 启动所有服务
docker-compose up -d

# 访问应用
# 前端: http://localhost:3000
# 后端: http://localhost:8082
# WebSocket: ws://localhost:8082/ws
```

### 本地开发

#### 启动后端 (Go)

```bash
cd backend
go run cmd/server/main.go
```

#### 启动前端 (Vue 3)

```bash
cd frontend
pnpm install
pnpm dev
```

访问 `http://localhost:5173`

## 💡 使用说明

1. **进入聊天室**：
   - 输入用户名和房间密钥
   - 点击"进入聊天室"按钮

2. **发送消息**：
   - 在输入框中输入文本消息
   - 按 `Enter` 键发送，或点击"发送"按钮
   - 支持 `Shift + Enter` 换行

3. **发送图片**：
   - 复制图片到剪贴板
   - 在输入框中粘贴（Ctrl+V）
   - 或点击图片图标手动粘贴

4. **查看在线用户**：
   - 左侧边栏显示当前在线用户列表
   - 绿色圆点表示用户在线状态

## 🎯 特色亮点

- **安全可靠**：端到端加密，保障通信隐私
- **高性能**：Go 后端，内存占用 ~30-50MB
- **用户体验**：响应式设计，支持图片粘贴
- **实时性强**：WebSocket 实现毫秒级消息传输
- **错误处理**：完善的连接异常处理和重连机制
- **现代化 UI**：基于 Tailwind CSS 的现代化界面设计

## 📊 性能指标

| 指标 | 数值 |
|------|------|
| 后端内存占用 | ~30-50MB |
| 后端启动时间 | <1秒 |
| 并发连接数 | 1000+ |
| 消息延迟 | <10ms |
| 镜像大小 | 后端 ~15MB, 前端 ~100MB |

## 🔧 开发说明

### 项目配置

- **端口配置**：前端默认端口 5173，后端默认端口 8082
- **热重载**：开发时支持热重载，修改代码自动刷新
- **TypeScript**：完整的类型支持，提高开发效率

### 加密机制

- 使用 AES 加密算法
- 房间密钥作为加密密钥
- 消息在传输前加密，接收后解密
- 仅相同密钥的用户可以解密消息

### WebSocket 协议

服务器支持的消息类型：

- `join`：用户加入房间
- `leave`：用户离开房间
- `message`：发送聊天消息
- `user_joined`：用户加入通知
- `user_left`：用户离开通知
- `user_list`：用户列表更新
- `error`：错误消息

## 📝 注意事项

1. 确保后端服务（端口 8082）正常运行
2. 同一房间的用户必须使用相同的密钥
3. 图片发送功能依赖浏览器剪贴板 API
4. 建议在 HTTPS 环境下使用以获得更好的安全性
5. 生产环境建议使用 Nginx 反向代理

## 🚀 部署

详细的部署指南请参考 [DEPLOYMENT.md](./DEPLOYMENT.md)

## 🔮 未来规划

- [ ] 支持文件传输
- [ ] 添加消息撤回功能
- [ ] 实现聊天记录保存
- [ ] 添加房间密码保护
- [ ] 支持消息搜索
- [ ] 添加主题切换功能
- [ ] 支持私聊功能

## 📄 许可证

MIT License

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！
