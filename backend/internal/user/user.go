package user

import (
	"sync"

	"chat-server/pkg/types"
	"github.com/gorilla/websocket"
)

// User 表示一个连接的用户
type User struct {
	ID       string
	Username string
	RoomKey  string
	JoinedAt int64
	IP       string
	Conn     *websocket.Conn
	Send     chan []byte // 发送消息的通道
}

// UserManager 用户管理器
type UserManager struct {
	users map[string]*User
	mu    sync.RWMutex
}

// NewUserManager 创建用户管理器
func NewUserManager() *UserManager {
	return &UserManager{
		users: make(map[string]*User),
	}
}

// Add 添加用户
func (um *UserManager) Add(user *User) {
	um.mu.Lock()
	defer um.mu.Unlock()
	um.users[user.ID] = user
}

// Get 获取用户
func (um *UserManager) Get(id string) (*User, bool) {
	um.mu.RLock()
	defer um.mu.RUnlock()
	user, ok := um.users[id]
	return user, ok
}

// Remove 移除用户
func (um *UserManager) Remove(id string) {
	um.mu.Lock()
	defer um.mu.Unlock()
	if user, ok := um.users[id]; ok {
		close(user.Send)
		delete(um.users, id)
	}
}

// GetAll 获取所有用户
func (um *UserManager) GetAll() []*User {
	um.mu.RLock()
	defer um.mu.RUnlock()
	users := make([]*User, 0, len(um.users))
	for _, user := range um.users {
		users = append(users, user)
	}
	return users
}

// GetUsersByRoom 获取房间内的所有用户
func (um *UserManager) GetUsersByRoom(roomKey string) []*User {
	um.mu.RLock()
	defer um.mu.RUnlock()
	users := make([]*User, 0)
	for _, user := range um.users {
		if user.RoomKey == roomKey {
			users = append(users, user)
		}
	}
	return users
}

// GetUserInfo 获取用户信息（不含连接）
func (um *UserManager) GetUserInfo(id string) (*types.UserInfo, bool) {
	um.mu.RLock()
	defer um.mu.RUnlock()
	user, ok := um.users[id]
	if !ok {
		return nil, false
	}
	return &types.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		RoomKey:  user.RoomKey,
		JoinedAt: user.JoinedAt,
		IP:       user.IP,
	}, true
}

// Exists 检查用户是否存在
func (um *UserManager) Exists(id string) bool {
	um.mu.RLock()
	defer um.mu.RUnlock()
	_, ok := um.users[id]
	return ok
}

// GetByConn 根据连接获取用户
func (um *UserManager) GetByConn(conn *websocket.Conn) (*User, bool) {
	um.mu.RLock()
	defer um.mu.RUnlock()
	for _, user := range um.users {
		if user.Conn == conn {
			return user, true
		}
	}
	return nil, false
}
