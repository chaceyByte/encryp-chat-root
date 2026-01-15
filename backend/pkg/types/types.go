package types

import "time"

// WebSocketMessage WebSocket 消息结构
type WebSocketMessage struct {
	Type        string     `json:"type"` // join, leave, message, user_joined, user_left, user_list, error
	UserID      string     `json:"userId,omitempty"`
	Username    string     `json:"username,omitempty"`
	RoomKey     string     `json:"roomKey,omitempty"`
	Content     string     `json:"content,omitempty"`
	MessageID   string     `json:"messageId,omitempty"`
	MessageType string     `json:"messageType,omitempty"` // text, image
	Timestamp   int64      `json:"timestamp,omitempty"`
	Encrypted   bool       `json:"encrypted,omitempty"`
	IP          string     `json:"ip,omitempty"`
	Users       []UserInfo `json:"users,omitempty"`
	Error       string     `json:"error,omitempty"`
}

// UserInfo 用户信息（不包含连接）
type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	RoomKey  string `json:"roomKey"`
	JoinedAt int64  `json:"joinedAt"`
	IP       string `json:"ip,omitempty"`
}

// ChatServerConfig 服务器配置
type ChatServerConfig struct {
	Port           int           `json:"port"`
	ReadTimeout    time.Duration `json:"readTimeout"`
	WriteTimeout   time.Duration `json:"writeTimeout"`
	PingPeriod     time.Duration `json:"pingPeriod"`
	MaxMessageSize int64         `json:"maxMessageSize"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *ChatServerConfig {
	return &ChatServerConfig{
		Port:           8082,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		PingPeriod:     54 * time.Second,
		MaxMessageSize: 512 << 20, // 512MB
	}
}
