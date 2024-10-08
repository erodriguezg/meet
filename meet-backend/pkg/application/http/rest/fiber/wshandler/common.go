package wshandler

import (
	"encoding/json"
	"fmt"

	"github.com/erodriguezg/meet/pkg/application/http/security"
	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type WebSocketFiberHandler interface {
	RegisterWebSocket(wsRoot string, appFiber *fiber.App)
}

// MessageObject Basic chat message object
type MessageObject struct {
	Data  json.RawMessage `json:"data"`
	From  string          `json:"from"`
	Event string          `json:"event"`
	To    *string         `json:"to,omitempty"`
}

var wsClients map[string]string

func NewMiddlewareFunction(httpSecurityService security.HttpSecurityService) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)

			anonymous := false
			identity, err := httpSecurityService.GetIdentity(c)
			if err != nil {
				anonymous = true
				identity = nil
			}

			// Your authentication process goes here. Get the Token from header and validate it
			// Extract the claims from the token and set them to the Locals
			// This is because you cannot access headers in the websocket.Conn object below
			c.Locals("ANONYMOUS", anonymous)
			c.Locals("IDENTITY", identity)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}
}

func InitWebSocketsHandlers(wsRoot string, appFiber *fiber.App) {
	wsClients = make(map[string]string)

	// Multiple event handling supported
	socketio.On(socketio.EventConnect, func(ep *socketio.EventPayload) {
		fmt.Printf("Connection event 1 - User: %s \n", ep.Kws.GetStringAttribute("user_id"))
	})

	// Custom event handling supported
	socketio.On("CUSTOM_EVENT", func(ep *socketio.EventPayload) {
		fmt.Printf("Custom event - User: %s \n", ep.Kws.GetStringAttribute("user_id"))
		// --->

		// DO YOUR BUSINESS HERE

		// --->
	})

	// On message event
	socketio.On(socketio.EventMessage, func(ep *socketio.EventPayload) {

		fmt.Printf("Message event - User: %s - Message: %s \n", ep.Kws.GetStringAttribute("user_id"), string(ep.Data))

		message := MessageObject{}

		// Unmarshal the json message
		// {
		//  "from": "<user-id>",
		//  "to": "<recipient-user-id>",
		//  "event": "CUSTOM_EVENT",
		//  "data": "hello"
		//}
		err := json.Unmarshal(ep.Data, &message)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Fire custom event based on some
		// business logic
		if message.Event != "" {
			ep.Kws.Fire(message.Event, []byte(message.Data))
		}

		// Emit the message directly to specified user
		if message.To != nil {
			err = ep.Kws.EmitTo(wsClients[*message.To], ep.Data, socketio.TextMessage)
			if err != nil {
				fmt.Println(err)
			}
		}

	})

	// On disconnect event
	socketio.On(socketio.EventDisconnect, func(ep *socketio.EventPayload) {
		// Remove the user from the local clients
		delete(wsClients, ep.Kws.GetStringAttribute("user_id"))
		fmt.Printf("Disconnection event - User: %s \n", ep.Kws.GetStringAttribute("user_id"))
	})

	// On close event
	// This event is called when the server disconnects the user actively with .Close() method
	socketio.On(socketio.EventClose, func(ep *socketio.EventPayload) {
		// Remove the user from the local clients
		delete(wsClients, ep.Kws.GetStringAttribute("user_id"))
		fmt.Printf("Close event - User: %s \n", ep.Kws.GetStringAttribute("user_id"))
	})

	// On error event
	socketio.On(socketio.EventError, func(ep *socketio.EventPayload) {
		fmt.Printf("Error event - User: %s \n", ep.Kws.GetStringAttribute("user_id"))
	})

	wsHandlers := [...]WebSocketFiberHandler{
		newMeetWebSocketFiberHandler(),
	}

	for _, wsHandler := range wsHandlers {
		wsHandler.RegisterWebSocket(wsRoot, appFiber)
	}

}
