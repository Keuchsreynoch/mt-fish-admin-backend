package gameconfig

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type GameConfigRoute struct {
	App               *fiber.App
	DBPool            *sqlx.DB
	GameConfigHandler *GameConfigHandlerImpl
}

func NewGameConfigRoute(app *fiber.App, dbPool *sqlx.DB) *GameConfigRoute {
	return &GameConfigRoute{
		App:               app,
		DBPool:            dbPool,
		GameConfigHandler: NewGameConfigHandler(dbPool),
	}
}

func (r *GameConfigRoute) RegisterGameConfigRoute() *GameConfigRoute {
	gameConfig := r.App.Group("/api/v1/admin/game-config")
	gameConfig.Get("/", r.GameConfigHandler.GetGameConfig)
	gameConfig.Put("/", r.GameConfigHandler.UpdateGameConfig)
	return r
}
