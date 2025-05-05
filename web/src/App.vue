<template>
  <div>
    <div class="page-title">
      <span>聊天室</span>
    </div>
    <div class="chat-container">
      <div v-if="!username" class="register-modal">
        <h3>请输入用户名</h3>
        <input
          type="text"
          v-model="inputUsername"
          @keyup.enter="registerUser"
        />
        <button @click="registerUser">进入聊天室</button>
      </div>
      <div v-show="username" class="chat-main">
        <div class="user-list">
          <h3>在线用户 ({{ users.length }})</h3>
          <div v-for="user in users" :key="user.name" class="user-item">
            {{ user }}
          </div>
        </div>
        <div class="chat-content">
          <div class="message-list" ref="messageList">
            <div v-for="msg in messages" :key="msg.id" :class="msg.user==username?'msg-self':'msg-other'">
              <span>{{ msg.user }}</span> 
              <div>{{ msg.content }}</div>
            </div>
          </div>
          <div class="input-area">
            <input
              type="text"
              v-model="message"
              @keyup.enter="sendMessage"
            />
            <button @click="sendMessage">发送</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import useWebSocket from '@/hooks/useWebSocket'

const { 
  message,
  username,
  inputUsername,
  users,
  messages,
  registerUser,
  sendMessage
} = useWebSocket()


</script>

<style scoped>
.page-title {
  height: 150px;
  font-size: 70px;
  font-weight: 900;
  color:cornflowerblue;
  text-align: center;
  line-height: 150px;
}

.chat-container {
  margin: 0px auto;
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 150px;
}

.register-modal {
  background: white;
  padding: 3rem;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  text-align: center;
}

.register-modal input {
  margin: 10px 0;
  padding: 8px;
  width: 200px;
}

.register-modal button {
  padding: 10px 20px;
  background-color: #4caf50;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.chat-main {
  display: flex;
  margin-top: 20px;
  width: 50%;
  border: 1px solid #ddd;
  min-width: 550px;
  max-width: 40%;
}

.user-list {
  width: 20%;
  background: rgb(236, 236, 236);
  padding: 15px;
  border-radius: 8px;
  min-width: 150px;
}

.user-list h3 {
  margin: 0;
  padding-bottom: 10px;
}

.user-list div {
  margin: 1px auto;
  padding: 8px;
  border-radius: 4px;
  color: #4caf50;
}

.user-list div:hover {
  cursor: pointer;
  color: black;
  background: white;
}

.chat-content {
  width: 100%;
}

.message-list {
  height: 400px;
  overflow-y: auto;
  padding: 15px;
  background: rgb(245, 245, 245);
}

.msg-other {
  margin: 8px 0;
  padding: 8px;
  border-radius: 4px;
  display: flex;
  gap: 8px;
  align-items: flex-start;
}

.msg-other span {
  color: #4caf50;
  white-space: nowrap;
  line-height: 1.5;
}

.msg-other span:hover {
  cursor: pointer;
  color: black;
}

.msg-other div {
  margin-top: 0;
  white-space: pre-wrap;
  word-break: break-word;
  width: fit-content;
  max-width: 80%;
  min-width: 60px;
  line-height: 1.5;
  font-family: 'Segoe UI', system-ui, sans-serif;
  background: white;
  padding: 0 7px;
  border-radius: 7px;
}

.msg-self {
  margin: 8px 0;
  padding: 8px;
  border-radius: 4px;
  display: flex;
  flex-direction: row-reverse;
  gap: 8px;
  align-items: flex-start;
}

.msg-self span {
  color: #4caf50;
  white-space: nowrap;
  line-height: 1.5;
}

.msg-self span:hover {
  cursor: pointer;
  color: black;
}

.msg-self div {
  margin-top: 0;
  white-space: pre-wrap;
  word-break: break-word;
  width: fit-content;
  max-width: 80%;
  min-width: 60px;
  line-height: 1.5;
  font-family: 'Segoe UI', system-ui, sans-serif;
  background: white;
  padding: 0 7px;
  border-radius: 7px;
  text-align: right;
}

.input-area {
  display: flex;
  padding: 15px;
  border-top: 1px solid #ddd;
}

.input-area input {
  flex: 1;
  padding: 10px;
  margin-right: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
}
</style>
