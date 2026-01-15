package websocket

import (
	"encoding/json"
	"log"
	"net/http"

	"chat-server/pkg/types"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源
	},
}

// Handler WebSocket 处理器
type Handler struct {
	hub    *Hub
	config *types.ChatServerConfig
}

// NewHandler 创建处理器
func NewHandler(hub *Hub, cfg *types.ChatServerConfig) *Handler {
	return &Handler{
		hub:    hub,
		config: cfg,
	}
}

// HandleWebSocket 处理 WebSocket 连接
func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 获取客户端 IP
	ip := getClientIP(r)
	log.Printf("新的客户端连接，IP: %s", ip)

	// 升级 HTTP 连接为 WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket 升级失败: %v", err)
		return
	}

	// 创建 WebSocket 连接
	wsConn := NewWSConnection(conn, ip, h.config)

	// 清理同一 IP 的旧 session
	h.cleanupOldSessions(ip, conn)

	// 启动连接处理
	go wsConn.WritePump()
	go wsConn.ReadPump(h.hub)
}

// ServeHTTP HTTP 处理方法
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/ws" {
		h.HandleWebSocket(w, r)
		return
	}

	if r.URL.Path == "/health" {
		h.HealthCheck(w, r)
		return
	}

	http.NotFound(w, r)
}

// HealthCheck 健康检查
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

// cleanupOldSessions 清理同一 IP 的旧 session
func (h *Handler) cleanupOldSessions(ip string, newConn interface{}) {
	users := h.hub.GetUserManager().GetAll()
	for _, u := range users {
		if u.IP == ip && u.Conn != newConn && u.Conn != nil {
			// 关闭旧连接
			u.Conn.Close()
			log.Printf("清理旧session: %s (%s)", u.Username, u.ID)
			h.hub.RemoveUserFromRoom(u.ID, u.RoomKey)
		}
	}
}

// getClientIP 获取客户端 IP 地址
func getClientIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// X-Forwarded-For 可能包含多个 IP，取第一个
		if len(forwarded) > 0 {
			return forwarded
		}
	}

	if forwarded := r.Header.Get("X-Real-IP"); forwarded != "" {
		return forwarded
	}

	// 从 RemoteAddr 获取
	return r.RemoteAddr
}
