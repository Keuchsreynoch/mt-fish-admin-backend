package websocket

import (
	types "fish_shooting_admin_backend/pkg/model"
	"fish_shooting_admin_backend/pkg/responses"

	"github.com/jmoiron/sqlx"
)

type webSocketCreator interface {
	GetLoginSession(login_session string) (bool, *responses.ErrorResponse)
}

type webSocketService struct {
	user_context types.UserContext
	db_Pool      *sqlx.DB
}

func NewwebSocketService(user types.UserContext, db_Pool *sqlx.DB) *webSocketService {
	return &webSocketService{
		user_context: user,
		db_Pool:      db_Pool,
	}
}
