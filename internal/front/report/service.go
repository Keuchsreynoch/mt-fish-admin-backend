package report

import (
	"fish_shooting_admin_backend/pkg/responses"

	"github.com/jmoiron/sqlx"
)

type ReportServiceCreator interface {
	GetMemberReports(req ReportShowRequest) (*ReportMemberListResponse, *responses.ErrorResponse)
}

type ReportService struct {
	DBPool    *sqlx.DB
	ReportRepo ReportRepo
}

func NewReportService(dbPool *sqlx.DB) ReportServiceCreator {
	return &ReportService{
		DBPool:    dbPool,
		ReportRepo: NewReportRepoImpl(dbPool),
	}
}

func (s *ReportService) GetMemberReports(req ReportShowRequest) (*ReportMemberListResponse, *responses.ErrorResponse) {
	return s.ReportRepo.GetMemberReports(req)
}
