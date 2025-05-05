import { ref, reactive } from 'vue'

export default function useWebSocket() {
  const socket = ref(null)
  const message = ref('')
  const username = ref('')
  const inputUsername = ref('')
  const users = reactive([])
  const messages = reactive([])

  const connect = () => {

    
    socket.value = new WebSocket('ws://localhost:8888/ws_test?username=' + username.value)
    
    socket.value.onopen = () => {}

    socket.value.onclose = (event) => {
      // 清理连接状态
      console.log('连接已关闭', event.code, event.reason)
      console.log('当前用户列表更新前:', JSON.parse(JSON.stringify(users)))
      users.splice(0, users.length)
      username.value = ''
      console.log('用户列表已清空')
      if (event.code === 1006) {
        console.error('连接异常关闭')
        alert('连接异常关闭')
      }
    }

    socket.value.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        console.log('收到消息:', data);
        
        if (data.type === 'users') {
          // console.log('收到用户列表更新:', data.users)
          users.splice(0, users.length, ...data.users)
          console.log(users);
          
        } else if (data.type === 'message') {
          messages.push({
            user: data.user,
            content: data.content,
            time: new Date().toLocaleTimeString()
          })
        } else {
          console.log('未处理的消息类型:', data.type)
        }
      } catch (error) {
        console.error('消息解析失败:', error)
        console.log('原始消息内容:', event.data)
      }
    }
  }

  const registerUser = () => {
    if (inputUsername.value.trim()) {
      username.value = inputUsername.value.trim()
      connect()
    }
  }

  const sendMessage = () => {
    if (socket.value && message.value.trim() && socket.value.readyState === WebSocket.OPEN) {
      socket.value.send(JSON.stringify({
        type: 'message',
        user: username.value,
        content: message.value,
        time: new Date().toLocaleTimeString()
      }))
      message.value = ''
    }
  }

return {
  message,
  username,
  inputUsername,
  users,
  messages,
  registerUser,
  sendMessage,
}
}