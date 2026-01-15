package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"chat-server/internal/config"
	"chat-server/internal/websocket"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 加载配置
	cfg := config.LoadConfig()

	// 创建 Hub
	hub := websocket.NewHub()

	// 创建处理器
	handler := websocket.NewHandler(hub, cfg)

	// 设置 HTTP 路由
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", handler.HandleWebSocket)
	mux.HandleFunc("/health", handler.HealthCheck)

	// 创建 HTTP 服务器
	server := &http.Server{
		Addr:         ":" + strconv.Itoa(cfg.Port),
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	// 启动服务器
	go func() {
		log.Printf("🚀 WebSocket 服务器运行在端口 %d", cfg.Port)
		log.Printf("🔗 WebSocket URL: ws://localhost:%d/ws", cfg.Port)
		log.Printf("💚 健康检查: http://localhost:%d/health", cfg.Port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 正在关闭服务器...")

	// 优雅关闭
	if err := server.Shutdown(nil); err != nil {
		log.Printf("服务器关闭失败: %v", err)
	}

	log.Println("✅ 服务器已关闭")
}
