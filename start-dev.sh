#!/bin/bash

# 开发环境启动脚本

echo "🚀 启动加密聊天室开发环境..."

# 检查是否在项目根目录
if [ ! -f "docker-compose.yml" ]; then
    echo "❌ 请在项目根目录运行此脚本"
    exit 1
fi

# 启动后端 (Go)
echo "📡 启动后端服务..."
cd backend
if [ ! -f "server" ]; then
    echo "📦 编译后端..."
    go build -o server cmd/server/main.go
fi
./server &
BACKEND_PID=$!
echo "✅ 后端服务已启动 (PID: $BACKEND_PID)"
cd ..

# 等待后端启动
sleep 2

# 启动前端 (Vue 3)
echo "💻 启动前端服务..."
cd frontend
if [ ! -d "node_modules" ]; then
    echo "📦 安装前端依赖..."
    pnpm install
fi
pnpm dev > /dev/null 2>&1 &
FRONTEND_PID=$!
echo "✅ 前端服务已启动 (PID: $FRONTEND_PID)"
cd ..

echo ""
echo "✨ 服务已全部启动！"
echo "📍 前端地址: http://localhost:5173"
echo "📍 后端地址: http://localhost:8082"
echo "📍 WebSocket: ws://localhost:8082/ws"
echo ""
echo "按 Ctrl+C 停止所有服务"

# 等待用户中断
trap "echo ''; echo '🛑 正在停止服务...'; kill $BACKEND_PID $FRONTEND_PID 2>/dev/null; echo '✅ 所有服务已停止'; exit 0" INT TERM

wait
