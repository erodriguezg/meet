<script setup lang="ts">
import { onMounted, onBeforeUnmount, ref, nextTick } from 'vue'
import { GeneralUtils } from '../../utils/GeneralUtils'
import { useRoute } from 'vue-router'
import Button from 'primevue/button'
import userConnectedSoundFile from './../../assets/sounds/user-online.mp3'
import userDisconnectionSoundFile from './../../assets/sounds/user-disconnection.mp3'
import newMessageSoundFile from './../../assets/sounds/new-msg.mp3'

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

// Chat variables
const messages = ref<Message[]>([])
const newMessage = ref<string>('')
const chatContainer = ref<HTMLDivElement | null>(null)

// WebRTC variables
const videoStreams = ref<MediaStream[]>([])
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

const mediaType = ref<'video' | 'audio' | 'both' | 'none'>('none')
const isTransmitting = ref(false)

const userConnectedSound = new Audio(userConnectedSoundFile)
const userDisconnectionSound = new Audio(userDisconnectionSoundFile)
const newMessageSound = new Audio(newMessageSoundFile)

onMounted(async () => {
  console.log(`room: ${roomId}`)
  username.value = await getUsername()
  configWebSocket()

  // Configuración del stream remoto
  peerConnection.ontrack = (event) => {
    if (event.streams && event.streams[0]) {
      const remoteStream = new MediaStream()
      videoStreams.value.push(remoteStream)
      // Agregar las pistas remotas al stream remoto
      event.streams[0].getTracks().forEach(track => remoteStream.addTrack(track))
    }
  }
})

async function startTransmission () {
  await setupLocalMedia()
  startWebRTCSocket()
  isTransmitting.value = true
}

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

    conn.onmessage = async function (message) {
      const jsonEventMsg = JSON.parse(message.data)
      if (jsonEventMsg.event === wsEventChatMSG) {
        if (jsonEventMsg.from === username.value) {
          appendChatMessage(meUser, jsonEventMsg.message)
        } else {
          if (document.hidden) {
            newMessageSound.play()
          }
          appendChatMessage(jsonEventMsg.from, jsonEventMsg.message)
        }
      } else if (jsonEventMsg.event === wsEventChatInfo) {
        processChatInfoMessage(jsonEventMsg)
      } else if (jsonEventMsg.event === wsEventWebRTC) {
        await handleWebRTCSignalingData(jsonEventMsg.message)
      }
    }
  } else {
    appendChatMessage(systemUser, 'Your browser does not support WebSockets')
  }
}

function processChatInfoMessage (event: any) {
  const msg = event.message as string
  appendChatMessage(systemUser, msg)
  if (msg.toLowerCase().includes('new user connected')) {
    userConnectedSound.play()
    setTimeout(sendOfferIfConnected, 1000)
  } else if (msg.toLocaleLowerCase().includes('user disconnected')) {
    userDisconnectionSound.play()
  }
}

// Configuración de los medios locales (video y audio)
async function setupLocalMedia () {
  try {
    let mediaConstraints: MediaStreamConstraints = {}

    if (mediaType.value === 'video') {
      mediaConstraints = { video: true, audio: false }
    } else if (mediaType.value === 'audio') {
      mediaConstraints = { video: false, audio: true }
    } else if (mediaType.value === 'both') {
      mediaConstraints = { video: true, audio: true }
    } else {
      mediaConstraints = { video: false, audio: false }
    }

    const localStream = await navigator.mediaDevices.getUserMedia(mediaConstraints)
    if (localStream) {
      videoStreams.value.push(localStream)
      localStream.getTracks().forEach(track => peerConnection.addTrack(track, localStream))
    }
  } catch (error) {
    console.error('Error al obtener medios locales:', error)
  }
}

// Función para manejar WebSocket
async function handleWebRTCSignalingData (data: any) {
  switch (data.type) {
    case 'offer':
      // Al recibir una oferta, establecer la descripción remota y crear una respuesta
      try {
        await peerConnection.setRemoteDescription(new RTCSessionDescription(data.offer))
        console.log('Remote description set for offer')
        await createAndSendAnswer()
      } catch (error) {
        console.error('Error handling offer:', error)
      }
      break

    case 'answer':
      // Al recibir una respuesta, establecer la descripción remota
      try {
        await peerConnection.setRemoteDescription(new RTCSessionDescription(data.answer))
        console.log('Remote description set for answer')
      } catch (error) {
        console.error('Error setting remote description for answer:', error)
      }
      break

    case 'candidate': {
      const candidate = new RTCIceCandidate(data.candidate)
      try {
        await peerConnection.addIceCandidate(candidate)
        console.log('ICE candidate added')
      } catch (error) {
        console.error('Error adding ICE candidate:', error)
      }
    }
      break

    default:
      console.warn('Unknown signaling data type:', data.type)
      break
  }
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

// Enviar mensajes de señalización
function sendWebRTCSignalingMessage (message: any) {
  conn!.send(JSON.stringify(message))
}

// Crear y enviar oferta
async function createAndSendOffer (options?: RTCOfferOptions) {
  console.log('=> createAndSendOffer')
  const offer = await peerConnection.createOffer(options)
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
  console.log('=> createAndSendAnswer')
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

async function sendOfferIfConnected () {
  if (isTransmitting.value && peerConnection.signalingState !== 'closed') {
    await createAndSendOffer({
      offerToReceiveAudio: true,
      offerToReceiveVideo: true
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
  nextTick(() => {
    chatContainer.value!.scrollTop = chatContainer.value!.scrollHeight
  })
}

function cleanupStreams () {
  videoStreams.value.forEach((stream) => {
    stream.getTracks().forEach((track) => {
      track.stop()
    })
  })
}

onBeforeUnmount(() => {
  cleanupStreams()
})

</script>

<template>
  <div class="meet-page chat-video-container">

    <div class="controls">
      <label for="mediaType">Seleccionar qué transmitir:</label>
      <select v-model="mediaType" id="mediaType">
        <option value="video">Solo video</option>
        <option value="audio">Solo audio</option>
        <option value="both">Video y audio</option>
      </select>
      <Button @click="startTransmission" v-if="!isTransmitting">Iniciar transmisión</Button>
    </div>

    <!-- Conference Section -->
    <div v-if="videoStreams.length > 0" class="conference-section">
      <div v-for="(stream, index) in videoStreams" :key="index" class="video-container">
        <video
          :id="'remoteVideo' + index"
          ref="videoElement"
          controls
          autoplay
          playsinline
          muted
          :srcObject="stream"
          class="video-box"
        ></video>
      </div>
    </div>

    <div v-else class="conference-section">
      <p>No se han iniciado transmisiones</p>
    </div>
    <!-- Chat Section -->
    <div class="chat-section">

      <!-- Chat Messages -->
      <div class="messages-box" ref="chatContainer">
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
  max-height: 100vh;
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
  overflow-y: scroll;
  padding: 10px;
  border-bottom: 1px solid #dddddd;
  height: calc(100vh - 200px);
  max-height: calc(100vh - 200px);
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

.controls {
  margin-bottom: 20px;
  position: fixed;
  bottom: 15px;
}

.controls select,
.controls button {
  margin-right: 10px;
  padding: 10px;
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

  .messages-box {
    height: calc(50vh - 200px);
    max-height: calc(50vh - 200px);
  }

  .controls {
    bottom: -15px;
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

  .messages-box {
    height: calc(70vh - 200px);
    max-height: calc(70vh - 200px);
  }

  .controls {
    bottom: -15px;
  }
}
</style>
