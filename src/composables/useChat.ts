import { ref, onUnmounted } from 'vue'
import CryptoJS from 'crypto-js'

interface User {
  id: string
  username: string
  ip?: string
}

interface Message {
  id: string
  userId: string
  username: string
  content: string
  type: 'text' | 'image' | 'system'
  timestamp: number
  encrypted?: boolean
  ip?: string
}

export function useChat(username: string, roomKey: string) {
  const ws = ref<WebSocket | null>(null)
  const isConnected = ref(false)
  const currentUserId = ref('')
  const onlineUsers = ref<User[]>([])
  const userMessages = ref<Message[]>([])
  const systemMessages = ref<Message[]>([])
  const reconnectAttempts = ref(0)
  const maxReconnectAttempts = 5

  // 生成用户ID
  const generateUserId = () => {
    return 'user_' + Math.random().toString(36).substr(2, 9)
  }

  // 加密消息
  const encryptMessage = (content: string): string => {
    return CryptoJS.AES.encrypt(content, roomKey).toString()
  }

  // 解密消息
  const decryptMessage = (encryptedContent: string): string => {
    try {
      const bytes = CryptoJS.AES.decrypt(encryptedContent, roomKey)
      return bytes.toString(CryptoJS.enc.Utf8)
    } catch (error) {
      console.error('解密失败:', error)
      return '[加密消息解密失败]'
    }
  }

  // 连接WebSocket
  const connect = () => {
    // 使用当前页面的hostname和端口8082
    // 在Docker环境中，前端和服务器都在同一台宿主机上
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const wsUrl = `${protocol}//${window.location.hostname}:8082`
    
    try {
      ws.value = new WebSocket(wsUrl)
      
      ws.value.onopen = () => {
        console.log('WebSocket连接成功')
        isConnected.value = true
        reconnectAttempts.value = 0
        
        // 发送加入房间消息
        const joinMessage = {
          type: 'join',
          userId: currentUserId.value,
          username: username,
          roomKey: roomKey
        }
        ws.value?.send(JSON.stringify(joinMessage))
      }

      ws.value.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          handleMessage(data)
        } catch (error) {
          console.error('消息解析失败:', error)
        }
      }

      ws.value.onclose = () => {
        console.log('WebSocket连接关闭')
        isConnected.value = false
        
        // 自动重连
        if (reconnectAttempts.value < maxReconnectAttempts) {
          setTimeout(() => {
            reconnectAttempts.value++
            console.log(`尝试重连 (${reconnectAttempts.value}/${maxReconnectAttempts})`)
            connect()
          }, 2000)
        }
      }

      ws.value.onerror = (error) => {
        console.error('WebSocket错误:', error)
      }
    } catch (error) {
      console.error('WebSocket连接失败:', error)
    }
  }

  // 处理接收到的消息
  const handleMessage = (data: any) => {
    switch (data.type) {
      case 'user_joined':
        addSystemMessage(`${data.username} 加入了聊天室`)
        updateOnlineUsers(data.users)
        break
        
      case 'user_left':
        addSystemMessage(`${data.username} 离开了聊天室`)
        updateOnlineUsers(data.users)
        break
        
      case 'user_list':
        updateOnlineUsers(data.users)
        break
        
      case 'message':
        const decryptedContent = data.encrypted ? decryptMessage(data.content) : data.content
        addUserMessage({
          id: data.messageId,
          userId: data.userId,
          username: data.username,
          content: decryptedContent,
          type: data.messageType || 'text',
          timestamp: data.timestamp,
          encrypted: data.encrypted,
          ip: data.ip // 接收IP信息
        })
        break
        
      case 'error':
        addSystemMessage(`错误: ${data.message}`)
        break
    }
  }

  // 添加系统消息
  const addSystemMessage = (content: string) => {
    systemMessages.value.push({
      id: 'sys_' + Date.now(),
      userId: 'system',
      username: '系统',
      content: content,
      type: 'system',
      timestamp: Date.now()
    })
  }

  // 添加用户消息
  const addUserMessage = (message: Omit<Message, 'id'>) => {
    userMessages.value.push({
      ...message,
      id: 'msg_' + Date.now() + '_' + Math.random().toString(36).substr(2, 5)
    })
  }

  // 更新在线用户列表
  const updateOnlineUsers = (users: User[]) => {
    onlineUsers.value = users
  }

  // 发送消息
  const sendMessage = (content: string, type: 'text' | 'image' = 'text', imageData?: string) => {
    if (!ws.value || ws.value.readyState !== WebSocket.OPEN) {
      addSystemMessage('连接未就绪，无法发送消息')
      return
    }

    const messageContent = type === 'image' ? imageData || content : content
    const encryptedContent = encryptMessage(messageContent)
    
    const message = {
      type: 'message',
      userId: currentUserId.value,
      username: username,
      roomKey: roomKey,
      content: encryptedContent,
      messageType: type,
      timestamp: Date.now(),
      encrypted: true
    }

    ws.value.send(JSON.stringify(message))
    
    // 本地显示已发送的消息（不加密）
    if (type === 'text') {
      addUserMessage({
        userId: currentUserId.value,
        username: username,
        content: content,
        type: 'text',
        timestamp: Date.now()
      })
    } else if (type === 'image') {
      addUserMessage({
        userId: currentUserId.value,
        username: username,
        content: imageData || content,
        type: 'image',
        timestamp: Date.now()
      })
    }
  }

  // 断开连接
  const disconnect = () => {
    if (ws.value) {
      // 发送离开消息
      const leaveMessage = {
        type: 'leave',
        userId: currentUserId.value,
        username: username,
        roomKey: roomKey
      }
      ws.value.send(JSON.stringify(leaveMessage))
      
      ws.value.close()
      ws.value = null
    }
    isConnected.value = false
    onlineUsers.value = []
    userMessages.value = []
    systemMessages.value = []
  }

  // 初始化
  const init = () => {
    currentUserId.value = generateUserId()
    connect()
  }

  // 组件卸载时清理
  onUnmounted(() => {
    disconnect()
  })

  // 立即初始化
  init()

  return {
    isConnected,
    currentUserId,
    onlineUsers,
    userMessages,
    systemMessages,
    sendMessage,
    disconnect
  }
}