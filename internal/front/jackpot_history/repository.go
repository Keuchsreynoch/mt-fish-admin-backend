package jackpothistory

import (
	custom_log "fish_shooting_admin_backend/pkg/logs"
	types "fish_shooting_admin_backend/pkg/model"
	"fish_shooting_admin_backend/pkg/postgres"
	"fish_shooting_admin_backend/pkg/responses"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type JackpotHistoryRepo interface {
	GetAllJackpotHistories(req JackpotHistoryShowRequest) (*JackpotHistoryListResponse, *responses.ErrorResponse)
}

type JackpotHistoryRepoImpl struct {
	DBPool *sqlx.DB
}

func NewJackpotHistoryRepoImpl(dbPool *sqlx.DB) JackpotHistoryRepo {
	return &JackpotHistoryRepoImpl{
		DBPool: dbPool,
	}
}

func (r *JackpotHistoryRepoImpl) GetAllJackpotHistories(req JackpotHistoryShowRequest) (*JackpotHistoryListResponse, *responses.ErrorResponse) {
	sqlLimit := postgres.BuildPaging(req.PageOptions.Page, req.PageOptions.Perpage)

	filters := append([]types.Filter{}, req.Filters...)
	sqlFilters, argsFilters := postgres.BuildSQLFilter(filters)

	sqlOrderBy := postgres.BuildSort(req.Sorts)
	if sqlOrderBy == "" {
		sqlOrderBy = "ORDER BY h.created_at DESC, h.id DESC"
	}

	whereClause := `
		WHERE h.deleted_at IS NULL
	`
	if sqlFilters != "" {
		whereClause += " AND " + sqlFilters
	}

	baseFromClause := `
		FROM tbl_jackpot_history h
		LEFT JOIN tbl_jackpot_global jg
			ON jg.id = h.jackpot_global_id
			AND jg.deleted_at IS NULL
		LEFT JOIN tbl_members m
			ON m.id = h.member_id
			AND m.deleted_at IS NULL
		LEFT JOIN tbl_fish_types ft
			ON ft.id = h.fish_type_id
			AND ft.deleted_at IS NULL
	`

	total := 0
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		%s
		%s
	`, baseFromClause, whereClause)

	if err := r.DBPool.Get(&total, countQuery, argsFilters...); err != nil {
		custom_log.NewCustomLog("get_jackpot_history_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("get_jackpot_history_failed", fmt.Errorf("error_database"))
	}

	rows := make([]jackpotHistoryRow, 0)
	query := fmt.Sprintf(`
		SELECT
			COALESCE(h.id, 0) AS id,
			COALESCE(h.jackpot_global_id, 0) AS jackpot_global_id,
			COALESCE(h.fish_type_id, 0) AS fish_type_id,
			h.member_id,

			COALESCE(h.jackpot_type, 1) AS jackpot_type,
			COALESCE(m.user_name, '') AS member_name,
			COALESCE(ft.fish_type_name, '') AS fish_type_name,
			COALESCE(h.payout_coin, 0) AS payout_coin,
			COALESCE(h.pool_before, 0) AS pool_before,
			COALESCE(h.pool_after, 0) AS pool_after,
			COALESCE(h.created_at, NOW()) AS created_at,
			h.created_by,
			h.updated_at,
			h.updated_by,
			h.deleted_at,
			h.deleted_by
		%s
		%s
		%s
		%s
	`, baseFromClause, whereClause, sqlOrderBy, sqlLimit)

	if err := r.DBPool.Select(&rows, query, argsFilters...); err != nil {
		custom_log.NewCustomLog("get_jackpot_history_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("get_jackpot_history_failed", fmt.Errorf("error_database"))
	}

	histories := make([]JackpotHistoryResponse, 0, len(rows))
	for _, row := range rows {
		histories = append(histories, buildJackpotHistoryResponse(row))
	}

	return &JackpotHistoryListResponse{
		Histories: histories,
		Total:     total,
	}, nil
}
