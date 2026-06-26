package jackpot

import (
	"database/sql"
	custom_log "fish_shooting_admin_backend/pkg/logs"
	types "fish_shooting_admin_backend/pkg/model"
	"fish_shooting_admin_backend/pkg/postgres"
	"fish_shooting_admin_backend/pkg/responses"
	"fish_shooting_admin_backend/pkg/utils"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type JackpotRepo interface {
	GetCurrentJackpotGlobal() (*JackpotGlobalResponse, *responses.ErrorResponse)
	UpdateJackpotGlobal(req UpdateJackpotGlobalRequest, currentUser types.UserContext) (*JackpotGlobalResponse, *responses.ErrorResponse)
	GetAllCompanyTopups(req JackpotCompanyTopupShowRequest) (*JackpotCompanyTopupListResponse, *responses.ErrorResponse)
	CreateCompanyTopup(req CreateJackpotCompanyTopupRequest, currentUser types.UserContext) (*JackpotCompanyTopupResponse, *responses.ErrorResponse)
	GetAllLedgers(req JackpotLedgerShowRequest) (*JackpotLedgerListResponse, *responses.ErrorResponse)
	UpdateJackpotRate(req UpdateJackpotRateRequest, currentUser types.UserContext) (*UpdateJackpotRateResponse, *responses.ErrorResponse)
}

type JackpotRepoImpl struct {
	DBPool *sqlx.DB
}

func NewJackpotRepoImpl(dbPool *sqlx.DB) JackpotRepo {
	return &JackpotRepoImpl{
		DBPool: dbPool,
	}
}

func (r *JackpotRepoImpl) GetCurrentJackpotGlobal() (*JackpotGlobalResponse, *responses.ErrorResponse) {
	row := jackpotGlobalRow{}
	query := `
		SELECT
			COALESCE(id, 0) AS id,
			COALESCE(current_amount, 0) AS current_amount,
			COALESCE(threshold_amount, 0) AS threshold_amount,
			COALESCE(chance_denom, 0) AS chance_denom,
			COALESCE(payout_percent, 0) AS payout_percent,
			COALESCE(min_eligible_bet_amount, 0) AS min_eligible_bet_amount,
			COALESCE(company_topup_amount, 0) AS company_topup_amount,
			COALESCE(jackpot_fixed_payout_amount, 0) AS jackpot_fixed_payout_amount,
			COALESCE(status_id, 0) AS status_id,
			"order",
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM tbl_jackpot_global
		WHERE deleted_at IS NULL
		ORDER BY id DESC
		LIMIT 1
	`

	if err := r.DBPool.Get(&row, query); err != nil {
		custom_log.NewCustomLog("get_jackpot_global_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		if err == sql.ErrNoRows {
			return nil, e.NewErrorResponse("not_found", fmt.Errorf("not_found"))
		}
		return nil, e.NewErrorResponse("get_jackpot_global_failed", fmt.Errorf("error_database"))
	}

	resp := buildJackpotGlobalResponse(row)
	return &resp, nil
}

func (r *JackpotRepoImpl) UpdateJackpotGlobal(req UpdateJackpotGlobalRequest, currentUser types.UserContext) (*JackpotGlobalResponse, *responses.ErrorResponse) {
	thresholdAmount, err := decimal.NewFromString(strings.TrimSpace(req.ThresholdAmount))
	if err != nil {
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("invalid_threshold_amount", err)
	}
	if thresholdAmount.LessThanOrEqual(decimal.Zero) {
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("threshold_amount_must_be_greater_than_zero", fmt.Errorf("threshold amount must be greater than zero"))
	}
	if req.ChanceDenom <= 0 {
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("chance_denom_must_be_greater_than_zero", fmt.Errorf("chance denom must be greater than zero"))
	}
	payoutPercent, err := decimal.NewFromString(strings.TrimSpace(req.PayoutPercent))
	if err != nil {
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("invalid_payout_percent", err)
	}
	if payoutPercent.LessThan(decimal.Zero) || payoutPercent.GreaterThan(decimal.NewFromInt(100)) {
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("payout_percent_must_be_between_0_and_100", fmt.Errorf("payout percent must be between 0 and 100"))
	}
	minEligibleBetAmount, err := decimal.NewFromString(strings.TrimSpace(req.MinEligibleBetAmount))
	if err != nil {
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("invalid_min_eligible_bet_amount", err)
	}
	if minEligibleBetAmount.LessThan(decimal.Zero) {
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("min_eligible_bet_amount_must_be_greater_than_or_equal_to_zero", fmt.Errorf("min eligible bet amount must be greater than or equal to zero"))
	}
	jackpotFixedPayoutAmount, err := decimal.NewFromString(strings.TrimSpace(req.JackpotFixedPayoutAmount))
	if err != nil {
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("invalid_jackpot_fixed_payout_amount", err)
	}
	if jackpotFixedPayoutAmount.LessThan(decimal.Zero) {
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("jackpot_fixed_payout_amount_must_be_greater_than_or_equal_to_zero", fmt.Errorf("jackpot fixed payout amount must be greater than or equal to zero"))
	}

	tx, err := r.DBPool.Beginx()
	if err != nil {
		custom_log.NewCustomLog("update_jackpot_global_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("update_jackpot_global_failed", fmt.Errorf("error_database"))
	}
	defer func() {
		if tx != nil {
			_ = tx.Rollback()
		}
	}()

	now, err := utils.GetCurrentAppTime()
	if err != nil {
		custom_log.NewCustomLog("update_jackpot_global_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("update_jackpot_global_failed", fmt.Errorf("error_database"))
	}

	row := jackpotGlobalRow{}
	if err := tx.Get(&row, `
		SELECT
			COALESCE(id, 0) AS id,
			COALESCE(current_amount, 0) AS current_amount,
			COALESCE(threshold_amount, 0) AS threshold_amount,
			COALESCE(chance_denom, 0) AS chance_denom,
			COALESCE(payout_percent, 0) AS payout_percent,
			COALESCE(min_eligible_bet_amount, 0) AS min_eligible_bet_amount,
			COALESCE(company_topup_amount, 0) AS company_topup_amount,
			COALESCE(jackpot_fixed_payout_amount, 0) AS jackpot_fixed_payout_amount,
			COALESCE(status_id, 0) AS status_id,
			"order",
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM tbl_jackpot_global
		WHERE deleted_at IS NULL
		ORDER BY id DESC
		LIMIT 1
		FOR UPDATE
	`); err != nil {
		custom_log.NewCustomLog("update_jackpot_global_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		if err == sql.ErrNoRows {
			return nil, e.NewErrorResponse("not_found", fmt.Errorf("not_found"))
		}
		return nil, e.NewErrorResponse("update_jackpot_global_failed", fmt.Errorf("error_database"))
	}

	if _, err := tx.Exec(`
		UPDATE tbl_jackpot_global
		SET threshold_amount = $1,
			chance_denom = $2,
			payout_percent = $3,
			min_eligible_bet_amount = $4,
			jackpot_fixed_payout_amount = $5,
			updated_at = $6,
			updated_by = $7
		WHERE id = $8
	`, thresholdAmount, req.ChanceDenom, payoutPercent, minEligibleBetAmount, jackpotFixedPayoutAmount, now, currentUser.Id, row.ID); err != nil {
		custom_log.NewCustomLog("update_jackpot_global_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("update_jackpot_global_failed", fmt.Errorf("error_database"))
	}

	if err := tx.Commit(); err != nil {
		custom_log.NewCustomLog("update_jackpot_global_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("update_jackpot_global_failed", fmt.Errorf("error_database"))
	}
	tx = nil
	row.ThresholdAmount = thresholdAmount
	row.ChanceDenom = req.ChanceDenom
	row.PayoutPercent = payoutPercent
	row.MinEligibleBetAmount = minEligibleBetAmount
	row.JackpotFixedPayoutAmount = jackpotFixedPayoutAmount
	row.UpdatedAt = &now
	row.UpdatedBy = &currentUser.Id

	resp := buildJackpotGlobalResponse(row)
	return &resp, nil
}

func (r *JackpotRepoImpl) GetAllCompanyTopups(req JackpotCompanyTopupShowRequest) (*JackpotCompanyTopupListResponse, *responses.ErrorResponse) {
	sqlLimit := postgres.BuildPaging(req.PageOptions.Page, req.PageOptions.Perpage)
	sqlFilters, argsFilters := postgres.BuildSQLFilter(req.Filters)
	sqlOrderBy := postgres.BuildSort(req.Sorts)
	if sqlOrderBy == "" {
		sqlOrderBy = "ORDER BY t.created_at DESC, t.id DESC"
	}

	whereClause := `
		WHERE t.deleted_at IS NULL
	`
	if sqlFilters != "" {
		whereClause += " AND " + sqlFilters
	}

	baseFromClause := `
		FROM tbl_jackpot_company_topup t
		LEFT JOIN tbl_jackpot_global g
			ON g.id = t.jackpot_global_id
			AND g.deleted_at IS NULL
		LEFT JOIN tbl_users u
			ON u.id = t.created_by
			AND u.deleted_at IS NULL
	`

	total := 0
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		%s
		%s
	`, baseFromClause, whereClause)
	if err := r.DBPool.Get(&total, countQuery, argsFilters...); err != nil {
		custom_log.NewCustomLog("get_jackpot_company_topups_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("get_jackpot_company_topups_failed", fmt.Errorf("error_database"))
	}

	rows := make([]jackpotCompanyTopupRow, 0)
	query := fmt.Sprintf(`
		SELECT
			COALESCE(t.id, 0) AS id,
			COALESCE(u.user_name, '') AS username,
			COALESCE(t.amount, 0) AS amount,
			COALESCE(t.current_amount_before, 0) AS current_amount_before,
			COALESCE(t.current_amount_after, 0) AS current_amount_after,
			COALESCE(t.note, '') AS note,
			COALESCE(t."order", 0) AS "order",
			COALESCE(t.created_at, NOW()) AS created_at,
			t.created_by,
			t.updated_at,
			t.updated_by,
			t.deleted_at,
			t.deleted_by
		%s
		%s
		%s
		%s
	`, baseFromClause, whereClause, sqlOrderBy, sqlLimit)

	if err := r.DBPool.Select(&rows, query, argsFilters...); err != nil {
		custom_log.NewCustomLog("get_jackpot_company_topups_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("get_jackpot_company_topups_failed", fmt.Errorf("error_database"))
	}

	topups := make([]JackpotCompanyTopupResponse, 0, len(rows))
	for _, row := range rows {
		topups = append(topups, buildJackpotCompanyTopupResponse(row))
	}

	return &JackpotCompanyTopupListResponse{
		Topups: topups,
		Total:  total,
	}, nil
}

func (r *JackpotRepoImpl) CreateCompanyTopup(req CreateJackpotCompanyTopupRequest, currentUser types.UserContext) (*JackpotCompanyTopupResponse, *responses.ErrorResponse) {
	amount, err := decimal.NewFromString(strings.TrimSpace(req.Amount))
	if err != nil {
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("invalid_amount", err)
	}
	if !amount.GreaterThan(decimal.Zero) {
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("amount_must_be_greater_than_zero", fmt.Errorf("amount must be greater than zero"))
	}

	tx, err := r.DBPool.Beginx()
	if err != nil {
		custom_log.NewCustomLog("create_jackpot_company_topup_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("create_jackpot_company_topup_failed", fmt.Errorf("error_database"))
	}
	defer func() {
		if tx != nil {
			_ = tx.Rollback()
		}
	}()

	now, err := utils.GetCurrentAppTime()
	if err != nil {
		custom_log.NewCustomLog("create_jackpot_company_topup_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("create_jackpot_company_topup_failed", fmt.Errorf("error_database"))
	}

	globalRow := jackpotGlobalRow{}
	if err := tx.Get(&globalRow, `
		SELECT
			COALESCE(id, 0) AS id,
			COALESCE(current_amount, 0) AS current_amount,
			COALESCE(threshold_amount, 0) AS threshold_amount,
			COALESCE(chance_denom, 0) AS chance_denom,
			COALESCE(payout_percent, 0) AS payout_percent,
			COALESCE(min_eligible_bet_amount, 0) AS min_eligible_bet_amount,
			COALESCE(company_topup_amount, 0) AS company_topup_amount,
			COALESCE(jackpot_fixed_payout_amount, 0) AS jackpot_fixed_payout_amount,
			COALESCE(status_id, 0) AS status_id,
			"order",
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM tbl_jackpot_global
		WHERE deleted_at IS NULL
		ORDER BY id DESC
		LIMIT 1
		FOR UPDATE
	`); err != nil {
		custom_log.NewCustomLog("create_jackpot_company_topup_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		if err == sql.ErrNoRows {
			return nil, e.NewErrorResponse("not_found", fmt.Errorf("not_found"))
		}
		return nil, e.NewErrorResponse("create_jackpot_company_topup_failed", fmt.Errorf("error_database"))
	}

	currentBefore := globalRow.CurrentAmount
	currentAfter := currentBefore.Add(amount)

	topupID := int64(0)
	if err := tx.QueryRowx(`
		INSERT INTO tbl_jackpot_company_topup (
			jackpot_global_id,
			amount,
			current_amount_before,
			current_amount_after,
			note,
			"order",
			created_at,
			created_by
		) VALUES ($1, $2, $3, $4, $5, 1, $6, $7)
		RETURNING id
	`, globalRow.ID, amount, currentBefore, currentAfter, req.Note, now, currentUser.Id).Scan(&topupID); err != nil {
		custom_log.NewCustomLog("create_jackpot_company_topup_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("create_jackpot_company_topup_failed", fmt.Errorf("error_database"))
	}

	if _, err := tx.Exec(`
		UPDATE tbl_jackpot_global
		SET current_amount = $1,
			company_topup_amount = COALESCE(company_topup_amount, 0) + $2,
			updated_at = $3,
			updated_by = $4
		WHERE id = $5
	`, currentAfter, amount, now, currentUser.Id, globalRow.ID); err != nil {
		custom_log.NewCustomLog("create_jackpot_company_topup_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("create_jackpot_company_topup_failed", fmt.Errorf("error_database"))
	}

	if err := tx.Commit(); err != nil {
		custom_log.NewCustomLog("create_jackpot_company_topup_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("create_jackpot_company_topup_failed", fmt.Errorf("error_database"))
	}
	tx = nil

	response := JackpotCompanyTopupResponse{
		ID:                  topupID,
		Username:            currentUser.UserName,
		Amount:              amount.Round(2).StringFixed(2),
		CurrentAmountBefore: currentBefore.Round(2).StringFixed(2),
		CurrentAmountAfter:  currentAfter.Round(2).StringFixed(2),
		Note:                req.Note,
		Order:               1,
		CreatedBy:           &currentUser.Id,
	}

	return &response, nil
}

func (r *JackpotRepoImpl) GetAllLedgers(req JackpotLedgerShowRequest) (*JackpotLedgerListResponse, *responses.ErrorResponse) {
	sqlLimit := postgres.BuildPaging(req.PageOptions.Page, req.PageOptions.Perpage)
	sqlFilters, argsFilters := postgres.BuildSQLFilter(req.Filters)
	sqlOrderBy := postgres.BuildSort(req.Sorts)
	if sqlOrderBy == "" {
		sqlOrderBy = "ORDER BY l.created_at DESC, l.id DESC"
	}

	whereClause := `
		WHERE l.deleted_at IS NULL
	`
	if sqlFilters != "" {
		whereClause += " AND " + sqlFilters
	}

	baseFromClause := `
		FROM tbl_jackpot_ledger l
		LEFT JOIN tbl_fish_types ft
			ON ft.id = l.fish_type_id
			AND ft.deleted_at IS NULL
	`

	total := 0
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		%s
		%s
	`, baseFromClause, whereClause)
	if err := r.DBPool.Get(&total, countQuery, argsFilters...); err != nil {
		custom_log.NewCustomLog("get_jackpot_ledger_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("get_jackpot_ledger_failed", fmt.Errorf("error_database"))
	}

	rows := make([]jackpotLedgerRow, 0)
	query := fmt.Sprintf(`
		SELECT
			COALESCE(l.id, 0) AS id,
			COALESCE(l.jackpot_global_id, 0) AS jackpot_global_id,
			COALESCE(l.fish_type_id, 0) AS fish_type_id,
			COALESCE(ft.fish_type_name, '') AS fish_type_name,
			l.member_id,
			l.bet_id,
			l.ticket_id,
			l.statement_id,
			COALESCE(l.source_type, '') AS source_type,
			COALESCE(l.global_contribution_coin, 0) AS global_contribution_coin,
			COALESCE(l.pool_before, 0) AS pool_before,
			COALESCE(l.pool_after, 0) AS pool_after,
			COALESCE(l.created_at, NOW()) AS created_at,
			l.created_by,
			l.updated_at,
			l.updated_by,
			l.deleted_at,
			l.deleted_by
		%s
		%s
		%s
		%s
	`, baseFromClause, whereClause, sqlOrderBy, sqlLimit)

	if err := r.DBPool.Select(&rows, query, argsFilters...); err != nil {
		custom_log.NewCustomLog("get_jackpot_ledger_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("get_jackpot_ledger_failed", fmt.Errorf("error_database"))
	}

	ledgers := make([]JackpotLedgerResponse, 0, len(rows))
	for _, row := range rows {
		ledgers = append(ledgers, buildJackpotLedgerResponse(row))
	}

	return &JackpotLedgerListResponse{
		Ledgers: ledgers,
		Total:   total,
	}, nil
}

func (r *JackpotRepoImpl) UpdateJackpotRate(req UpdateJackpotRateRequest, currentUser types.UserContext) (*UpdateJackpotRateResponse, *responses.ErrorResponse) {
	jackpotRate, err := decimal.NewFromString(strings.TrimSpace(req.JackpotRate))
	if err != nil {
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("invalid_jackpot_rate", err)
	}
	if jackpotRate.LessThanOrEqual(decimal.Zero) || jackpotRate.GreaterThan(decimal.NewFromInt(1)) {
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("jackpot_rate_must_be_between_0_and_1", fmt.Errorf("jackpot rate must be between 0 and 1"))
	}

	tx, err := r.DBPool.Beginx()
	if err != nil {
		custom_log.NewCustomLog("update_jackpot_rate_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("update_jackpot_rate_failed", fmt.Errorf("error_database"))
	}
	defer func() {
		if tx != nil {
			_ = tx.Rollback()
		}
	}()

	row := jackpotGameConfigRow{}
	if err := tx.Get(&row, `
		SELECT
			id,
			game_code,
			COALESCE(jackpot_rate, 0) AS jackpot_rate,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM tbl_game_configs
		WHERE game_code = $1
			AND deleted_at IS NULL
		FOR UPDATE
	`, req.GameCode); err != nil {
		custom_log.NewCustomLog("update_jackpot_rate_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
	if err == sql.ErrNoRows {
			return nil, e.NewErrorResponse("not_found", fmt.Errorf("not_found"))
		}
		return nil, e.NewErrorResponse("update_jackpot_rate_failed", fmt.Errorf("error_database"))
	}

	now, err := utils.GetCurrentAppTime()
	if err != nil {
		custom_log.NewCustomLog("update_jackpot_rate_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("update_jackpot_rate_failed", fmt.Errorf("error_database"))
	}

	if _, err := tx.Exec(`
		UPDATE tbl_game_configs
		SET jackpot_rate = $1,
			updated_at = $2,
			updated_by = $3
		WHERE id = $4
	`, jackpotRate, now, currentUser.Id, row.ID); err != nil {
		custom_log.NewCustomLog("update_jackpot_rate_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("update_jackpot_rate_failed", fmt.Errorf("error_database"))
	}

	if err := tx.Commit(); err != nil {
		custom_log.NewCustomLog("update_jackpot_rate_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("update_jackpot_rate_failed", fmt.Errorf("error_database"))
	}
	tx = nil

	return &UpdateJackpotRateResponse{
		ID:         row.ID,
		GameCode:   row.GameCode,
		JackpotRate: jackpotRate.StringFixed(6),
		UpdatedAt:  now,
		UpdatedBy:  &currentUser.Id,
	}, nil
}
