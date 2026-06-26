package websocket

import (
	"encoding/json"
	custom_log "fish_shooting_admin_backend/pkg/logs"
	types "fish_shooting_admin_backend/pkg/model"
	"fmt"
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// WebSocketHandler struct
type WebSocketHandler struct {
	db_Pool          *sqlx.DB
	webSocketService func(*websocket.Conn) *webSocketService
	writeChannel     chan []byte
}

func NewHandler(db_Pool *sqlx.DB) *WebSocketHandler {
	handler := &WebSocketHandler{
		db_Pool: db_Pool,
		webSocketService: func(c *websocket.Conn) *webSocketService {
			raw := c.Locals("UserContext")
			var userContext types.UserContext
			if contextMap, ok := raw.(types.UserContext); ok {
				userContext = contextMap
			} else {
				custom_log.NewCustomLog("user_context_failed", "extract user context failed", "warn")
				userContext = types.UserContext{}
			}

			// store the connection in the global map
			ClientsMutex.Lock()
			Clients[userContext.UserUuid] = c
			ClientsMutex.Unlock()

			return NewwebSocketService(userContext, db_Pool)
		},
		writeChannel: make(chan []byte), 
	}
	go handler.websocketWriter()
	return handler
}

// global map to store connections
var (
	Clients      = make(map[string]*websocket.Conn)
	ClientsMutex = &sync.Mutex{}
)

func GetClient(userUuid string) (*websocket.Conn, bool) {
	ClientsMutex.Lock()
	defer ClientsMutex.Unlock()

	conn, ok := Clients[userUuid]
	return conn, ok
}

func AddClient(userUuid string, conn *websocket.Conn) {
	ClientsMutex.Lock()
	defer ClientsMutex.Unlock()

	Clients[userUuid] = conn
}

func RemoveClient(userUuid string) {
	ClientsMutex.Lock()
	defer ClientsMutex.Unlock()

	delete(Clients, userUuid)
}

func (h *WebSocketHandler) websocketWriter() {
	for message := range h.writeChannel {
		ClientsMutex.Lock()
		for _, conn := range Clients {
			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("error writing to websocket:", err)
			}
		}
		ClientsMutex.Unlock()
	}
}

func (h *WebSocketHandler) HandleWebSocket(c *websocket.Conn) {
	h.webSocketService(c)

	defer func() {
		raw := c.Locals("UserContext")
		mCtx, _ := raw.(types.UserContext)
		ClientsMutex.Lock()
		delete(Clients, mCtx.UserUuid)
		ClientsMutex.Unlock()
		c.Close()
	}()

	var (
		msg []byte
		err error
	)
	for {
		if _, msg, err = c.ReadMessage(); err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", msg)
		h.writeChannel <- msg
	}
}

func (h *WebSocketHandler) BroadcastToUser(c *fiber.Ctx) error {
	ClientsMutex.Lock()
	defer ClientsMutex.Unlock()

	userUuid := "c5b66b62-2cb0-4a2e-b704-1da97d8ed10d"
	message := "Broadcast hello"

	if conn, ok := Clients[userUuid]; ok {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("user %s not connected", userUuid)
}

func WriteMessage(message string) error {
	fmt.Println("hello websocket", message)

	ClientsMutex.Lock()
	defer ClientsMutex.Unlock()

	for _, client := range Clients {
		if err := client.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			fmt.Println("err", err)
		}
	}

	return nil
}

func WriteToUser(user_uuid string, message string) error {
	fmt.Println("message", message)

	ClientsMutex.Lock()
	defer ClientsMutex.Unlock()

	client, ok := Clients[user_uuid]
	if !ok {
		fmt.Println("user not found")
	} else {
		if err := client.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			fmt.Println("err", err)
		}
	}

	return nil
}

func WriteJSONToUser(userUUID string, payload interface{}) error {
	message, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return WriteToUser(userUUID, string(message))
}

func PublishNotificationToUser(userUUID string, payload FrontNotificationPayload) error {
	envelope := types.BroadcastResponse{
		Topic: TopicFrontNotificationCreated,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	envelope.Data = data
	return WriteJSONToUser(userUUID, envelope)
}

func PublishCoinUpdatedToUser(userUUID string, payload FrontCoinUpdatePayload) error {
	envelope := types.BroadcastResponse{
		Topic: TopicFrontCoinUpdate,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	envelope.Data = data
	return WriteJSONToUser(userUUID, envelope)
}
