package member

import (
	custom_log "fish_shooting_admin_backend/pkg/logs"
	"fish_shooting_admin_backend/pkg/postgres"
	"fish_shooting_admin_backend/pkg/responses"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type MemberRepo interface {
	GetAllMembers(req MemberShowRequest) (*GetAllMembersResponse, *responses.ErrorResponse)
}

type MemberRepoImpl struct {
	DBPool *sqlx.DB
}

func NewMemberRepoImpl(dbPool *sqlx.DB) *MemberRepoImpl {
	return &MemberRepoImpl{
		DBPool: dbPool,
	}
}

func (m *MemberRepoImpl) GetAllMembers(req MemberShowRequest) (*GetAllMembersResponse, *responses.ErrorResponse) {
	sqlLimit := postgres.BuildPaging(req.PageOptions.Page, req.PageOptions.Perpage)
	sqlFilters, argsFilters := postgres.BuildSQLFilter(req.Filters)
	sqlOrderBy := postgres.BuildSort(req.Sorts)
	if sqlOrderBy == "" {
		sqlOrderBy = "ORDER BY created_at DESC, id DESC"
	}
	whereClause := `
		WHERE m.deleted_at IS NULL
	`
	if sqlFilters != "" {
		whereClause += " AND " + sqlFilters
	}

	baseQuery := `
		FROM tbl_members m
	`

	total := 0
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		%s
		%s
	`, baseQuery, whereClause)

	if err := m.DBPool.Get(&total, countQuery, argsFilters...); err != nil {
		custom_log.NewCustomLog("get_all_members_failed", err.Error(), "error")
		errResp := &responses.ErrorResponse{}
		return nil, errResp.NewErrorResponse("get_all_members_failed", fmt.Errorf("error_database"))
	}

	var members []MemberInfo
	query := fmt.Sprintf(`
		SELECT
			m.id,
			m.user_uuid,
			m.user_name,
			m.login_id,
			m.phone_number,
			m.profile_photo,
			m.language_id,
			m.currency_id,
			COALESCE(mc.coin_amount, 0)::text AS coin_amount,
			m.remark,
			m.nickname,
			m.login_session,
			m.last_login_at,
			m.is_online,
			m.status_id,
			m.timezone,
			m.pattern,
			m."order",
			m.created_by,
			m.created_at,
			m.updated_by,
			m.updated_at
		%s
		LEFT JOIN LATERAL (
			SELECT coin_amount
			FROM tbl_members_coins
			WHERE member_id = m.id
				AND deleted_at IS NULL
			ORDER BY id DESC
			LIMIT 1
		) mc ON TRUE
		%s
		%s
		%s
	`, baseQuery, whereClause, sqlOrderBy, sqlLimit)

	if err := m.DBPool.Select(&members, query, argsFilters...); err != nil {
		custom_log.NewCustomLog("get_all_members_failed", err.Error(), "error")
		errResp := &responses.ErrorResponse{}
		return nil, errResp.NewErrorResponse("get_all_members_failed", fmt.Errorf("error_database"))
	}

	if len(members) == 0 {
		return &GetAllMembersResponse{
			Members: members,
			Total:   total,
		}, nil
	}

	memberIDs := make([]int, 0, len(members))
	memberIndexByID := make(map[int]int, len(members))
	for idx := range members {
		memberIDs = append(memberIDs, members[idx].ID)
		memberIndexByID[members[idx].ID] = idx
	}

	balanceQuery, balanceArgs, err := sqlx.In(`
		SELECT
			mb.member_id,
			mb.currency_id,
			c.currency_name,
			c.currency_code,
			c.currency_symbol,
			mb.balance
		FROM tbl_members_balances mb
		JOIN tbl_currencies c
			ON c.id = mb.currency_id
			AND c.deleted_at IS NULL
		WHERE mb.deleted_at IS NULL
			AND mb.member_id IN (?)
		ORDER BY mb.member_id ASC, mb.currency_id ASC
	`, memberIDs)
	if err != nil {
		custom_log.NewCustomLog("get_all_members_failed", err.Error(), "error")
		errResp := &responses.ErrorResponse{}
		return nil, errResp.NewErrorResponse("get_all_members_failed", fmt.Errorf("error_database"))
	}

	balanceQuery = m.DBPool.Rebind(balanceQuery)

	type memberBalanceRow struct {
		MemberID       int             `db:"member_id"`
		CurrencyID     int             `db:"currency_id"`
		CurrencyName   string          `db:"currency_name"`
		CurrencyCode   string          `db:"currency_code"`
		CurrencySymbol string          `db:"currency_symbol"`
		Balance        decimal.Decimal `db:"balance"`
	}

	var balances []memberBalanceRow
	if err := m.DBPool.Select(&balances, balanceQuery, balanceArgs...); err != nil {
		custom_log.NewCustomLog("get_all_members_failed", err.Error(), "error")
		errResp := &responses.ErrorResponse{}
		return nil, errResp.NewErrorResponse("get_all_members_failed", fmt.Errorf("error_database"))
	}

	for _, balance := range balances {
		memberIdx, ok := memberIndexByID[balance.MemberID]
		if !ok {
			continue
		}

		members[memberIdx].Balances = append(members[memberIdx].Balances, MemberCurrencyBalance{
			CurrencyID:     balance.CurrencyID,
			CurrencyName:   balance.CurrencyName,
			CurrencyCode:   balance.CurrencyCode,
			CurrencySymbol: balance.CurrencySymbol,
			Balance:        balance.Balance.Round(3).StringFixed(3),
		})
	}

	return &GetAllMembersResponse{
		Members: members,
		Total:   total,
	}, nil
}
