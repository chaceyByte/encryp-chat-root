import { WebSocketServer, WebSocket } from 'ws'
import { createServer } from 'http'

interface User {
  id: string
  username: string
  ws: WebSocket
  roomKey: string
  joinedAt: number
  ip?: string
}

interface Room {
  key: string
  users: User[]
}

interface WebSocketMessage {
  type: 'join' | 'leave' | 'message'
  userId?: string
  username?: string
  roomKey?: string
  content?: string
  messageType?: 'text' | 'image'
  timestamp?: number
  encrypted?: boolean
}

class ChatServer {
  private wss: WebSocketServer
  private rooms: Map<string, Room> = new Map()
  private users: Map<string, User> = new Map()

  constructor(port: number = 8080) {
    const server = createServer()
    this.wss = new WebSocketServer({ server })
    
    server.listen(port, () => {
      console.log(`🚀 WebSocket服务器运行在端口 ${port}`)
    })

    this.setupEventHandlers()
  }

  private setupEventHandlers() {
    this.wss.on('connection', (ws: WebSocket, request) => {
      // 获取客户端IP地址
      const forwarded = request.headers['x-forwarded-for']
      const ip = forwarded 
        ? (typeof forwarded === 'string' ? forwarded.split(',')[0] : forwarded[0].split(',')[0])
        : request.socket.remoteAddress || 'unknown'
      
      console.log(`新的客户端连接，IP: ${ip}`)

      // 清理同一IP的旧session
      this.cleanupOldSessions(ip, ws)

      ws.on('message', (data: Buffer) => {
        try {
          const message: WebSocketMessage = JSON.parse(data.toString())
          this.handleMessage(ws, message, ip)
        } catch (error) {
          console.error('消息解析错误:', error)
          this.sendError(ws, '消息格式错误')
        }
      })

      ws.on('close', () => {
        this.handleDisconnect(ws)
      })

      ws.on('error', (error) => {
        console.error('WebSocket错误:', error)
        this.handleDisconnect(ws)
      })
    })
  }

  private handleMessage(ws: WebSocket, message: WebSocketMessage, ip: string) {
    switch (message.type) {
      case 'join':
        this.handleJoin(ws, message, ip)
        break
      case 'leave':
        this.handleLeave(ws, message)
        break
      case 'message':
        this.handleChatMessage(ws, message)
        break
      default:
        this.sendError(ws, '未知的消息类型')
    }
  }

  private handleJoin(ws: WebSocket, message: WebSocketMessage, ip: string) {
    if (!message.userId || !message.username || !message.roomKey) {
      this.sendError(ws, '加入房间消息缺少必要字段')
      return
    }

    const userId = message.userId
    const username = message.username
    const roomKey = message.roomKey

    // 检查用户是否已存在
    if (this.users.has(userId)) {
      this.sendError(ws, '用户ID已存在')
      return
    }

    // 创建或获取房间
    let room = this.rooms.get(roomKey)
    if (!room) {
      room = { key: roomKey, users: [] }
      this.rooms.set(roomKey, room)
      console.log(`创建新房间: ${roomKey}`)
    }

    // 创建用户
    const user: User = {
      id: userId,
      username: username,
      ws: ws,
      roomKey: roomKey,
      joinedAt: Date.now(),
      ip: ip
    }

    this.users.set(userId, user)
    room.users.push(user)

    console.log(`用户 ${username} (${userId}) 加入房间 ${roomKey}`)

    // 通知房间内其他用户
    this.broadcastToRoom(roomKey, {
      type: 'user_joined',
      userId: userId,
      username: username,
      users: this.getRoomUsers(roomKey)
    }, userId)

    // 发送当前用户列表给新用户
    ws.send(JSON.stringify({
      type: 'user_list',
      users: this.getRoomUsers(roomKey)
    }))
  }

  private handleLeave(ws: WebSocket, message: WebSocketMessage) {
    if (!message.userId || !message.roomKey) {
      return
    }

    const userId = message.userId
    const roomKey = message.roomKey

    this.removeUserFromRoom(userId, roomKey)
  }

  private handleChatMessage(ws: WebSocket, message: WebSocketMessage) {
    if (!message.userId || !message.content || !message.roomKey) {
      this.sendError(ws, '聊天消息缺少必要字段')
      return
    }

    const userId = message.userId
    const roomKey = message.roomKey
    const user = this.users.get(userId)

    if (!user || user.roomKey !== roomKey) {
      this.sendError(ws, '用户不在该房间中')
      return
    }

    // 广播消息给房间内所有用户（排除发送者）
    this.broadcastToRoom(roomKey, {
      type: 'message',
      messageId: `msg_${Date.now()}_${Math.random().toString(36).substr(2, 5)}`,
      userId: userId,
      username: user.username,
      content: message.content,
      messageType: message.messageType || 'text',
      timestamp: message.timestamp || Date.now(),
      encrypted: message.encrypted || false,
      ip: user.ip // 添加IP信息
    }, userId)

    console.log(`用户 ${user.username} 在房间 ${roomKey} 发送消息`)
  }

  private handleDisconnect(ws: WebSocket) {
    // 查找断开连接的用户
    for (const [userId, user] of this.users.entries()) {
      if (user.ws === ws) {
        console.log(`用户 ${user.username} 断开连接`)
        this.removeUserFromRoom(userId, user.roomKey)
        break
      }
    }
  }

  private removeUserFromRoom(userId: string, roomKey: string) {
    const user = this.users.get(userId)
    if (!user) return

    const room = this.rooms.get(roomKey)
    if (!room) return

    // 从房间中移除用户
    room.users = room.users.filter(u => u.id !== userId)
    this.users.delete(userId)

    // 如果房间为空，删除房间
    if (room.users.length === 0) {
      this.rooms.delete(roomKey)
      console.log(`房间 ${roomKey} 已被删除（无用户）`)
    } else {
      // 通知房间内其他用户
      this.broadcastToRoom(roomKey, {
        type: 'user_left',
        userId: userId,
        username: user.username,
        users: this.getRoomUsers(roomKey)
      })
    }

    console.log(`用户 ${user.username} 离开房间 ${roomKey}`)
  }

  private getRoomUsers(roomKey: string): Omit<User, 'ws'>[] {
    const room = this.rooms.get(roomKey)
    if (!room) return []

    return room.users.map(user => ({
      id: user.id,
      username: user.username,
      roomKey: user.roomKey,
      joinedAt: user.joinedAt,
      ip: user.ip // 包含IP信息
    }))
  }

  private broadcastToRoom(roomKey: string, message: any, excludeUserId?: string) {
    const room = this.rooms.get(roomKey)
    if (!room) return

    room.users.forEach(user => {
      if (user.id !== excludeUserId && user.ws.readyState === WebSocket.OPEN) {
        user.ws.send(JSON.stringify(message))
      }
    })
  }

  private sendError(ws: WebSocket, errorMessage: string) {
    if (ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({
        type: 'error',
        message: errorMessage
      }))
    }
  }

  // 清理同一IP的旧session（但保留当前连接的WebSocket）
  private cleanupOldSessions(ip: string, currentWs: WebSocket) {
    for (const [userId, user] of this.users.entries()) {
      // 只清理同一IP且WebSocket连接已关闭的旧session
      if (user.ip === ip && user.ws !== currentWs && user.ws.readyState !== WebSocket.OPEN) {
        console.log(`清理旧session: ${user.username} (${userId})`)
        this.removeUserFromRoom(userId, user.roomKey)
      }
    }
  }
}

// 启动服务器
const server = new ChatServer(8082)

export default server