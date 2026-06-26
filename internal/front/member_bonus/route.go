package memberbonus

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type JackpotMemberBonusRoute struct {
	App                       *fiber.App
	DBPool                    *sqlx.DB
	JackpotMemberBonusHandler *JackpotMemberBonusHandlerImpl
}

func NewJackpotMemberBonusRoute(app *fiber.App, dbPool *sqlx.DB) *JackpotMemberBonusRoute {
	return &JackpotMemberBonusRoute{
		App:                       app,
		DBPool:                    dbPool,
		JackpotMemberBonusHandler: NewJackpotMemberBonusHandler(dbPool),
	}
}

func (r *JackpotMemberBonusRoute) RegisterJackpotMemberBonusRoute() *JackpotMemberBonusRoute {
	memberBonus := r.App.Group("/api/v1/admin/jackpot")
	memberBonus.Get("/member-bonuses", r.JackpotMemberBonusHandler.GetAllMemberBonuses)
	memberBonus.Post("/member-bonuses", r.JackpotMemberBonusHandler.CreateMemberBonus)
	return r
}
