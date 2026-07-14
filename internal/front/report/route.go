package report

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type ReportRoute struct {
	App           *fiber.App
	DBPool        *sqlx.DB
	ReportHandler *ReportHandlerImpl
}

func NewReportRoute(app *fiber.App, dbPool *sqlx.DB) *ReportRoute {
	return &ReportRoute{
		App:           app,
		DBPool:        dbPool,
		ReportHandler: NewReportHandler(dbPool),
	}
}

func (s *ReportRoute) RegisterReportRoute() *ReportRoute {
	report := s.App.Group("/api/v1/admin/reports")
	report.Get("/", s.ReportHandler.GetMemberReports)
	return s
}
