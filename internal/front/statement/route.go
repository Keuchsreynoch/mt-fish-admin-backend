package statement

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type StatementRoute struct {
	App              *fiber.App
	DBPool           *sqlx.DB
	StatementHandler *StatementHandlerImpl
}

func NewStatementRoute(app *fiber.App, dbPool *sqlx.DB) *StatementRoute {
	return &StatementRoute{
		App:              app,
		DBPool:           dbPool,
		StatementHandler: NewStatementHandler(dbPool),
	}
}

func (s *StatementRoute) RegisterStatementRoute() *StatementRoute {
	statement := s.App.Group("/api/v1/admin/statements")
	statement.Get("/", s.StatementHandler.GetAllStatements)
	statement.Get("/:member_uuid", s.StatementHandler.GetStatementsByMemberUUID)
	return s
}
