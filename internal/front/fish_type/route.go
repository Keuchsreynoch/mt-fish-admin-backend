package fishtype

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type FishTypeRoute struct {
	App             *fiber.App
	DBPool          *sqlx.DB
	FishTypeHandler *FishTypeHandlerImpl
}

func NewFishTypeRoute(app *fiber.App, dbPool *sqlx.DB) *FishTypeRoute {
	return &FishTypeRoute{
		App:           app,
		DBPool:        dbPool,
		FishTypeHandler: NewFishTypeHandler(dbPool),
	}
}

func (f *FishTypeRoute) RegisterFishTypeRoute() *FishTypeRoute {
	fishType := f.App.Group("/api/v1/admin/fish-type")
	fishType.Get("/", f.FishTypeHandler.GetFishTypes)
	return f
}
