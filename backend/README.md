# 🔒 聊天室后端服务 (Go)

基于 Go 语言实现的高性能 WebSocket 聊天室后端服务。

## ✨ 特性

- 🚀 **高性能**: 内存占用低，支持高并发连接
- 💬 **实时通信**: 基于 WebSocket 的即时消息传输
- 🔐 **密钥认证**: 用户通过相同密钥进入同一聊天室
- 👥 **房间管理**: 支持多房间、用户加入/离开通知
- 🛡️ **连接管理**: 自动清理过期连接，IP 地址记录
- 📊 **健康检查**: 提供 HTTP 健康检查端点

## 🛠️ 技术栈

- **语言**: Go 1.21
- **WebSocket**: github.com/gorilla/websocket
- **并发**: goroutine + channel
- **容器**: Docker + Alpine

## 📁 项目结构

```
backend/
├── cmd/
│   └── server/
│       └── main.go           # 入口文件
├── internal/
│   ├── websocket/            # WebSocket 处理
│   │   ├── handler.go        # HTTP 处理器
│   │   ├── hub.go           # 连接管理中心
│   │   └── connection.go    # 连接管理
│   ├── room/                # 房间管理
│   │   └── room.go
│   ├── user/                # 用户管理
│   │   └── user.go
│   ├── message/             # 消息处理
│   │   └── message.go
│   └── config/              # 配置
│       └── config.go
├── pkg/
│   └── types/               # 共享类型
│       └── types.go
├── go.mod
├── go.sum
├── Dockerfile
└── README.md
```

## 🚀 快速开始

### 本地运行

```bash
# 进入后端目录
cd backend

# 运行服务器
go run cmd/server/main.go
```

服务器将在 `http://localhost:8082` 启动。

### Docker 运行

```bash
# 构建镜像
docker build -t chat-server .

# 运行容器
docker run -p 8082:8082 chat-server
```

## ⚙️ 配置

通过环境变量配置：

| 变量 | 说明 | 默认值 |
|------|------|--------|
| PORT | 服务端口 | 8082 |
| NODE_ENV | 运行环境 | development |

## 📡 API 端点

### WebSocket

- **URL**: `ws://localhost:8082/ws`

### HTTP

- **健康检查**: `GET http://localhost:8082/health`

## 🔄 WebSocket 消息协议

### 客户端发送消息类型

#### Join - 加入房间
```json
{
  "type": "join",
  "userId": "user_123",
  "username": "张三",
  "roomKey": "room_abc"
}
```

#### Leave - 离开房间
```json
{
  "type": "leave",
  "userId": "user_123",
  "roomKey": "room_abc"
}
```

#### Message - 发送消息
```json
{
  "type": "message",
  "userId": "user_123",
  "roomKey": "room_abc",
  "content": "你好！",
  "messageType": "text",
  "encrypted": true
}
```

### 服务端推送消息类型

#### UserJoined - 用户加入通知
```json
{
  "type": "user_joined",
  "userId": "user_123",
  "username": "张三",
  "users": [...]
}
```

#### UserLeft - 用户离开通知
```json
{
  "type": "user_left",
  "userId": "user_123",
  "username": "张三",
  "users": [...]
}
```

#### UserList - 用户列表
```json
{
  "type": "user_list",
  "users": [...]
}
```

#### Message - 聊天消息
```json
{
  "type": "message",
  "messageId": "msg_1234567890_abcde",
  "userId": "user_123",
  "username": "张三",
  "content": "你好！",
  "roomKey": "room_abc",
  "messageType": "text",
  "timestamp": 1234567890123,
  "encrypted": true,
  "ip": "192.168.1.1"
}
```

#### Error - 错误消息
```json
{
  "type": "error",
  "error": "错误信息"
}
```

## 🔧 开发说明

### 构建

```bash
# 本地构建
go build -o server cmd/server/main.go

# 交叉构建
GOOS=linux GOARCH=amd64 go build -o server-linux cmd/server/main.go
```

### 测试

```bash
# 运行测试
go test ./...

# 运行测试并查看覆盖率
go test -cover ./...
```

## 📊 性能指标

- **内存占用**: ~30-50MB (10 用户在线)
- **并发连接**: 支持 1000+ 并发
- **消息延迟**: <10ms
- **启动时间**: <1秒

## 🐳 Docker 镜像

- **基础镜像**: alpine:latest
- **镜像大小**: ~15MB
- **启动时间**: <1秒

## 🚨 注意事项

1. 确保 WebSocket 连接定期发送心跳
2. 生产环境建议使用 HTTPS/WSS
3. 建议配置反向代理 (Nginx)
4. 注意防火墙和端口配置

## 📝 与前端协议完全兼容

本后端与 Vue 3 前端完全兼容，无需修改前端代码。

## 📄 许可证

MIT License
