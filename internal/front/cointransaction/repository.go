package cointransaction

import (
	custom_log "fish_shooting_admin_backend/pkg/logs"
	"fish_shooting_admin_backend/pkg/postgres"
	"fish_shooting_admin_backend/pkg/responses"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type CoinTransactionRepo interface {
	GetTransactions(req CoinTransactionShowRequest) (*CoinTransactionListResponse, *responses.ErrorResponse)
}

type CoinTransactionRepoImpl struct {
	DBPool *sqlx.DB
}

func NewCoinTransactionRepoImpl(dbPool *sqlx.DB) CoinTransactionRepo {
	return &CoinTransactionRepoImpl{
		DBPool: dbPool,
	}
}

func (r *CoinTransactionRepoImpl) GetTransactions(req CoinTransactionShowRequest) (*CoinTransactionListResponse, *responses.ErrorResponse) {
	sqlLimit := postgres.BuildPaging(req.PageOptions.Page, req.PageOptions.Perpage)
	sqlFilters, argsFilters := postgres.BuildSQLFilter(req.Filters)
	sqlOrderBy := postgres.BuildSort(req.Sorts)
	if sqlOrderBy == "" {
		sqlOrderBy = "ORDER BY COALESCE(ct.transaction_date, ct.created_at) DESC, ct.id DESC"
	}

	whereClause := `
		WHERE ct.deleted_at IS NULL
	`
	if sqlFilters != "" {
		whereClause += " AND " + sqlFilters
	}

	baseFromClause := `
		FROM tbl_members_coins_transactions ct
		LEFT JOIN tbl_members m
			ON m.id = ct.member_id
			AND m.deleted_at IS NULL
	`

	total := 0
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		%s
		%s
	`, baseFromClause, whereClause)

	if err := r.DBPool.Get(&total, countQuery, argsFilters...); err != nil {
		custom_log.NewCustomLog("get_coin_transactions_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("get_coin_transactions_failed", fmt.Errorf("error_database"))
	}

	rows := make([]coinTransactionRow, 0)
	query := fmt.Sprintf(`
		SELECT
			COALESCE(ct.member_id, 0) AS member_id,
			COALESCE(m.user_uuid::text, '') AS member_uuid,
			COALESCE(m.user_name, '') AS username,
			COALESCE(ct.member_coin_id, 0) AS member_coin_id,
			COALESCE(ct.before_coin, 0) AS before_coin,
			COALESCE(ct.amount, 0) AS amount,
			COALESCE(ct.transaction_type_id, 0) AS transaction_type_id,
			COALESCE(ct.transaction_group_type_id, 0) AS transaction_group_type_id,
			COALESCE(ct.transaction_date, ct.created_at) AS transaction_date,
			COALESCE(ct.require_approval, FALSE) AS require_approval,
			COALESCE(ct.reference, '') AS reference,
			COALESCE(ct.remark, '') AS remark,
			COALESCE(ct.status_id, 0) AS status_id,
			ct."order" AS "order",
			COALESCE(ct.created_at, NOW()) AS created_at
		%s
		%s
		%s
		%s
	`, baseFromClause, whereClause, sqlOrderBy, sqlLimit)

	if err := r.DBPool.Select(&rows, query, argsFilters...); err != nil {
		custom_log.NewCustomLog("get_coin_transactions_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("get_coin_transactions_failed", fmt.Errorf("error_database"))
	}

	transactions := make([]CoinTransactionResponse, 0, len(rows))
	for _, row := range rows {
		transactions = append(transactions, buildCoinTransactionResponse(row))
	}

	return &CoinTransactionListResponse{
		Transactions: transactions,
		Total:        total,
	}, nil
}
