package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type AdminRoute struct {
	App          *fiber.App
	DBPool       *sqlx.DB
	AdminHandler *AdminHandlerImpl
}

func NewAdminRoute(dbPool *sqlx.DB, app *fiber.App) *AdminRoute {
	return &AdminRoute{
		App:          app,
		DBPool:       dbPool,
		AdminHandler: NewAdminHandler(dbPool),
	}
}

func (r *AdminRoute) RegisterAdminRoute() *AdminRoute {
	admin := r.App.Group("/api/v1/admin")

	admin.Get("/dashboard", r.AdminHandler.GetDashboard)
	admin.Get("/users", r.AdminHandler.ListUsers)
	admin.Post("/users", r.AdminHandler.CreateUser)
	admin.Put("/users/:id", r.AdminHandler.UpdateUser)
	admin.Get("/me", r.AdminHandler.GetUserInfo)
	admin.Get("/me/menus", r.AdminHandler.GetMyMenus)
	admin.Get("/users/:id/menus", r.AdminHandler.GetUserMenus)
	admin.Put("/users/:id/menus", r.AdminHandler.AssignMenus)
	admin.Get("/menus", r.AdminHandler.ListMenus)

	return r
}
