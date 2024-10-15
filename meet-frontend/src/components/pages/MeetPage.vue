<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { GeneralUtils } from '../../utils/GeneralUtils'
import { useRoute } from 'vue-router'

interface Message {
    id: number
    user: string
    content: string
}

const route = useRoute()

const roomId: string = route.params.roomId as string
const userQueryParam = route.query.user as string | undefined

const wsEventChatMSG: string = 'CHAT_MSG'
const wsEventChatInfo: string = 'CHAT_INFO'
const wsEventWebRTC: string = 'WEBRTC_SIGNALING'

const backendUrl: string = GeneralUtils.getWebSocketBaseUrl()
const systemUser: string = 'System'
const meUser: string = 'Me'

let conn: WebSocket | undefined

const username = ref<string>('')
const messages = ref<Message[]>([])
const newMessage = ref<string>('')

// WebRTC variables
const localStream = ref<MediaStream | null>(null)
const remoteStream = ref<MediaStream | null>(null)
const localVideo = ref<HTMLVideoElement | null>(null)
const remoteVideo = ref<HTMLVideoElement | null>(null)
const peerConnection = new RTCPeerConnection({
  iceServers: [
    {
      urls: [
        'stun:stun.l.google.com:19302',
        'stun:stun1.l.google.com:19302',
        'stun:stun2.l.google.com:19302',
        'stun:stun3.l.google.com:19302',
        'stun:stun4.l.google.com:19302'
      ]
    }
  ]
})

onMounted(async () => {
  console.log(`room: ${roomId}`)
  username.value = await getUsername()
  await setupLocalMedia()
  configWebSocket()

  // Configuración del stream remoto
  peerConnection.ontrack = (event) => {
    if (!remoteStream.value) {
      remoteStream.value = new MediaStream()
      if (remoteVideo.value) {
        remoteVideo.value.srcObject = remoteStream.value
      }
    }
    event.streams[0].getTracks().forEach(track => remoteStream.value?.addTrack(track))
  }

  startWebRTCSocket()
})

function startWebRTCSocket () {
  setTimeout(async () => {
    if (conn!.readyState === conn!.OPEN) {
      await createAndSendOffer()
    } else {
      startWebRTCSocket()
    }
  }, 1000)
}

function getUsername (): Promise<string> {
  return new Promise((resolve, reject) => {
    if (userQueryParam) {
      resolve(userQueryParam)
    } else {
      const usernamePrompt = prompt('username?') ?? ''
      if (usernamePrompt) {
        resolve(usernamePrompt)
      } else {
        reject(new Error('error getting username'))
      }
    }
  })
}

function configWebSocket () {
  if (window.WebSocket) {
    const wsURL = `${backendUrl}/ws/${username.value}`

    conn = new WebSocket(wsURL)
    conn.onclose = function (evt) {
      appendChatMessage(systemUser, 'Connection closed')
    }

    conn.onmessage = function (message) {
      const jsonEventMsg = JSON.parse(message.data)
      if (jsonEventMsg.event === wsEventChatMSG) {
        if (jsonEventMsg.from === username.value) {
          appendChatMessage(meUser, jsonEventMsg.message)
        } else {
          appendChatMessage(jsonEventMsg.from, jsonEventMsg.message)
        }
      } else if (jsonEventMsg.event === wsEventChatInfo) {
        appendChatMessage(systemUser, jsonEventMsg.message)
      } else if (jsonEventMsg.event === wsEventWebRTC) {
        handleWebRTCSignalingData(jsonEventMsg.message)
      }
    }
  } else {
    appendChatMessage(systemUser, 'Your browser does not support WebSockets')
  }
}

// Configuración de los medios locales (video y audio)
async function setupLocalMedia () {
  try {
    localStream.value = await navigator.mediaDevices.getUserMedia({ video: true, audio: true })
    if (localVideo.value && localStream.value) {
      localVideo.value.srcObject = localStream.value
    }
    // Agregar las pistas locales a la conexión WebRTC
    localStream.value?.getTracks().forEach(track => peerConnection.addTrack(track, localStream.value!))
  } catch (error) {
    console.error('Error al obtener medios locales:', error)
  }
}

// Función para manejar WebSocket
function handleWebRTCSignalingData (data: any) {
  switch (data.type) {
    case 'offer':
      peerConnection.setRemoteDescription(new RTCSessionDescription(data.offer))
      createAndSendAnswer()
      break
    case 'answer':
      peerConnection.setRemoteDescription(new RTCSessionDescription(data.answer))
      break
    case 'candidate':
      peerConnection.addIceCandidate(new RTCIceCandidate(data.candidate))
      break
    default:
      break
  }
}

// Enviar mensajes de señalización
function sendWebRTCSignalingMessage (message: any) {
  conn!.send(JSON.stringify(message))
}

// Crear y enviar oferta
async function createAndSendOffer () {
  const offer = await peerConnection.createOffer()
  await peerConnection.setLocalDescription(offer)
  sendWebRTCSignalingMessage({
    event: wsEventWebRTC,
    data: {
      type: 'offer',
      offer
    }
  })
}

// Crear y enviar respuesta
async function createAndSendAnswer () {
  const answer = await peerConnection.createAnswer()
  await peerConnection.setLocalDescription(answer)
  sendWebRTCSignalingMessage({
    event: wsEventWebRTC,
    data: {
      type: 'answer',
      answer
    }
  })
}

// Manejar el intercambio de ICE candidates
peerConnection.onicecandidate = (event) => {
  if (event.candidate) {
    sendWebRTCSignalingMessage({
      event: wsEventWebRTC,
      data: {
        type: 'candidate',
        candidate: event.candidate
      }
    })
  }
}

// Función para enviar el mensaje
function sendMessage () {
  if (newMessage.value.trim() !== '') {
    const wsMsg = {
      from: username.value,
      event: wsEventChatMSG,
      data: newMessage.value.trim()
    }
        conn!.send(JSON.stringify(wsMsg))
        newMessage.value = ''
  }
}

function appendChatMessage (usernameIn: string, messageIn: string) {
  messages.value.push({
    id: Date.now(),
    user: usernameIn,
    content: messageIn
  })
}
</script>

<template>
    <div class="meet-page chat-video-container">

        <!-- Conference Section -->
        <div class="conference-section">
            <div class="video-container">
                <video ref="localVideo" autoplay playsinline class="video-box"></video>
                <p>Tu video</p>
            </div>
            <div class="video-container">
                <video ref="remoteVideo" autoplay playsinline class="video-box"></video>
                <p>Video remoto</p>
            </div>
        </div>

        <!-- Chat Section -->
        <div class="chat-section">

            <!-- Chat Messages -->
            <div class="messages-box">
                <div v-for="message in messages" :key="message.id" class="message"
                    :class="{ 'user-message': message.user === meUser, 'system-message': message.user === systemUser }">
                    <strong>{{ message.user }}:</strong> {{ message.content }}
                </div>
            </div>

            <!-- Input Box -->
            <div class="input-box">
                <input v-model="newMessage" @keyup.enter="sendMessage" type="text" placeholder="Type a message..."
                    class="message-input" />
                <button v-if="conn" @click="sendMessage" class="send-button">Send</button>
            </div>
        </div>

    </div>
</template>

<style lang="scss" scoped>
/* Layout de dos columnas */
.chat-video-container {
    display: flex;
    flex-direction: row;
    width: 100%;
    max-width: 1200px;
    border: 1px solid #ddd;
    margin: 1rem auto;
}

/* Columna de chat (izquierda) */
.chat-section {
    flex: 1;
    display: flex;
    flex-direction: column;
    border-right: 1px solid #ddd;
    padding: 10px;
}

.messages-box {
    flex-grow: 1;
    overflow-y: auto;
    padding: 10px;
    border-bottom: 1px solid #dddddd;
}

.message {
    padding: 8px;
    margin-bottom: 5px;
    background-color: #4c4e02;
    border-radius: 5px;
}

.user-message {
    background-color: #033400;
    align-self: flex-end;
}

.system-message {
  background-color: #171658;
  align-self: flex-end;
}

.input-box {
    display: flex;
    justify-content: space-between;
    padding-top: 10px;
}

.message-input {
    flex-grow: 1;
    padding: 10px;
    border: 1px solid #ddd;
    border-radius: 5px;
    margin-right: 10px;
    background-color: black;
    color: white;
}

.send-button {
    padding: 10px 15px;
    background-color: #007bff;
    color: white;
    border: none;
    border-radius: 5px;
    cursor: pointer;
}

.send-button:hover {
    background-color: #0056b3;
}

/* Columna de video (derecha) */
.conference-section {
    flex: 1;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    padding: 10px;
}

.video-container {
    flex: 1;
    text-align: center;
}

.video-box {
    width: 100%;
    max-height: 400px;
    border: 1px solid #ddd;
    background-color: black;
}

/* Media Queries para diseño responsivo */
@media (max-width: 768px) {
    .chat-video-container {
        flex-direction: column;
    }

    .conference-section {
        flex-direction: row;
        justify-content: space-between;
    }

    .video-box {
        height: 150px;
    }
}

@media (max-width: 480px) {
    .video-box {
        height: 120px;
    }

    .message-input {
        padding: 8px;
    }

    .send-button {
        padding: 8px 12px;
    }
}
</style>
