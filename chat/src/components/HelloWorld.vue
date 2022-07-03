<script setup>
import { ref, onMounted, onUnmounted, reactive, nextTick } from 'vue'

defineProps({
  msg: String
})


let socket = null;
const path = "ws://localhost:8080/ws"; // 服务器地址，服务器代码在下方
const textValue = ref("");
const chatBox = ref(null);
const texta = ref(null);
const name = new Date().getTime().toString();
const bg = randomRgb();
const chatArr = reactive([]);

// WebSocket初始化
function init() {
  if (typeof WebSocket === "undefined") {
    alert("您的浏览器不支持socket");
  } else {
    socket = new WebSocket(path);
    socket.onopen = open;
    socket.onerror = error;
    socket.onclose = closed;
    socket.onmessage = getMessage;
    window.onbeforeunload = function (e) {
      e = e || window.event;
      if (e) {
        e.returnValue = "关闭提示";
        socket.close();
      }
      socket.close();
      return "关闭提示";
    };
  }
}

function open() {
  console.log("socket连接成功");
}

function error() {
  console.log("连接错误");
}

function closed() {
  console.log("socket关闭");
}
// 监听信息
async function getMessage(msg) {
  const obj = JSON.parse(msg.data);
  chatArr.push(obj);
  await nextTick(); // 异步更新DOM
  chatBox.value.scrollTop = chatBox.value.scrollHeight; // 保持滚动条在底部
}
// 随机获取头像背景
function randomRgb() {
  let R = Math.floor(Math.random() * 130 + 110);
  let G = Math.floor(Math.random() * 130 + 110);
  let B = Math.floor(Math.random() * 130 + 110);
  return "rgb(" + R + "," + G + "," + B + ")";
}
// 发送消息
function send() {
  if (textValue.value.trim().length > 0) {
    const obj = {
      type: "send",
      uid: name,
      message: textValue.value,
      bg: bg
    };
    socket.send(JSON.stringify(obj));
    textValue.value = "";
    texta.value.focus();
  }
}

function close() {
  alert("socket已经关闭");
}

onMounted(() => {
  init();
});

onUnmounted(() => {
  socket.onclose = close;
});

// return {
//   send,
//   textValue,
//   chatArr,
//   name,
//   bg,
//   chatBox,
//   texta,
//   randomRgb
// };

</script>

<template>
  <div class="home">
    <div class="content">
      <div class="chat-box" ref="chatBox">
        <div v-for="(item, index) in chatArr" :key="index" class="chat-item">
          <div v-if="item.uid === name" class="chat-msg mine">
            <p class="msg mineBg">{{ item.message }}</p>
            <p class="user" :style="{ background: item.bg }">
              {{ item.uid }}
            </p>
          </div>
          <div v-else class="chat-msg other">
            <p class="user" :style="{ background: item.bg }">
              {{ item.uid }}
            </p>
            <p class="msg otherBg">{{ item.message }}</p>
          </div>
        </div>
      </div>
    </div>
    <div class="footer">
      <textarea placeholder="说点什么..." v-model="textValue" autofocus ref="texta" @keyup.enter="send"></textarea>
      <div class="send-box">
        <p class="send active" @click="send">发送</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
a {
  color: #42b983;
}

html,
body {
  background-color: #e8e8e8;
  user-select: none;
}

::-webkit-scrollbar {
  width: 8px;
  height: 8px;
  display: none;
}

::-webkit-scrollbar-thumb {
  background-color: #D1D1D1;
  border-radius: 3px;
  -webkit-border-radius: 3px;
  border-left: 2px solid transparent;
  border-top: 2px solid transparent;
}

* {
  margin: 0;
  padding: 0;
}

.mine {
  justify-content: flex-end;
}

.other {
  justify-content: flex-start;
}

.mineBg {
  background: #98e165;
}

.otherBg {
  background: #fff;
}

.home {
  position: fixed;
  top: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 100%;
  height: 100%;
  min-width: 360px;
  min-height: 430px;
  box-shadow: 0 0 24px 0 rgb(19 70 80 / 25%);
}

.count {
  height: 5%;
  display: flex;
  justify-content: center;
  align-items: center;
  background: #EEEAE8;
  font-size: 16px;
}

.content {
  width: 100%;
  height: 80%;
  background-color: #f4f4f4;
  overflow: hidden;
}

.footer {
  position: fixed;
  bottom: 0;
  width: 100%;
  height: 15%;
  background-color: #fff;
}

.footer textarea {
  width: 100%;
  height: 50%;
  background: #fff;
  border: 0;
  box-sizing: border-box;
  resize: none;
  outline: none;
  padding: 10px;
  font-size: 16px;
}

.send-box {
  display: flex;
  height: 40%;
  justify-content: flex-end;
  align-items: center;
}

.send {
  margin-right: 20px;
  cursor: pointer;
  border-radius: 3px;
  background: #f5f5f5;
  z-index: 21;
  font-size: 16px;
  padding: 8px 20px;
}

.send:hover {
  filter: brightness(110%);
}

.active {
  background: #98e165;
  color: #fff;
}

.chat-box {
  height: 100%;
  padding: 0 20px;
  overflow-y: auto;
}

.chat-msg {
  display: flex;
  align-items: center;
}

.user {
  font-weight: bold;
  color: #fff;
  position: relative;
  word-wrap: break-word;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  width: 60px;
  height: 60px;
  line-height: 60px;
  border-radius: 8px;
  text-align: center;
  overflow: hidden;
}

.msg {
  margin: 0 5px;
  max-width: 74%;
  white-space: normal;
  word-break: break-all;
  color: #333;
  border-radius: 8px;
  padding: 10px;
  text-align: justify;
  font-size: 16px;
  box-shadow: 0px 0px 10px #f4f4f4;
}

.chat-item {
  margin: 20px 0;
  animation: up-down 1s both;
}

@keyframes up-down {
  0% {
    opacity: 0;
    transform: translate3d(0, 20px, 0);
  }

  100% {
    opacity: 1;
    transform: none;
  }
}
</style>
