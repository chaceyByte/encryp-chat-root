// 聊天室相关类型定义

export interface User {
  id: string
  username: string
  joinedAt: number
}

export interface Message {
  id: string
  userId: string
  username: string
  content: string
  type: 'text' | 'image' | 'system'
  timestamp: number
  encrypted?: boolean
}

export interface Room {
  key: string
  users: User[]
  messages: Message[]
  createdAt: number
}

export interface WebSocketMessage {
  type: 'join' | 'leave' | 'message' | 'user_joined' | 'user_left' | 'user_list' | 'error'
  userId?: string
  username?: string
  roomKey?: string
  content?: string
  messageType?: 'text' | 'image'
  timestamp?: number
  encrypted?: boolean
  users?: User[]
  messageId?: string
  message?: string
}

export interface ChatState {
  isConnected: boolean
  currentUserId: string
  onlineUsers: User[]
  userMessages: Message[]
  systemMessages: Message[]
}