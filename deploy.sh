#!/bin/bash

# 加密聊天室 Docker 部署脚本

echo "🚀 开始部署加密聊天室..."

# 检查 Docker 是否安装
if ! command -v docker &> /dev/null; then
    echo "❌ Docker 未安装，请先安装 Docker"
    exit 1
fi

# 检查 Docker Compose 是否安装
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose 未安装，请先安装 Docker Compose"
    exit 1
fi

# 构建前端应用
echo "📦 构建前端应用..."
pnpm build

if [ $? -ne 0 ]; then
    echo "❌ 前端构建失败"
    exit 1
fi

# 构建 Docker 镜像
echo "🐳 构建 Docker 镜像..."
docker-compose build

if [ $? -ne 0 ]; then
    echo "❌ Docker 镜像构建失败"
    exit 1
fi

# 启动服务
echo "🚀 启动服务..."
docker-compose up -d

if [ $? -ne 0 ]; then
    echo "❌ 服务启动失败"
    exit 1
fi

# 等待服务启动
echo "⏳ 等待服务启动..."
sleep 10

# 检查服务状态
echo "🔍 检查服务状态..."
docker-compose ps

echo ""
echo "🎉 部署完成！"
echo ""
echo "📊 服务信息："
echo "  前端应用: http://localhost:3000"
echo "  WebSocket服务器: localhost:8082"
echo ""
echo "🔧 常用命令："
echo "  查看日志: docker-compose logs -f"
echo "  停止服务: docker-compose down"
echo "  重启服务: docker-compose restart"
echo "  查看状态: docker-compose ps"
echo ""
echo "💡 提示：确保端口 3000 和 8082 未被占用"