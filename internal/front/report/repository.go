package report

import (
	custom_log "fish_shooting_admin_backend/pkg/logs"
	types "fish_shooting_admin_backend/pkg/model"
	"fish_shooting_admin_backend/pkg/postgres"
	"fish_shooting_admin_backend/pkg/responses"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ReportRepo interface {
	GetMemberReports(req ReportShowRequest) (*ReportMemberListResponse, *responses.ErrorResponse)
}

type ReportRepoImpl struct {
	DBPool *sqlx.DB
}

func NewReportRepoImpl(dbPool *sqlx.DB) ReportRepo {
	return &ReportRepoImpl{
		DBPool: dbPool,
	}
}

func (r *ReportRepoImpl) GetMemberReports(req ReportShowRequest) (*ReportMemberListResponse, *responses.ErrorResponse) {
	sqlLimit := postgres.BuildPaging(req.PageOptions.Page, req.PageOptions.Perpage)

	filters := append([]types.Filter{}, req.Filters...)
	sqlFilters, argsFilters := postgres.BuildSQLFilter(filters)
	sqlOrderBy := postgres.BuildSort(req.Sorts)
	if sqlOrderBy == "" {
		sqlOrderBy = "ORDER BY jackpot_win_amount DESC, total_win_lose DESC, member_id DESC"
	}

	whereClause := `
		WHERE s.deleted_at IS NULL
	`
	if sqlFilters != "" {
		whereClause += " AND " + sqlFilters
	}

	baseFromClause := `
		FROM tbl_statements s
		LEFT JOIN tbl_tickets t
			ON t.id = s.ticket_id
			AND t.deleted_at IS NULL
		LEFT JOIN tbl_sessions ss
			ON ss.id = s.session_id
			AND ss.deleted_at IS NULL
		LEFT JOIN tbl_bets b
			ON b.ticket_id = t.id
			AND b.deleted_at IS NULL
		LEFT JOIN tbl_members m
			ON m.id = COALESCE(s.member_id, t.member_id, b.member_id)
			AND m.deleted_at IS NULL
	`

	total := 0
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM (
			SELECT
				COALESCE(m.id, 0) AS member_id,
				COALESCE(m.user_uuid::text, '') AS member_uuid,
				COALESCE(m.user_name, '') AS member_name
			%s
			%s
			GROUP BY COALESCE(m.id, 0), COALESCE(m.user_uuid::text, ''), COALESCE(m.user_name, '')
		) grouped_reports
	`, baseFromClause, whereClause)

	if err := r.DBPool.Get(&total, countQuery, argsFilters...); err != nil {
		custom_log.NewCustomLog("get_statements_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("get_statements_failed", fmt.Errorf("error_database"))
	}

	totalReport := TotalReport{}
	totalReportQuery := fmt.Sprintf(`
		SELECT
			COALESCE(SUM(COALESCE(s.total_bet_amount, t.bet_amount, b.bet_amount, 0)), 0) AS total_bet,
			COALESCE(SUM(COALESCE(s.total_bet_amount, t.bet_amount, b.bet_amount, 0) - COALESCE(s.total_bet_invalid_amount, 0)), 0) AS total_valid_bet,
			COALESCE(SUM(COALESCE(s.payout_amount, t.payout_amount, b.payout_amount, 0) - COALESCE(s.total_bet_amount, t.bet_amount, b.bet_amount, 0)), 0) AS total_winlose
		%s
		%s
	`, baseFromClause, whereClause)

	if err := r.DBPool.Get(&totalReport, totalReportQuery, argsFilters...); err != nil {
		custom_log.NewCustomLog("get_statements_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("get_statements_failed", fmt.Errorf("error_database"))
	}

	rows := make([]reportMemberRow, 0)
	query := fmt.Sprintf(`
		SELECT
			COALESCE(m.id, 0) AS member_id,
			COALESCE(m.user_uuid::text, '') AS member_uuid,
			COALESCE(m.user_name, '') AS member_name,
			COALESCE(SUM(COALESCE(s.total_bet_amount, t.bet_amount, b.bet_amount, 0)), 0) AS total_bet_amount,
			COALESCE(SUM(GREATEST(COALESCE(s.total_bet_amount, t.bet_amount, b.bet_amount, 0) - COALESCE(s.total_bet_invalid_amount, 0), 0)), 0) AS total_valid_bet,
			COALESCE(SUM(COALESCE(s.payout_amount, t.payout_amount, b.payout_amount, 0) - COALESCE(s.total_bet_amount, t.bet_amount, b.bet_amount, 0)), 0) AS total_win_lose,
			COALESCE(SUM(COALESCE(s.jackpot_win_amount, t.jackpot_win_amount, b.jackpot_payout_amount, 0)), 0) AS jackpot_win_amount
		%s
		%s
		GROUP BY COALESCE(m.id, 0), COALESCE(m.user_uuid::text, ''), COALESCE(m.user_name, '')
		%s
		%s
	`, baseFromClause, whereClause, sqlOrderBy, sqlLimit)

	if err := r.DBPool.Select(&rows, query, argsFilters...); err != nil {
		custom_log.NewCustomLog("get_statements_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("get_statements_failed", fmt.Errorf("error_database"))
	}

	reports := make([]ReportMemberResponse, 0, len(rows))
	for _, row := range rows {
		reports = append(reports, buildReportMemberResponse(row))
	}

	return &ReportMemberListResponse{
		Reports:     reports,
		TotalReport: buildReportTotalResponse(totalReport.TotalBet, totalReport.TotalValidBet, totalReport.TotalWinLose),
		Total:       total,
	}, nil
}
