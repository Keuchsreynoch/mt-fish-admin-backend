package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type AuthRoute struct {
	App         *fiber.App
	DBPool      *sqlx.DB
	AuthHandler *AuthHandler
}

func NewRoute(db_pool *sqlx.DB, app *fiber.App) *AuthRoute {
	return &AuthRoute{
		App:         app,
		DBPool:      db_pool,
		AuthHandler: NewAuthHandler(db_pool),
	}
}

func (au *AuthRoute) RegisterAuthRoute() *AuthRoute {
	adminAuthLogin := au.App.Group("/api/v1/admin/auth/login")

	adminAuthLogin.Post("/", au.AuthHandler.Login)
	return au
}
