package admin

import (
	"database/sql"
	custom_log "fish_shooting_admin_backend/pkg/logs"
	"fish_shooting_admin_backend/pkg/postgres"
	"fish_shooting_admin_backend/pkg/responses"
	"fish_shooting_admin_backend/pkg/utils"
	"fmt"

	types "fish_shooting_admin_backend/pkg/model"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type AdminRepo interface {
	GetDashboard() (*DashboardSummaryResponse, *responses.ErrorResponse)
	ListUsers() (*UserListResponse, *responses.ErrorResponse)
	ListMenus() (*MenuListResponse, *responses.ErrorResponse)
	GetUserInfo(userUUID string) (*UserInfoResponse, *responses.ErrorResponse)
	GetUserMenus(userID int) (*UserMenusResponse, *responses.ErrorResponse)
	CreateUser(req NewUserRequest) (*NewUser, *responses.ErrorResponse)
	UpdateUser(userID int, req UpdateAdminUserRequest) (*UpdateAdminUserResponse, *responses.ErrorResponse)
	AssignMenus(userID int, req AssignMenusRequest) (*UserMenusResponse, *responses.ErrorResponse)
}

type AdminRepoImpl struct {
	DBPool      *sqlx.DB
	UserContext *types.UserContext
}

func NewAdminRepoImpl(dbPool *sqlx.DB, userContext *types.UserContext) AdminRepo {
	return &AdminRepoImpl{
		DBPool:      dbPool,
		UserContext: userContext,
	}
}

func (r *AdminRepoImpl) GetDashboard() (*DashboardSummaryResponse, *responses.ErrorResponse) {

	baseFromClause := `
		FROM tbl_statements s
		LEFT JOIN tbl_tickets t
			ON t.id = s.ticket_id
			AND t.deleted_at IS NULL
		LEFT JOIN tbl_bets b
			ON b.ticket_id = t.id
			AND b.deleted_at IS NULL
		LEFT JOIN tbl_members m
			ON m.id = COALESCE(s.member_id, t.member_id, b.member_id)
			AND m.deleted_at IS NULL
	`
	todayWhereClause := `
		WHERE s.deleted_at IS NULL
			AND COALESCE(s.statement_at, s.created_at) >= DATE_TRUNC('day', CURRENT_TIMESTAMP)
			AND COALESCE(s.statement_at, s.created_at) < DATE_TRUNC('day', CURRENT_TIMESTAMP) + INTERVAL '1 day'
	`

	dashboardTotals := DashboardTotals{}
	if err := r.DBPool.Get(&dashboardTotals, fmt.Sprintf(`
		SELECT
			COALESCE(SUM(COALESCE(s.total_bet_amount, t.bet_amount, b.bet_amount, 0)), 0) AS total_turnover,
			COALESCE(SUM(COALESCE(s.payout_amount, t.payout_amount, b.payout_amount, 0)), 0) AS total_payout,
			COALESCE(SUM(COALESCE(s.total_bet_amount, t.bet_amount, b.bet_amount, 0)), 0)
				- COALESCE(SUM(COALESCE(s.payout_amount, t.payout_amount, b.payout_amount, 0)), 0) AS total_company_profit
		%s
		%s
	`, baseFromClause, todayWhereClause)); err != nil {
		custom_log.NewCustomLog("get_admin_dashboard_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("get_admin_dashboard_failed", fmt.Errorf("error_database"))
	}

	currentJackpot := struct {
		CurrentPoolJackpot decimal.Decimal `db:"current_pool_jackpot"`
		ThresholdAmount    decimal.Decimal `db:"threshold_amount"`
	}{}
	if err := r.DBPool.Get(&currentJackpot, `
		SELECT COALESCE(current_amount, 0) AS current_pool_jackpot
			, COALESCE(threshold_amount, 0) AS threshold_amount
		FROM tbl_jackpot_global
		WHERE deleted_at IS NULL
		ORDER BY id DESC
		LIMIT 1
	`); err != nil {
		e := &responses.ErrorResponse{}
		if err == sql.ErrNoRows {
			currentJackpot.CurrentPoolJackpot = decimal.Zero
		} else {
			custom_log.NewCustomLog("get_admin_dashboard_failed", err.Error(), "error")
			return nil, e.NewErrorResponse("get_admin_dashboard_failed", fmt.Errorf("error_database"))
		}
	}

	rows := make([]dashboardTopWinMemberRow, 0)
	if err := r.DBPool.Select(&rows, fmt.Sprintf(`
		SELECT
			COALESCE(m.id, 0) AS member_id,
			COALESCE(m.user_uuid::text, '') AS member_uuid,
			COALESCE(m.user_name, '') AS username,
			COALESCE(SUM(COALESCE(s.total_bet_amount, t.bet_amount, b.bet_amount, 0)), 0) AS total_bet_amount,
			COALESCE(SUM(COALESCE(s.payout_amount, t.payout_amount, b.payout_amount, 0)), 0) AS total_payout_amount,
			COALESCE(SUM(GREATEST(COALESCE(s.payout_amount, t.payout_amount, b.payout_amount, 0) - COALESCE(s.total_bet_amount, t.bet_amount, b.bet_amount, 0), 0)), 0) AS total_win_amount,
			COALESCE(COUNT(*) FILTER (
				WHERE COALESCE(s.payout_amount, t.payout_amount, b.payout_amount, 0) > COALESCE(s.total_bet_amount, t.bet_amount, b.bet_amount, 0)
			), 0) AS win_count
		%s
		%s
		GROUP BY COALESCE(m.id, 0), COALESCE(m.user_uuid::text, ''), COALESCE(m.user_name, '')
		ORDER BY total_win_amount DESC, total_payout_amount DESC, member_id DESC
		LIMIT 5
	`, baseFromClause, todayWhereClause)); err != nil {
		custom_log.NewCustomLog("get_admin_dashboard_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("get_admin_dashboard_failed", fmt.Errorf("error_database"))
	}

	topMembers := make([]DashboardTopWinMemberResponse, 0, len(rows))
	for _, row := range rows {
		topMembers = append(topMembers, DashboardTopWinMemberResponse{
			MemberID:          row.MemberID,
			MemberUUID:        row.MemberUUID,
			Username:          row.Username,
			TotalBetAmount:    row.TotalBetAmount.Round(3).StringFixed(3),
			TotalPayoutAmount: row.TotalPayoutAmount.Round(3).StringFixed(3),
			TotalWinAmount:    row.TotalWinAmount.Round(3).StringFixed(3),
			WinCount:          row.WinCount,
		})
	}

	return &DashboardSummaryResponse{
		TotalTurnover:      dashboardTotals.TotalTurnover.Round(3).StringFixed(3),
		TotalPayout:        dashboardTotals.TotalPayout.Round(3).StringFixed(3),
		TotalCompanyProfit: dashboardTotals.TotalCompanyProfit.Round(3).StringFixed(3),
		CurrentPoolJackpot: currentJackpot.CurrentPoolJackpot.Round(3).StringFixed(3),
		ThresholdAmount:    currentJackpot.ThresholdAmount.Round(3).StringFixed(3),
		TopWinMembers:      topMembers,
	}, nil
}

func (r *AdminRepoImpl) ListUsers() (*UserListResponse, *responses.ErrorResponse) {
	rows := make([]User, 0)
	query := `
		SELECT
			id,
			user_uuid,
			user_name,
			login_id,
			COALESCE(email, '') AS email,
			COALESCE(password, '') AS password,
			COALESCE(nickname, '') AS nickname,
			COALESCE(profile, '') AS profile,
			COALESCE(role_id, 0) AS role_id,
			COALESCE(role_name, '') AS role_name,
			COALESCE(is_admin, FALSE) AS is_admin,
			COALESCE(login_session, '') AS login_session,
			COALESCE(status_id, 0) AS status_id,
			COALESCE("order", 0) AS "order",
			COALESCE(created_by, 0) AS created_by,
			COALESCE(created_at, NOW()) AS created_at,
			COALESCE(updated_by, created_by, 0) AS updated_by,
			COALESCE(updated_at, created_at, NOW()) AS updated_at
		FROM tbl_users
		WHERE deleted_at IS NULL
			AND id <> $1
		ORDER BY id ASC
	`

	if err := r.DBPool.Select(&rows, query, r.UserContext.Id); err != nil {
		custom_log.NewCustomLog("get_admin_users_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("get_admin_users_failed", fmt.Errorf("error_database"))
	}

	return &UserListResponse{Users: rows}, nil
}

func (r *AdminRepoImpl) ListMenus() (*MenuListResponse, *responses.ErrorResponse) {
	rows := make([]UserMenu, 0)
	query := `
		SELECT
			id,
			menu_uuid,
			name,
			icon,
			COALESCE(path, '') AS path,
			COALESCE(parent_id, 0) AS parent_id,
			COALESCE(status_id, 0) AS status_id,
			COALESCE("order", 0) AS "order",
			COALESCE(created_by, 0) AS created_by,
			COALESCE(created_at, NOW()) AS created_at,
			COALESCE(updated_by, created_by, 0) AS updated_by,
			COALESCE(updated_at, created_at, NOW()) AS updated_at,
			deleted_by,
			deleted_at
		FROM tbl_menus
		WHERE deleted_at IS NULL
		ORDER BY parent_id ASC, "order" ASC, id ASC
	`

	if err := r.DBPool.Select(&rows, query); err != nil {
		custom_log.NewCustomLog("get_admin_menus_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("get_admin_menus_failed", fmt.Errorf("error_database"))
	}

	return &MenuListResponse{Menus: rows}, nil
}

func (r *AdminRepoImpl) GetUserInfo(userUUID string) (*UserInfoResponse, *responses.ErrorResponse) {
	user := User{}
	userQuery := `
		SELECT
			id,
			user_uuid,
			user_name,
			login_id,
			COALESCE(email, '') AS email,
			COALESCE(password, '') AS password,
			COALESCE(nickname, '') AS nickname,
			COALESCE(profile, '') AS profile,
			COALESCE(role_id, 0) AS role_id,
			COALESCE(role_name, '') AS role_name,
			COALESCE(is_admin, FALSE) AS is_admin,
			COALESCE(login_session, '') AS login_session,
			COALESCE(status_id, 0) AS status_id,
			COALESCE("order", 0) AS "order",
			COALESCE(created_by, 0) AS created_by,
			COALESCE(created_at, NOW()) AS created_at,
			COALESCE(updated_by, created_by, 0) AS updated_by,
			COALESCE(updated_at, created_at, NOW()) AS updated_at
		FROM tbl_users
		WHERE deleted_at IS NULL
			AND user_uuid = $1
		LIMIT 1
	`

	if err := r.DBPool.Get(&user, userQuery, userUUID); err != nil {
		custom_log.NewCustomLog("get_admin_user_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("get_userinfo_failed", fmt.Errorf("error_database"))
	}

	menus, errResp := r.GetUserMenus(user.ID)
	if errResp != nil {
		return nil, errResp
	}

	return &UserInfoResponse{
		User:  user,
		Menus: menus.Menus,
	}, nil
}

func (r *AdminRepoImpl) GetUserMenus(userID int) (*UserMenusResponse, *responses.ErrorResponse) {
	rows := make([]UserMenu, 0)
	query := `
		SELECT
			m.id,
			m.menu_uuid,
			m.name,
			m.icon,
			COALESCE(m.path, '') AS path,
			COALESCE(m.parent_id, 0) AS parent_id,
			COALESCE(m.status_id, 0) AS status_id,
			COALESCE(m."order", 0) AS "order",
			COALESCE(m.created_by, 0) AS created_by,
			COALESCE(m.created_at, NOW()) AS created_at,
			COALESCE(m.updated_by, m.created_by, 0) AS updated_by,
			COALESCE(m.updated_at, m.created_at, NOW()) AS updated_at,
			m.deleted_by,
			m.deleted_at
		FROM tbl_users_menus um
		JOIN tbl_menus m
			ON m.id = um.menu_id
		WHERE um.deleted_at IS NULL
			AND um.is_allow = TRUE
			AND m.deleted_at IS NULL
			AND um.user_id = $1
		ORDER BY m.parent_id ASC, m."order" ASC, m.id ASC
	`

	if err := r.DBPool.Select(&rows, query, userID); err != nil {
		custom_log.NewCustomLog("get_user_menus_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("get_user_menus_failed", fmt.Errorf("error_database"))
	}

	return &UserMenusResponse{Menus: rows}, nil
}

func (r *AdminRepoImpl) CreateUser(req NewUserRequest) (*NewUser, *responses.ErrorResponse) {
	tx, err := r.DBPool.Beginx()
	if err != nil {
		custom_log.NewCustomLog("create_admin_user_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("create_admin_user_failed", fmt.Errorf("error_database"))
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var newUser NewUser
	errResp := newUser.New(req, r.UserContext, tx)
	if errResp != nil {
		return nil, errResp
	}

	insertQuery := `
		INSERT INTO tbl_users (
			user_uuid,
			user_name,
			login_id,
			email,
			password,
			nickname,
			profile,
			role_id,
			role_name,
			is_admin,
			status_id,
			"order",
			created_by,
			created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			1, 'ADMIN', TRUE, 1, 0, $8, $9
		)
		RETURNING id
	`

	var createdID int
	if err := tx.Get(
		&createdID,
		insertQuery,
		newUser.User.UserUUID,
		newUser.User.UserName,
		newUser.User.LoginID,
		newUser.User.Email,
		req.Password,
		newUser.User.Nickname,
		newUser.User.Profile,
		r.UserContext.Id,
		newUser.User.CreatedAt,
	); err != nil {
		custom_log.NewCustomLog("create_admin_user_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("create_admin_user_failed", fmt.Errorf("error_database"))
	}

	if _, err := tx.Exec(`UPDATE tbl_users SET "order" = id WHERE id = $1`, createdID); err != nil {
		custom_log.NewCustomLog("create_admin_user_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("create_admin_user_failed", fmt.Errorf("error_database"))
	}

	if err := syncUserMenusTx(tx, createdID, nil, r.UserContext.Id, true); err != nil {
		custom_log.NewCustomLog("create_admin_user_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("create_admin_user_failed", fmt.Errorf("error_database"))
	}

	if err := tx.Commit(); err != nil {
		custom_log.NewCustomLog("create_admin_user_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("create_admin_user_failed", fmt.Errorf("error_database"))
	}

	user, errResp := r.getUserByID(createdID)
	if errResp != nil {
		return nil, errResp
	}

	menus, errResp := r.GetUserMenus(createdID)
	if errResp != nil {
		return nil, errResp
	}

	return &NewUser{
		User:  *user,
		Menus: menus.Menus,
	}, nil
}

func (r *AdminRepoImpl) UpdateUser(userID int, req UpdateAdminUserRequest) (*UpdateAdminUserResponse, *responses.ErrorResponse) {
	exists, err := postgres.IsExists("tbl_users", "id", userID, r.DBPool)
	if err != nil {
		custom_log.NewCustomLog("update_admin_user_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("update_admin_user_failed", fmt.Errorf("error_database"))
	}
	if !exists {
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("update_admin_user_failed", fmt.Errorf("user_not_found"))
	}

	tx, err := r.DBPool.Beginx()
	if err != nil {
		custom_log.NewCustomLog("update_admin_user_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("update_admin_user_failed", fmt.Errorf("error_database"))
	}
	defer func() {
		_ = tx.Rollback()
	}()

	now, err := utils.GetCurrentAppTime()
	if err != nil {
		custom_log.NewCustomLog("update_admin_user_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("update_admin_user_failed", fmt.Errorf("error_database"))
	}

	updateQuery := `
		UPDATE tbl_users SET
			password = $1,
			updated_by = $2,
			updated_at = $3
		WHERE deleted_at IS NULL
			AND id = $4
	`

	if _, err := tx.Exec(
		updateQuery,
		req.Password,
		r.UserContext.Id,
		now,
		userID,
	); err != nil {
		custom_log.NewCustomLog("update_admin_user_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("update_admin_user_failed", fmt.Errorf("error_database"))
	}

	if err := tx.Commit(); err != nil {
		custom_log.NewCustomLog("update_admin_user_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("update_admin_user_failed", fmt.Errorf("error_database"))
	}

	user, errResp := r.getUserByID(userID)
	if errResp != nil {
		return nil, errResp
	}

	menus, errResp := r.GetUserMenus(userID)
	if errResp != nil {
		return nil, errResp
	}

	return &UpdateAdminUserResponse{
		User:  *user,
		Menus: menus.Menus,
	}, nil
}

func (r *AdminRepoImpl) AssignMenus(userID int, req AssignMenusRequest) (*UserMenusResponse, *responses.ErrorResponse) {
	exists, err := postgres.IsExists("tbl_users", "id", userID, r.DBPool)
	if err != nil {
		custom_log.NewCustomLog("assign_menus_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("assign_menus_failed", fmt.Errorf("error_database"))
	}
	if !exists {
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("assign_menus_failed", fmt.Errorf("user_not_found"))
	}

	tx, err := r.DBPool.Beginx()
	if err != nil {
		custom_log.NewCustomLog("assign_menus_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("assign_menus_failed", fmt.Errorf("error_database"))
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err := syncUserMenusTx(tx, userID, req.MenuIDs, r.UserContext.Id, false); err != nil {
		custom_log.NewCustomLog("assign_menus_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("assign_menus_failed", fmt.Errorf("error_database"))
	}

	if err := tx.Commit(); err != nil {
		custom_log.NewCustomLog("assign_menus_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("assign_menus_failed", fmt.Errorf("error_database"))
	}

	return r.GetUserMenus(userID)
}

func (r *AdminRepoImpl) getUserByID(userID int) (*User, *responses.ErrorResponse) {
	user := User{}
	query := `
		SELECT
			id,
			user_uuid,
			user_name,
			login_id,
			COALESCE(email, '') AS email,
			COALESCE(password, '') AS password,
			COALESCE(nickname, '') AS nickname,
			COALESCE(profile, '') AS profile,
			COALESCE(role_id, 0) AS role_id,
			COALESCE(role_name, '') AS role_name,
			COALESCE(is_admin, FALSE) AS is_admin,
			COALESCE(login_session, '') AS login_session,
			COALESCE(status_id, 0) AS status_id,
			COALESCE("order", 0) AS "order",
			COALESCE(created_by, 0) AS created_by,
			COALESCE(created_at, NOW()) AS created_at,
			COALESCE(updated_by, created_by, 0) AS updated_by,
			COALESCE(updated_at, created_at, NOW()) AS updated_at
		FROM tbl_users
		WHERE deleted_at IS NULL
			AND id = $1
		LIMIT 1
	`

	if err := r.DBPool.Get(&user, query, userID); err != nil {
		custom_log.NewCustomLog("get_admin_user_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("get_admin_user_failed", fmt.Errorf("error_database"))
	}

	return &user, nil
}

func syncUserMenusTx(tx *sqlx.Tx, userID int, menuIDsInput []int, createdBy int, allowAll bool) error {
	now, err := utils.GetCurrentAppTime()
	if err != nil {
		return err
	}

	menuIDs := make([]int, 0)
	if err := tx.Select(&menuIDs, `
		SELECT id
		FROM tbl_menus
		WHERE deleted_at IS NULL
		ORDER BY parent_id ASC, "order" ASC, id ASC
	`); err != nil {
		return err
	}

	allowedSet := make(map[int]struct{}, len(menuIDsInput))
	if allowAll {
		for _, menuID := range menuIDs {
			allowedSet[menuID] = struct{}{}
		}
	} else {
		deniedSet := make(map[int]struct{}, len(menuIDsInput))
		for _, menuID := range menuIDsInput {
			deniedSet[menuID] = struct{}{}
		}

		for _, menuID := range menuIDs {
			if _, denied := deniedSet[menuID]; !denied {
				allowedSet[menuID] = struct{}{}
			}
		}
	}

	for _, menuID := range menuIDs {
		isAllow := false
		if _, ok := allowedSet[menuID]; ok {
			isAllow = true
		}

		if _, err := tx.Exec(
			`
			INSERT INTO tbl_users_menus (
				user_id,
				menu_id,
				is_allow,
				status_id,
				created_by,
				created_at
			) VALUES ($1, $2, $3, 1, $4, $5)
			ON CONFLICT (user_id, menu_id) DO UPDATE
				SET is_allow = EXCLUDED.is_allow,
					status_id = EXCLUDED.status_id,
					updated_by = EXCLUDED.created_by,
					updated_at = EXCLUDED.created_at,
					deleted_by = NULL,
					deleted_at = NULL
			`,
			userID,
			menuID,
			isAllow,
			createdBy,
			now,
		); err != nil {
			return err
		}
	}

	return nil
}
