package websocket

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// WebSocketHandler struct
type WebSocketRoute struct {
	app     *fiber.App
	db      *sqlx.DB
	handler *WebSocketHandler
}

func NewRoute(app *fiber.App, db *sqlx.DB) *WebSocketRoute {
	handler := NewHandler(db)
	return &WebSocketRoute{
		app:     app,
		db:      db,
		handler: handler,
	}
}

func (w *WebSocketRoute) RegisterWebSocketRoute() *WebSocketRoute {
	v1 := w.app.Group("/api/v1/front/")
	ws := v1.Group("/websocket")
	ws.Get("/ws", websocket.New(w.handler.HandleWebSocket))
	ws.Post("/ws/:useruuid", w.handler.BroadcastToUser)

	return w
}
