package websocket

import (
	"log"
	"time"

	"chat-server/internal/message"
	"chat-server/internal/room"
	"chat-server/internal/user"
	"chat-server/pkg/types"

	"github.com/gorilla/websocket"
)

// Hub 维护活跃的连接并广播消息
type Hub struct {
	userManager *user.UserManager
	roomManager *room.RoomManager
}

// NewHub 创建 Hub
func NewHub() *Hub {
	userManager := user.NewUserManager()
	roomManager := room.NewRoomManager(userManager)

	return &Hub{
		userManager: userManager,
		roomManager: roomManager,
	}
}

// HandleJoin 处理用户加入房间
func (h *Hub) HandleJoin(wsConn *WSConnection, userID, username, roomKey string) error {
	// 检查用户是否已存在
	if h.userManager.Exists(userID) {
		wsConn.SendError("用户ID已存在")
		return nil
	}

	// 创建用户，使用 WSConnection 的 Send channel
	newUser := &user.User{
		ID:       userID,
		Username: username,
		RoomKey:  roomKey,
		JoinedAt: getCurrentTimestamp(),
		IP:       wsConn.IP,
		Conn:     wsConn.Conn,
		Send:     wsConn.Send, // 使用 WSConnection 的 Send channel
	}

	// 添加到用户管理器
	h.userManager.Add(newUser)

	// 添加到房间
	h.roomManager.AddUser(roomKey, userID)

	log.Printf("用户 %s (%s) 加入房间 %s", username, userID, roomKey)

	// 通知房间内其他用户
	users := h.GetRoomUsersInfo(roomKey)
	userJoinedMsg := message.NewUserJoinedMessage(userID, username, users)
	h.BroadcastToRoom(roomKey, userJoinedMsg, userID)

	// 发送当前用户列表给新用户
	userListMsg := message.NewUserListMessage(users)
	wsConn.SendMessage(userListMsg)

	return nil
}

// HandleLeave 处理用户离开房间
func (h *Hub) HandleLeave(userID, roomKey string) {
	h.RemoveUserFromRoom(userID, roomKey)
}

// HandleMessage 处理聊天消息
func (h *Hub) HandleMessage(userID, roomKey, content string, messageType string, encrypted bool) error {
	usr, ok := h.userManager.Get(userID)
	if !ok {
		return nil
	}

	if usr.RoomKey != roomKey {
		return nil
	}

	// 广播消息给房间内所有用户
	chatMsg := message.NewChatMessage(
		userID,
		usr.Username,
		content,
		roomKey,
		usr.IP,
		messageType,
		encrypted,
	)

	h.BroadcastToRoom(roomKey, chatMsg, userID)

	log.Printf("用户 %s 在房间 %s 发送消息", usr.Username, roomKey)

	return nil
}

// RemoveUserFromRoom 从房间移除用户
func (h *Hub) RemoveUserFromRoom(userID, roomKey string) {
	usr, ok := h.userManager.Get(userID)
	if !ok {
		return
	}

	log.Printf("用户 %s 离开房间 %s", usr.Username, roomKey)

	// 从房间移除用户
	isEmpty := h.roomManager.RemoveUser(roomKey, userID)

	if isEmpty {
		log.Printf("房间 %s 已被删除（无用户）", roomKey)
	} else {
		// 通知房间内其他用户
		users := h.GetRoomUsersInfo(roomKey)
		userLeftMsg := message.NewUserLeftMessage(userID, usr.Username, users)
		h.BroadcastToRoom(roomKey, userLeftMsg, userID)
	}

	// 从用户管理器移除
	h.userManager.Remove(userID)
}

// BroadcastToRoom 向房间广播消息
func (h *Hub) BroadcastToRoom(roomKey string, msg *types.WebSocketMessage, excludeUserID string) {
	data, err := message.MarshalMessage(msg)
	if err != nil {
		log.Printf("消息序列化错误: %v", err)
		return
	}

	users := h.roomManager.GetRoomUsers(roomKey)
	for _, u := range users {
		if u.ID != excludeUserID {
			select {
			case u.Send <- data:
			default:
				log.Printf("用户 %s 的发送通道已满，跳过消息", u.Username)
			}
		}
	}
}

// GetRoomUsersInfo 获取房间内的用户信息列表
func (h *Hub) GetRoomUsersInfo(roomKey string) []types.UserInfo {
	users := h.roomManager.GetRoomUsers(roomKey)
	userInfos := make([]types.UserInfo, 0, len(users))

	for _, u := range users {
		userInfos = append(userInfos, types.UserInfo{
			ID:       u.ID,
			Username: u.Username,
			RoomKey:  u.RoomKey,
			JoinedAt: u.JoinedAt,
			IP:       u.IP,
		})
	}

	return userInfos
}

// GetUser 获取用户
func (h *Hub) GetUser(userID string) (*user.User, bool) {
	return h.userManager.Get(userID)
}

// GetByConn 根据连接获取用户
func (h *Hub) GetByConn(conn interface{}) (*user.User, bool) {
	wsConn, ok := conn.(*websocket.Conn)
	if !ok {
		return nil, false
	}
	return h.userManager.GetByConn(wsConn)
}

// GetUserManager 获取用户管理器
func (h *Hub) GetUserManager() *user.UserManager {
	return h.userManager
}

// getCurrentTimestamp 获取当前时间戳
func getCurrentTimestamp() int64 {
	return time.Now().UnixMilli()
}
