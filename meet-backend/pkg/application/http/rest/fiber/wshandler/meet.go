package wshandler

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/fiber/v2"
)

const (
	eventChatMsg   = "CHAT_MSG"
	eventChatInfo  = "CHAT_INFO"
	eventWebRTC    = "WEBRTC_SIGNALING"
	systemChatUser = "System"
)

type chatEventMessage struct {
	Event   string `json:"event"`
	From    string `json:"from"`
	Message string `json:"message"`
}

type webRtcEventMessage struct {
	Event   string          `json:"event"`
	Message json.RawMessage `json:"message"`
}

type meetWebSocketFiberHandler struct {
}

func newMeetWebSocketFiberHandler() WebSocketFiberHandler {
	return &meetWebSocketFiberHandler{}
}

// RegisterWebSocket implements WebSocketFiberHandler.
func (port *meetWebSocketFiberHandler) RegisterWebSocket(wsRoot string, appFiber *fiber.App) {
	// endpoints
	appFiber.Get(wsRoot+"/:id", socketio.New(port.meetWebSocketEndpoint))
	// events
	socketio.On(eventChatMsg, port.chatMessageEvent)
	socketio.On(eventWebRTC, port.webRtcSignalingMessageEvent)
}

func (port *meetWebSocketFiberHandler) meetWebSocketEndpoint(kws *socketio.Websocket) {

	// Retrieve the user id from endpoint
	userId := kws.Params("id")

	// Add the connection to the list of the connected clients
	// The UUID is generated randomly and is the key that allow
	// socketio to manage Emit/EmitTo/Broadcast
	wsClients[userId] = kws.UUID

	// Every websocket connection has an optional session key => value storage
	kws.SetAttribute("user_id", userId)

	//Broadcast to all the connected users the newcomer
	broadcastMsg := chatEventMessage{
		Event:   eventChatInfo,
		From:    systemChatUser,
		Message: fmt.Sprintf("New user connected: %s and UUID: %s", userId, kws.UUID),
	}
	broadcastMsgJson, err := json.Marshal(broadcastMsg)
	if err != nil {
		return
	}

	kws.Broadcast(broadcastMsgJson, true, socketio.TextMessage)

	//Write welcome message

	welcomeMsg := chatEventMessage{
		Event:   eventChatInfo,
		From:    systemChatUser,
		Message: fmt.Sprintf("Hello user: %s with UUID: %s", userId, kws.UUID),
	}

	welcomeMsgJson, err := json.Marshal(welcomeMsg)
	if err != nil {
		return
	}

	kws.Emit(welcomeMsgJson, socketio.TextMessage)
}

func (port *meetWebSocketFiberHandler) chatMessageEvent(ep *socketio.EventPayload) {
	userId := ep.Kws.GetAttribute("user_id")

	msg := chatEventMessage{
		Event:   eventChatMsg,
		From:    fmt.Sprintf("%s", userId),
		Message: string(ep.Data),
	}

	msgJson, err := json.Marshal(msg)
	if err != nil {
		return
	}

	ep.Kws.Broadcast(msgJson, false, socketio.TextMessage)
}

func (port *meetWebSocketFiberHandler) webRtcSignalingMessageEvent(ep *socketio.EventPayload) {

	msg := webRtcEventMessage{
		Event:   eventWebRTC,
		Message: ep.Data,
	}

	msgJson, err := json.Marshal(msg)
	if err != nil {
		return
	}

	ep.Kws.Broadcast(msgJson, true, socketio.TextMessage)
}
