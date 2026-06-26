package memberbonus

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

type JackpotMemberBonusRepo interface {
	GetAllMemberBonuses(req JackpotMemberBonusShowRequest) (*JackpotMemberBonusListResponse, *responses.ErrorResponse)
	CreateMemberBonus(req CreateJackpotMemberBonusRequest, currentUser types.UserContext) (*JackpotMemberBonusResponse, *responses.ErrorResponse)
}

type JackpotMemberBonusRepoImpl struct {
	DBPool *sqlx.DB
}

func NewJackpotMemberBonusRepoImpl(dbPool *sqlx.DB) JackpotMemberBonusRepo {
	return &JackpotMemberBonusRepoImpl{DBPool: dbPool}
}

func (r *JackpotMemberBonusRepoImpl) GetAllMemberBonuses(req JackpotMemberBonusShowRequest) (*JackpotMemberBonusListResponse, *responses.ErrorResponse) {
	sqlLimit := postgres.BuildPaging(req.PageOptions.Page, req.PageOptions.Perpage)
	sqlFilters, argsFilters := postgres.BuildSQLFilter(req.Filters)
	sqlOrderBy := postgres.BuildSort(req.Sorts)
	if sqlOrderBy == "" {
		sqlOrderBy = "ORDER BY b.created_at DESC, b.id DESC"
	}

	whereClause := `
		WHERE b.deleted_at IS NULL
	`
	if sqlFilters != "" {
		whereClause += " AND " + sqlFilters
	}

	baseFromClause := `
		FROM tbl_jackpot_member_bonus b
		LEFT JOIN tbl_members m
			ON m.id = b.member_id
			AND m.deleted_at IS NULL
		LEFT JOIN tbl_users u
			ON u.id = b.created_by
			AND u.deleted_at IS NULL
	`

	total := 0
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		%s
		%s
	`, baseFromClause, whereClause)
	if err := r.DBPool.Get(&total, countQuery, argsFilters...); err != nil {
		custom_log.NewCustomLog("get_jackpot_member_bonuses_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("get_jackpot_member_bonuses_failed", fmt.Errorf("error_database"))
	}

	rows := make([]jackpotMemberBonusRow, 0)
	query := fmt.Sprintf(`
		SELECT
			COALESCE(b.id, 0) AS id,
			COALESCE(b.member_id, 0) AS member_id,
			COALESCE(m.user_name, '') AS member_name,
			COALESCE(m.user_uuid::text, '') AS member_uuid,
			COALESCE(u.user_name, '') AS created_by_name,
			COALESCE(b.amount, 0) AS amount,
			COALESCE(b.note, '') AS note,
			COALESCE(b."order", 0) AS "order",
			COALESCE(b.status_id, 1) AS status_id,
			COALESCE(b.created_at, NOW()) AS created_at,
			b.created_by,
			b.updated_at,
			b.updated_by,
			b.deleted_at,
			b.deleted_by
		%s
		%s
		%s
		%s
	`, baseFromClause, whereClause, sqlOrderBy, sqlLimit)

	if err := r.DBPool.Select(&rows, query, argsFilters...); err != nil {
		custom_log.NewCustomLog("get_jackpot_member_bonuses_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("get_jackpot_member_bonuses_failed", fmt.Errorf("error_database"))
	}

	bonuses := make([]JackpotMemberBonusResponse, 0, len(rows))
	for _, row := range rows {
		bonuses = append(bonuses, buildJackpotMemberBonusResponse(row))
	}

	return &JackpotMemberBonusListResponse{
		Bonuses: bonuses,
		Total:   total,
	}, nil
}

func (r *JackpotMemberBonusRepoImpl) CreateMemberBonus(req CreateJackpotMemberBonusRequest, currentUser types.UserContext) (*JackpotMemberBonusResponse, *responses.ErrorResponse) {
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
		custom_log.NewCustomLog("create_jackpot_member_bonus_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("create_jackpot_member_bonus_failed", fmt.Errorf("error_database"))
	}
	defer func() {
		if tx != nil {
			_ = tx.Rollback()
		}
	}()

	now, err := utils.GetCurrentAppTime()
	if err != nil {
		custom_log.NewCustomLog("create_jackpot_member_bonus_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("create_jackpot_member_bonus_failed", fmt.Errorf("error_database"))
	}

	memberRow := struct {
		MemberName string `db:"member_name"`
	}{}
	if err := tx.Get(&memberRow, `
		SELECT COALESCE(user_name, '') AS member_name
		FROM tbl_members
		WHERE id = $1 AND deleted_at IS NULL
		FOR UPDATE
	`, req.MemberID); err != nil {
		custom_log.NewCustomLog("create_jackpot_member_bonus_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		if err == sql.ErrNoRows {
			return nil, e.NewErrorResponse("member_not_found", fmt.Errorf("member_not_found"))
		}
		return nil, e.NewErrorResponse("create_jackpot_member_bonus_failed", fmt.Errorf("error_database"))
	}

	bonusID := int64(0)
	if err := tx.QueryRowx(`
		INSERT INTO tbl_jackpot_member_bonus (
			member_id,
			amount,
			note,
			"order",
			status_id,
			created_at,
			created_by
		) VALUES ($1, $2, $3, 1, 1, $4, $5)
		RETURNING id
	`, req.MemberID, amount, req.Note, now, currentUser.Id).Scan(&bonusID); err != nil {
		custom_log.NewCustomLog("create_jackpot_member_bonus_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("create_jackpot_member_bonus_failed", fmt.Errorf("error_database"))
	}

	if err := tx.Commit(); err != nil {
		custom_log.NewCustomLog("create_jackpot_member_bonus_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("create_jackpot_member_bonus_failed", fmt.Errorf("error_database"))
	}
	tx = nil

	resp := JackpotMemberBonusResponse{
		ID:            bonusID,
		MemberID:      req.MemberID,
		MemberName:    memberRow.MemberName,
		CreatedByName: currentUser.UserName,
		Amount:        amount.Round(2).StringFixed(2),
		Note:          req.Note,
		Order:         1,
		StatusID:      1,
		CreatedBy:     &currentUser.Id,
	}

	return &resp, nil
}
