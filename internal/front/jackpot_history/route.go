package jackpothistory

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type JackpotHistoryRoute struct {
	App                  *fiber.App
	DBPool               *sqlx.DB
	JackpotHistoryHandler *JackpotHistoryHandlerImpl
}

func NewJackpotHistoryRoute(app *fiber.App, dbPool *sqlx.DB) *JackpotHistoryRoute {
	return &JackpotHistoryRoute{
		App:                  app,
		DBPool:               dbPool,
		JackpotHistoryHandler: NewJackpotHistoryHandler(dbPool),
	}
}

func (r *JackpotHistoryRoute) RegisterJackpotHistoryRoute() *JackpotHistoryRoute {
	jackpotHistory := r.App.Group("/api/v1/admin/jackpot-history")
	jackpotHistory.Get("/", r.JackpotHistoryHandler.GetAllJackpotHistories)
	return r
}
