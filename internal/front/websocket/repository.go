package websocket

import (
	"fish_shooting_admin_backend/pkg/responses"

	types "fish_shooting_admin_backend/pkg/model"

	"github.com/jmoiron/sqlx"
)

type WebSocketRepo interface {
	GetLoginSession(login_session string) (bool, *responses.ErrorResponse)
}

type WebSocketRepoImpl struct {
	user_context types.UserContext
	db_Pool      *sqlx.DB
}

func NewWebSocketRepoImpl(user types.UserContext, db_Pool *sqlx.DB) *WebSocketRepoImpl {
	return &WebSocketRepoImpl{
		user_context: user,
		db_Pool:      db_Pool,
	}
}
