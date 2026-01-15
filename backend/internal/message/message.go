package message

import (
	"encoding/json"
	"fmt"
	"time"

	"chat-server/pkg/types"
)

// ParseMessage 解析 WebSocket 消息
func ParseMessage(data []byte) (*types.WebSocketMessage, error) {
	var msg types.WebSocketMessage
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

// MarshalMessage 序列化 WebSocket 消息
func MarshalMessage(msg interface{}) ([]byte, error) {
	return json.Marshal(msg)
}

// NewUserJoinedMessage 创建用户加入通知消息
func NewUserJoinedMessage(userID, username string, users []types.UserInfo) *types.WebSocketMessage {
	return &types.WebSocketMessage{
		Type:     "user_joined",
		UserID:   userID,
		Username: username,
		Users:    users,
	}
}

// NewUserLeftMessage 创建用户离开通知消息
func NewUserLeftMessage(userID, username string, users []types.UserInfo) *types.WebSocketMessage {
	return &types.WebSocketMessage{
		Type:     "user_left",
		UserID:   userID,
		Username: username,
		Users:    users,
	}
}

// NewUserListMessage 创建用户列表消息
func NewUserListMessage(users []types.UserInfo) *types.WebSocketMessage {
	return &types.WebSocketMessage{
		Type:  "user_list",
		Users: users,
	}
}

// NewChatMessage 创建聊天消息
func NewChatMessage(userID, username, content, roomKey, ip string, messageType string, encrypted bool) *types.WebSocketMessage {
	return &types.WebSocketMessage{
		Type:        "message",
		MessageID:   generateMessageID(),
		UserID:      userID,
		Username:    username,
		Content:     content,
		RoomKey:     roomKey,
		MessageType: messageType,
		Timestamp:   getCurrentTimestamp(),
		Encrypted:   encrypted,
		IP:          ip,
	}
}

// NewErrorMessage 创建错误消息
func NewErrorMessage(errMsg string) *types.WebSocketMessage {
	return &types.WebSocketMessage{
		Type:  "error",
		Error: errMsg,
	}
}

// generateMessageID 生成消息ID
func generateMessageID() string {
	return fmt.Sprintf("msg_%d_%s", time.Now().UnixNano(), randomString(5))
}

// getCurrentTimestamp 获取当前时间戳
func getCurrentTimestamp() int64 {
	return time.Now().UnixMilli()
}

// randomString 生成随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
