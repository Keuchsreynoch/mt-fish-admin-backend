package member

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type MemberRoute struct {
	App           *fiber.App
	DBPool        *sqlx.DB
	MemberHandler *MemberHandler
}

func NewMemberRoute(app *fiber.App, dbPool *sqlx.DB) *MemberRoute {
	return &MemberRoute{
		App:           app,
		DBPool:        dbPool,
		MemberHandler: NewMemberHandler(dbPool),
	}
}

func (m *MemberRoute) RegisterMemberRoute() *MemberRoute {
	adminMember := m.App.Group("/api/v1/admin/members")
	adminMember.Get("/", m.MemberHandler.GetAllMembers)

	return m
}
