<template>
  <div class="flex h-screen bg-gray-100">
    <!-- 新消息提示 -->
    <div 
      v-if="showNewMessageAlert" 
      class="fixed top-4 left-1/2 transform -translate-x-1/2 z-50"
    >
      <div class="bg-green-500 text-white px-4 py-2 rounded-lg shadow-lg flex items-center space-x-2 animate-bounce">
        <i class="fas fa-bell"></i>
        <span>您有新消息</span>
        <button @click="hideNewMessageAlert" class="ml-2 hover:bg-green-600 rounded p-1">
          <i class="fas fa-times"></i>
        </button>
      </div>
    </div>
    
    <!-- 侧边栏 -->
    <div class="w-80 bg-white shadow-lg flex flex-col">
      <div class="p-6 border-b border-gray-200">
        <div class="flex items-center justify-between">
          <h2 class="text-xl font-semibold text-gray-800">在线用户</h2>
          <span class="bg-green-100 text-green-800 text-sm px-2 py-1 rounded-full">
            {{ onlineUsers.length }} 在线
          </span>
        </div>
      </div>
      
      <div class="flex-1 overflow-y-auto">
        <div v-for="user in onlineUsers" :key="user.id" class="p-4 border-b border-gray-100 hover:bg-gray-50 transition-colors">
          <div class="flex items-center space-x-3">
            <div class="w-3 h-3 bg-green-500 rounded-full"></div>
            <div>
              <p class="font-medium text-gray-800">{{ user.username }}</p>
              <p class="text-xs text-gray-500">{{ user.id === currentUserId ? '我' : '在线' }}</p>
              <!-- 管理员模式下显示IP -->
              <p v-if="props.isAdmin && user.ip" class="text-xs text-blue-500 font-mono">{{ user.ip }}</p>
            </div>
          </div>
        </div>
      </div>
      
      <div class="p-4 border-t border-gray-200">
        <button
          @click="disconnect"
          class="w-full bg-red-500 text-white py-2 rounded-lg font-medium hover:bg-red-600 transition-colors"
        >
          离开聊天室
        </button>
      </div>
    </div>
    
    <!-- 主聊天区域 -->
    <div class="flex-1 flex flex-col">
      <!-- 消息区域 -->
      <div ref="messagesContainer" class="flex-1 overflow-y-auto p-6 space-y-4">
        <!-- 系统消息 -->
        <div v-for="message in systemMessages" :key="message.id" class="text-center">
          <span class="inline-block bg-gray-200 text-gray-600 text-sm px-3 py-1 rounded-full">
            {{ message.content }}
          </span>
        </div>
        
        <!-- 用户消息 -->
        <div
          v-for="message in userMessages"
          :key="message.id"
          :class="[
            'flex flex-col',
            message.userId === currentUserId ? 'items-end' : 'items-start'
          ]"
        >
          <!-- 消息头部：发送者名称和时间 -->
          <div
            :class="[
              'flex items-center space-x-2 mb-1 px-2',
              message.userId === currentUserId ? 'flex-row-reverse space-x-reverse' : ''
            ]"
          >
            <span
              :class="[
                'text-sm font-medium',
                message.userId === currentUserId ? 'text-blue-600' : 'text-gray-600'
              ]"
            >
              {{ message.userId === currentUserId ? '我' : message.username }}
            </span>
            <span
              :class="[
                'text-xs',
                message.userId === currentUserId ? 'text-blue-400' : 'text-gray-400'
              ]"
            >
              {{ formatTime(message.timestamp) }}
            </span>
            <!-- 管理员模式下显示IP -->
            <span 
              v-if="props.isAdmin && message.ip && message.userId !== currentUserId"
              class="text-xs text-red-500 font-mono"
            >
              {{ message.ip }}
            </span>
          </div>
          
          <!-- 消息内容 -->
          <div
            :class="[
              'max-w-xs lg:max-w-md px-4 py-3 rounded-2xl shadow-sm',
              message.userId === currentUserId
                ? 'bg-blue-500 text-white rounded-br-none'
                : 'bg-white text-gray-800 rounded-bl-none border'
            ]"
          >
            <!-- 图片消息 -->
            <img
              v-if="message.type === 'image'"
              :src="message.content"
              class="max-w-full h-auto rounded-lg"
              @load="scrollToBottom"
            />
            
            <!-- 文本消息 -->
            <p v-else class="whitespace-pre-wrap break-words">{{ message.content }}</p>
          </div>
        </div>
      </div>
      
      <!-- 输入区域 -->
      <div class="p-6 border-t border-gray-200 bg-white">
        <div class="flex space-x-4">
          <div class="flex-1 relative">
            <textarea
              ref="messageInput"
              v-model="newMessage"
              placeholder="输入消息... (支持粘贴图片)"
              class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none transition-all"
              style="min-height: 50px; height: 50px;"
              @keydown="handleKeydown"
              @paste="handlePaste"
              @input="autoResize"
            ></textarea>
            
            <button
              @click="pasteImage"
              class="absolute right-3 top-3 text-gray-400 hover:text-gray-600 transition-colors"
              title="粘贴图片"
            >
              <i class="fas fa-image"></i>
            </button>
          </div>
          
          <button
            @click="sendMessage"
            :disabled="!newMessage.trim()"
            class="bg-blue-500 text-white px-6 py-3 rounded-lg font-medium hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors self-end"
          >
            发送
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useChat } from '../composables/useChat'
import { formatTime } from '../utils/formatTime'

interface Props {
  username: string
  roomKey: string
  isAdmin?: boolean
}

interface Emits {
  (e: 'disconnect'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const messagesContainer = ref<HTMLDivElement>()
const messageInput = ref<HTMLTextAreaElement>()
const newMessage = ref('')

const {
  isConnected,
  currentUserId,
  onlineUsers,
  userMessages,
  systemMessages,
  sendMessage: sendChatMessage,
  disconnect: disconnectChat
} = useChat(props.username, props.roomKey)

const showNewMessageAlert = ref(false)
let newMessageTimer: NodeJS.Timeout | null = null

const hideNewMessageAlert = () => {
  showNewMessageAlert.value = false
  if (newMessageTimer) {
    clearTimeout(newMessageTimer)
    newMessageTimer = null
  }
}

const showNewMessageNotification = () => {
  showNewMessageAlert.value = true
  
  // 3秒后自动隐藏
  if (newMessageTimer) {
    clearTimeout(newMessageTimer)
  }
  newMessageTimer = setTimeout(() => {
    showNewMessageAlert.value = false
  }, 3000)
}

const autoResize = () => {
  nextTick(() => {
    if (messageInput.value) {
      messageInput.value.style.height = 'auto'
      // 设置最大高度为200px，最小高度为48px
      const newHeight = Math.min(Math.max(messageInput.value.scrollHeight, 48), 200)
      messageInput.value.style.height = newHeight + 'px'
    }
  })
}

const handleKeydown = (event: KeyboardEvent) => {
  if (event.key === 'Enter' && !event.shiftKey) {
    event.preventDefault()
    sendMessage()
  }
}

const handlePaste = async (event: ClipboardEvent) => {
  const items = event.clipboardData?.items
  if (!items) return
  
  for (const item of items) {
    if (item.type.indexOf('image') !== -1) {
      event.preventDefault()
      const file = item.getAsFile()
      if (file) {
        await sendImage(file)
      }
      return
    }
  }
}

const pasteImage = async () => {
  try {
    const items = await navigator.clipboard.read()
    for (const item of items) {
      const imageTypes = item.types.filter(type => type.startsWith('image/'))
      if (imageTypes.length > 0) {
        const blob = await item.getType(imageTypes[0])
        await sendImage(new File([blob], 'pasted-image.png', { type: blob.type }))
        break
      }
    }
  } catch (error) {
    console.warn('无法读取剪贴板中的图片:', error)
  }
}

const sendImage = async (file: File) => {
  const reader = new FileReader()
  reader.onload = (e) => {
    const imageData = e.target?.result as string
    sendChatMessage('', 'image', imageData)
    // 发送图片后立即滚动到底部
    immediateScrollToBottom()
  }
  reader.readAsDataURL(file)
}

const sendMessage = () => {
  if (newMessage.value.trim()) {
    sendChatMessage(newMessage.value.trim())
    newMessage.value = ''
    autoResize()
    // 发送消息后立即滚动到底部
    immediateScrollToBottom()
  }
}

const scrollToBottom = () => {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

const smoothScrollToBottom = () => {
  nextTick(() => {
    if (messagesContainer.value) {
      const container = messagesContainer.value
      const isAtBottom = container.scrollHeight - container.scrollTop - container.clientHeight < 50
      
      // 如果用户已经在底部附近，或者有新消息，就平滑滚动到底部
      if (isAtBottom) {
        container.scrollTo({
          top: container.scrollHeight,
          behavior: 'smooth'
        })
      }
    }
  })
}

// 强制滚动到底部（用于新消息）
const forceScrollToBottom = () => {
  nextTick(() => {
    if (messagesContainer.value) {
      const container = messagesContainer.value
      container.scrollTo({
        top: container.scrollHeight,
        behavior: 'smooth'
      })
    }
  })
}

// 立即滚动到底部（无动画）
const immediateScrollToBottom = () => {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

const disconnect = () => {
  disconnectChat()
  emit('disconnect')
}

// 监听用户消息变化，自动滚动到底部
watch(userMessages, (newMessages, oldMessages) => {
  // 只要有消息就滚动到底部
  if (newMessages.length > 0) {
    // 检查是否是新消息（不是自己发送的）
    const lastMessage = newMessages[newMessages.length - 1]
    const isMyMessage = lastMessage.userId === currentUserId.value
    
    // 如果是别人发送的消息，显示新消息提示
    if (!isMyMessage) {
      showNewMessageNotification()
    }
    
    // 有新消息时立即滚动到底部
    immediateScrollToBottom()
  }
}, { deep: true, flush: 'post' })

// 监听系统消息变化
watch(systemMessages, (newMessages, oldMessages) => {
  if (newMessages.length > 0) {
    // 有系统消息时立即滚动到底部
    immediateScrollToBottom()
  }
}, { deep: true, flush: 'post' })

// 组件挂载时滚动到底部
onMounted(() => {
  scrollToBottom()
  
  // 添加键盘快捷键：Ctrl+Enter 发送消息
  const handleGlobalKeydown = (event: KeyboardEvent) => {
    if (event.ctrlKey && event.key === 'Enter') {
      event.preventDefault()
      sendMessage()
    }
  }
  
  document.addEventListener('keydown', handleGlobalKeydown)
  
  // 组件卸载时清理事件监听器
  onUnmounted(() => {
    document.removeEventListener('keydown', handleGlobalKeydown)
    disconnectChat()
  })
})
</script>