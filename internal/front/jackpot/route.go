package jackpot

import (
	memberbonus "fish_shooting_admin_backend/internal/front/member_bonus"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type JackpotRoute struct {
	App                     *fiber.App
	DBPool                  *sqlx.DB
	JackpotHandler          *JackpotHandlerImpl
	JackpotMemberBonusRoute *memberbonus.JackpotMemberBonusRoute
}

func NewJackpotRoute(app *fiber.App, dbPool *sqlx.DB) *JackpotRoute {
	return &JackpotRoute{
		App:                     app,
		DBPool:                  dbPool,
		JackpotHandler:          NewJackpotHandler(dbPool),
		JackpotMemberBonusRoute: memberbonus.NewJackpotMemberBonusRoute(app, dbPool),
	}
}

func (r *JackpotRoute) RegisterJackpotRoute() *JackpotRoute {
	jackpot := r.App.Group("/api/v1/admin/jackpot")
	jackpot.Get("/current", r.JackpotHandler.GetCurrentJackpotGlobal)
	jackpot.Get("/form", r.JackpotHandler.GetJackpotGlobalUpdateForm)
	jackpot.Put("/current", r.JackpotHandler.UpdateJackpotGlobal)
	jackpot.Get("/company-topups", r.JackpotHandler.GetAllCompanyTopups)
	jackpot.Post("/company-topups", r.JackpotHandler.CreateCompanyTopup)
	jackpot.Get("/ledger", r.JackpotHandler.GetAllLedgers)
	return r
}
