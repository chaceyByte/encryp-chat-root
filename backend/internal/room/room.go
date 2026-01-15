package room

import (
	"sync"

	"chat-server/internal/user"
)

// Room 表示一个聊天房间
type Room struct {
	Key     string
	UserIDs map[string]bool
}

// RoomManager 房间管理器
type RoomManager struct {
	rooms       map[string]*Room
	mu          sync.RWMutex
	userManager *user.UserManager
}

// NewRoomManager 创建房间管理器
func NewRoomManager(userManager *user.UserManager) *RoomManager {
	return &RoomManager{
		rooms:       make(map[string]*Room),
		userManager: userManager,
	}
}

// GetOrCreate 获取或创建房间
func (rm *RoomManager) GetOrCreate(roomKey string) *Room {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if room, ok := rm.rooms[roomKey]; ok {
		return room
	}

	room := &Room{
		Key:     roomKey,
		UserIDs: make(map[string]bool),
	}
	rm.rooms[roomKey] = room
	return room
}

// Get 获取房间
func (rm *RoomManager) Get(roomKey string) (*Room, bool) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	room, ok := rm.rooms[roomKey]
	return room, ok
}

// Delete 删除房间
func (rm *RoomManager) Delete(roomKey string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	delete(rm.rooms, roomKey)
}

// AddUser 向房间添加用户
func (rm *RoomManager) AddUser(roomKey string, userID string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if room, ok := rm.rooms[roomKey]; ok {
		room.UserIDs[userID] = true
	}
}

// RemoveUser 从房间移除用户
func (rm *RoomManager) RemoveUser(roomKey string, userID string) bool {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	room, ok := rm.rooms[roomKey]
	if !ok {
		return false
	}

	delete(room.UserIDs, userID)

	// 如果房间为空，删除房间
	if len(room.UserIDs) == 0 {
		delete(rm.rooms, roomKey)
		return true
	}
	return false
}

// GetUserCount 获取房间用户数
func (rm *RoomManager) GetUserCount(roomKey string) int {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	if room, ok := rm.rooms[roomKey]; ok {
		return len(room.UserIDs)
	}
	return 0
}

// GetRoomUsers 获取房间内的用户信息
func (rm *RoomManager) GetRoomUsers(roomKey string) []*user.User {
	return rm.userManager.GetUsersByRoom(roomKey)
}
