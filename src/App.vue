<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-500 to-purple-600">
    <div v-if="!isConnected" class="flex items-center justify-center min-h-screen p-4">
      <div 
        class="bg-white/90 backdrop-blur-lg rounded-2xl shadow-2xl p-8 w-full max-w-md"
        @click="handleFormClick"
      >
        <div class="text-center mb-6">
          <i class="fas fa-lock text-4xl text-blue-500 mb-4"></i>
          <h1 class="text-3xl font-bold text-gray-800 mb-2">加密聊天室</h1>
          <p class="text-gray-600">输入相同密钥的用户将进入同一聊天室</p>
          <!-- 管理员模式提示 -->
          <div v-if="isAdminMode" class="mt-4 p-3 bg-red-100 border border-red-300 rounded-lg">
            <div class="flex items-center justify-center space-x-2">
              <i class="fas fa-shield-alt text-red-500"></i>
              <span class="text-red-700 font-medium">管理员模式已激活</span>
            </div>
            <p class="text-xs text-red-600 mt-1">可以查看用户IP地址</p>
          </div>
        </div>
        
        <div class="space-y-4">
          <!-- 显示随机用户名 -->
          <div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm font-medium text-blue-800">您的随机昵称</p>
                <p class="text-lg font-bold text-blue-900">{{ randomUsername }}</p>
              </div>
              <button
                @click="generateRandomUsername"
                class="text-blue-500 hover:text-blue-700 transition-colors"
                title="重新生成昵称"
              >
                <i class="fas fa-redo"></i>
              </button>
            </div>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">房间密钥</label>
            <input
              v-model="roomKey"
              type="password"
              placeholder="请输入房间密钥"
              class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
              @keyup.enter="connect"
            />
          </div>
          
          <button
            @click="connect"
            :disabled="!roomKey"
            class="w-full bg-gradient-to-r from-blue-500 to-purple-600 text-white py-3 rounded-lg font-semibold transition-all hover:from-blue-600 hover:to-purple-700 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            进入聊天室
          </button>
        </div>
      </div>
    </div>
    
    <ChatRoom
      v-else
      :username="randomUsername"
      :roomKey="roomKey"
      :is-admin="isAdminMode"
      @disconnect="handleDisconnect"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import ChatRoom from './components/ChatRoom.vue'

// 随机昵称生成器
const adjectives = [
  '骄傲的', '开心的', '聪明的', '勇敢的', '优雅的', '活泼的', '神秘的', '温柔的',
  '调皮的', '稳重的', '热情的', '冷静的', '可爱的', '帅气的', '美丽的', '有趣的',
  '神秘的', '机智的', '阳光的', '浪漫的', '专注的', '勤奋的', '乐观的', '谦虚的'
]

const animals = [
  '小狮子', '小兔子', '小熊猫', '小猫咪', '小狗狗', '小狐狸', '小松鼠', '小鸟儿',
  '小海豚', '小企鹅', '小老虎', '小猴子', '小鹿儿', '小熊崽', '小蝴蝶', '小蜜蜂',
  '小海星', '小乌龟', '小金鱼', '小龙虾', '小刺猬', '小考拉', '小袋鼠', '小绵羊'
]

const generateRandomUsername = () => {
  const adj = adjectives[Math.floor(Math.random() * adjectives.length)]
  const animal = animals[Math.floor(Math.random() * animals.length)]
  return `${adj}${animal}`
}

const randomUsername = ref('')
const roomKey = ref('')
const isConnected = ref(false)
const isAdminMode = ref(false)
const clickCount = ref(0)
const lastClickTime = ref(0)

// 处理点击事件
const handleFormClick = () => {
  const now = Date.now()
  
  // 如果距离上次点击超过2秒，重置计数器
  if (now - lastClickTime.value > 2000) {
    clickCount.value = 0
  }
  
  clickCount.value++
  lastClickTime.value = now
  
  // 如果连续点击5次，激活管理员模式
  if (clickCount.value >= 5) {
    isAdminMode.value = true
    clickCount.value = 0
    console.log('管理员模式已激活')
  }
}

// 初始化时生成随机用户名
onMounted(() => {
  randomUsername.value = generateRandomUsername()
})

const connect = () => {
  if (roomKey.value) {
    isConnected.value = true
  }
}

const handleDisconnect = () => {
  isConnected.value = false
  roomKey.value = ''
  // 重新生成随机用户名
  randomUsername.value = generateRandomUsername()
}
</script>