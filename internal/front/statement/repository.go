package statement

import (
	custom_log "fish_shooting_admin_backend/pkg/logs"
	types "fish_shooting_admin_backend/pkg/model"
	"fish_shooting_admin_backend/pkg/postgres"
	"fish_shooting_admin_backend/pkg/responses"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type StatementRepo interface {
	GetAllStatements(req StatementShowRequest) (*StatementListResponse, *responses.ErrorResponse)
	GetStatementsByMemberUUID(req StatementShowRequest, member_uuid string) (*StatementListResponse, *responses.ErrorResponse)
}

type StatementRepoImpl struct {
	DBPool      *sqlx.DB
	UserContext *types.UserContext
}

func NewStatementRepoImpl(dbPool *sqlx.DB, userContext *types.UserContext) StatementRepo {
	return &StatementRepoImpl{
		DBPool:      dbPool,
		UserContext: userContext,
	}
}

func (r *StatementRepoImpl) GetAllStatements(req StatementShowRequest) (*StatementListResponse, *responses.ErrorResponse) {
	return r.getStatements(req, "")
}

func (r *StatementRepoImpl) GetStatementsByMemberUUID(req StatementShowRequest, member_uuid string) (*StatementListResponse, *responses.ErrorResponse) {
	if member_uuid == "" {
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("get_statements_failed", fmt.Errorf("missing_member_uuid"))
	}

	return r.getStatements(req, member_uuid)
}

func (r *StatementRepoImpl) getStatements(req StatementShowRequest, member_uuid string) (*StatementListResponse, *responses.ErrorResponse) {
	sqlLimit := postgres.BuildPaging(req.PageOptions.Page, req.PageOptions.Perpage)
	filters := append([]types.Filter{}, req.Filters...)
	if member_uuid != "" {
		filters = append(filters, types.Filter{
			Property: "m.user_uuid",
			Operator: "eq",
			Value:    member_uuid,
		})
	}
	sqlFilters, argsFilters := postgres.BuildSQLFilter(filters)
	sqlOrderBy := postgres.BuildSort(req.Sorts)
	if sqlOrderBy == "" {
		sqlOrderBy = "ORDER BY COALESCE(s.statement_at, s.created_at) DESC, s.id DESC"
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
		LEFT JOIN tbl_fish_types ft
			ON ft.id = b.fish_type_id
	`

	total := 0
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		%s
		%s
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

	rows := make([]statementRow, 0)
	query := fmt.Sprintf(`
		SELECT
			COALESCE(s.statement_uuid::text, '') AS statement_uuid,
			COALESCE(s.member_id, t.member_id, b.member_id, 0) AS member_id,
			COALESCE(m.user_uuid::text, '') AS member_uuid,
			COALESCE(m.user_name, '') AS username,
			COALESCE(s.session_id, t.session_id, b.session_id, 0) AS session_id,
			COALESCE(ss.session_no, '') AS session_no,
			COALESCE(s.ticket_id, t.id, 0) AS ticket_id,
			COALESCE(s.ticket_no, t.ticket_no, '') AS ticket_no,
			COALESCE(b.bet_no, '') AS bet_no,
			COALESCE(b.fish_type_id, 0) AS fish_type_id,
			COALESCE(ft.fish_type_name, '') AS fish_type_name,
			COALESCE(s.total_bet_amount, t.bet_amount, b.bet_amount, 0) AS bet_amount,
			GREATEST(COALESCE(s.total_bet_amount, t.bet_amount, b.bet_amount, 0) - COALESCE(s.total_bet_invalid_amount, 0), 0) AS bet_valid,
			COALESCE(s.total_bet_invalid_amount, 0) AS bet_invalid,
			COALESCE(s.total_bet_amount, 0) AS total_bet_amount,
			COALESCE(s.total_bet_invalid_amount, 0) AS total_bet_invalid_amount,
			COALESCE(s.is_kill, FALSE) AS is_kill,
			CASE
				WHEN COALESCE(s.payout_amount, t.payout_amount, b.payout_amount, 0) > 0 THEN 'win'
				ELSE 'lose'
			END AS win_lose,
			COALESCE(b.kill_reward, 0) AS kill_reward,
			COALESCE(s.miss_reward, t.miss_reward, b.miss_reward, 0) AS miss_reward,
			COALESCE(s.payout_amount, t.payout_amount, b.payout_amount, 0) AS payout_amount,
			COALESCE(s.jackpot_win_amount, t.jackpot_win_amount, b.jackpot_payout_amount, 0) AS jackpot_win_amount,
			COALESCE(s.jackpot_win_amount, t.jackpot_win_amount, b.jackpot_payout_amount, 0) AS jackpot,
			COALESCE(s.sync_id, 0) AS sync_id,
			COALESCE(s.is_synced, FALSE) AS is_synced,
			COALESCE(s.statement_at, s.created_at, NOW()) AS statement_at,
			COALESCE(s.status_id, 0) AS status_id,
			COALESCE(s."order", 0) AS "order",
			COALESCE(s.created_at, NOW()) AS created_at,
			s.created_by,
			s.updated_at,
			s.updated_by,
			s.deleted_at,
			s.deleted_by,
			COALESCE(s.payout_amount, t.payout_amount, b.payout_amount, 0) - COALESCE(s.total_bet_amount, t.bet_amount, b.bet_amount, 0) AS total_win_lose
		%s
		%s
		%s
		%s
	`, baseFromClause, whereClause, sqlOrderBy, sqlLimit)

	if err := r.DBPool.Select(&rows, query, argsFilters...); err != nil {
		custom_log.NewCustomLog("get_statements_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("get_statements_failed", fmt.Errorf("error_database"))
	}

	statements := make([]StatementResponse, 0, len(rows))
	for _, row := range rows {
		statements = append(statements, buildStatementResponse(row))
	}

	return &StatementListResponse{
		Statements:  statements,
		TotalReport: buildTotalReportResponse(totalReport.TotalBet, totalReport.TotalValidBet, totalReport.TotalWinLose),
		Total:       total,
	}, nil
}
