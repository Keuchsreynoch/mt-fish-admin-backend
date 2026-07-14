package admin

import (
	"errors"
	custom_log "fish_shooting_admin_backend/pkg/logs"
	"fish_shooting_admin_backend/pkg/postgres"
	"fmt"

	types "fish_shooting_admin_backend/pkg/model"
	"fish_shooting_admin_backend/pkg/utils"
	"strings"
	"time"

	"fish_shooting_admin_backend/pkg/responses"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type UserMenu struct {
	ID        int        `db:"id" json:"id"`
	MenuUUID  uuid.UUID  `db:"menu_uuid" json:"menu_uuid"`
	Name      string     `db:"name" json:"name"`
	Icon      string     `db:"icon" json:"icon"`
	Path      string     `db:"path" json:"path"`
	ParentID  int        `db:"parent_id" json:"parent_id"`
	StatusID  int        `db:"status_id" json:"-"`
	Order     int        `db:"order" json:"-"`
	CreatedBy int        `db:"created_by" json:"-"`
	CreatedAt time.Time  `db:"created_at" json:"-"`
	UpdatedBy int        `db:"updated_by" json:"-"`
	UpdatedAt time.Time  `db:"updated_at" json:"-"`
	DeletedBy *int       `db:"deleted_by" json:"-"`
	DeletedAt *time.Time `db:"deleted_at" json:"-"`
}

type User struct {
	ID           int       `db:"id" json:"id"`
	UserUUID     uuid.UUID `db:"user_uuid" json:"user_uuid"`
	UserName     string    `db:"user_name" json:"user_name"`
	LoginID      string    `db:"login_id" json:"login_id"`
	Email        string    `db:"email" json:"email"`
	Password     string    `db:"password" json:"-"`
	Nickname     string    `db:"nickname" json:"nickname"`
	Profile      string    `db:"profile" json:"profile"`
	RoleID       int       `db:"role_id" json:"role_id"`
	RoleName     string    `db:"role_name" json:"role_name"`
	IsAdmin      bool      `db:"is_admin" json:"is_admin"`
	LoginSession string    `db:"login_session" json:"login_session"`
	StatusID     int       `db:"status_id" json:"status_id"`
	Order        int       `db:"order" json:"order"`
	CreatedBy    int       `db:"created_by" json:"created_by"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedBy    int       `db:"updated_by" json:"updated_by"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type UserListResponse struct {
	Users []User `json:"users"`
}

type MenuListResponse struct {
	Menus []UserMenu `json:"menus"`
}

type UserInfoResponse struct {
	User  User       `json:"user"`
	Menus []UserMenu `json:"menus"`
}

type UpdateAdminUserResponse struct {
	User  User       `json:"user"`
	Menus []UserMenu `json:"menus"`
}

type UserMenusResponse struct {
	Menus []UserMenu `json:"menus"`
}

type NewUser struct {
	User  User       `json:"user"`
	Menus []UserMenu `json:"menus"`
}

type NewUserRequest struct {
	UserName string `json:"user_name" validate:"required,min=2,max=100"`
	LoginID  string `json:"login_id" validate:"required,min=3,max=100"`
	Password string `json:"password" validate:"required,min=6,max=255"`
	Email    string `json:"email" validate:"omitempty,email,max=50"`
	Nickname string `json:"nickname" validate:"omitempty,max=100"`
	Profile  string `json:"profile" validate:"omitempty,max=100"`
}

func (r *NewUserRequest) Bind(c *fiber.Ctx, v *utils.Validator) error {
	if err := c.BodyParser(r); err != nil {
		custom_log.NewCustomLog("create_admin_user_failed", err.Error(), "error")
		return errors.New(utils.Translate("invalid_body", nil, c))
	}

	if err := v.Validate(r, c); err != nil {
		custom_log.NewCustomLog("create_admin_user_failed", err.Error(), "error")
		return err
	}

	r.UserName = strings.TrimSpace(strings.ToUpper(r.UserName))
	r.LoginID = r.UserName
	r.Password = strings.TrimSpace(r.Password)
	r.Email = strings.TrimSpace(r.Email)
	r.Nickname = strings.TrimSpace(r.Nickname)
	r.Profile = strings.TrimSpace(r.Profile)

	return nil
}

func (r *NewUser) New(form_data NewUserRequest, user_context *types.UserContext, tx *sqlx.Tx) *responses.ErrorResponse {
	if form_data.LoginID == "" || form_data.UserName == "" || form_data.Password == "" {
		err_msg := &responses.ErrorResponse{}
		return err_msg.NewErrorResponse("create_admin_user_failed", fmt.Errorf("missing_required_fields"))
	}

	exists, err := postgres.IsExists("tbl_users", "user_name", form_data.UserName, tx)
	if err != nil {
		custom_log.NewCustomLog("create_admin_user_failed", err.Error(), "error")
		err_msg := &responses.ErrorResponse{}
		return err_msg.NewErrorResponse("create_admin_user_failed", fmt.Errorf("error_database"))
	}
	if exists {
		err_msg := &responses.ErrorResponse{}
		return err_msg.NewErrorResponse("create_admin_user_failed", fmt.Errorf("login_id_already_exists"))
	}

	// Get next member ID
	id, err := postgres.GetSeqNextVal("tbl_members_id_seq", tx)
	if err != nil {
		e := &responses.ErrorResponse{}
		return e.NewErrorResponse("cannot_get_next_sequence_value", fmt.Errorf("login_id_already_exists"))
	}

	now, err := utils.GetCurrentAppTime()
	if err != nil {
		err_msg := &responses.ErrorResponse{}
		return err_msg.NewErrorResponse("register_failed", fmt.Errorf("technical_error"))
	}
	userUUID := uuid.New()

	r.User.ID = *id
	r.User.UserUUID = userUUID
	r.User.UserName = form_data.UserName
	r.User.LoginID = form_data.LoginID
	r.User.Email = form_data.Email
	r.User.Nickname = form_data.Nickname
	r.User.Profile = form_data.Profile
	r.User.RoleID = 1
	r.User.IsAdmin = true
	r.User.StatusID = 1
	r.User.Order = 1
	r.User.CreatedBy = user_context.Id
	r.User.CreatedAt = now

	return nil
}

type UpdateAdminUserRequest struct {
	Password string `json:"password" validate:"required,min=6,max=255"`
}

func (r *UpdateAdminUserRequest) Bind(c *fiber.Ctx, v *utils.Validator) error {
	if err := c.BodyParser(r); err != nil {
		custom_log.NewCustomLog("update_admin_user_failed", err.Error(), "error")
		return errors.New(utils.Translate("invalid_body", nil, c))
	}

	if err := v.Validate(r, c); err != nil {
		custom_log.NewCustomLog("update_admin_user_failed", err.Error(), "error")
		return err
	}

	r.Password = strings.TrimSpace(r.Password)

	return nil
}

type AssignMenusRequest struct {
	MenuIDs []int `json:"menu_ids"`
}

func (r *AssignMenusRequest) Bind(c *fiber.Ctx, _ *utils.Validator) error {
	if err := c.BodyParser(r); err != nil {
		custom_log.NewCustomLog("assign_menus_failed", err.Error(), "error")
		return errors.New(utils.Translate("invalid_body", nil, c))
	}

	return nil
}

type DashboardSummaryResponse struct {
	TotalTurnover        string                          `json:"total_turnover"`
	TotalPayout          string                          `json:"total_payout"`
	TotalCompanyProfit   string                          `json:"total_company_profit"`
	RewardPool           string                          `json:"reward_pool"`
	CurrentCompanyProfit string                          `json:"current_company_profit"`
	CurrentPoolJackpot   string                          `json:"current_pool_jackpot"`
	ThresholdAmount      string                          `json:"threshold_amount"`
	TopWinMembers        []DashboardTopWinMemberResponse `json:"top_win_members"`
}

type DashboardTopWinMemberResponse struct {
	MemberID          int64  `json:"member_id"`
	MemberUUID        string `json:"member_uuid"`
	Username          string `json:"username"`
	TotalBetAmount    string `json:"total_bet_amount"`
	TotalPayoutAmount string `json:"total_payout_amount"`
	TotalWinAmount    string `json:"total_win_amount"`
	WinCount          int64  `json:"win_count"`
}

type dashboardTopWinMemberRow struct {
	MemberID          int64           `db:"member_id"`
	MemberUUID        string          `db:"member_uuid"`
	Username          string          `db:"username"`
	TotalBetAmount    decimal.Decimal `db:"total_bet_amount"`
	TotalPayoutAmount decimal.Decimal `db:"total_payout_amount"`
	TotalWinAmount    decimal.Decimal `db:"total_win_amount"`
	WinCount          int64           `db:"win_count"`
}

type DashboardTotals struct {
	TotalTurnover        decimal.Decimal `db:"total_turnover"`
	TotalPayout          decimal.Decimal `db:"total_payout"`
	TotalCompanyProfit   decimal.Decimal `db:"total_company_profit"`
	RewardPool           decimal.Decimal `db:"reward_pool"`
	CurrentCompanyProfit decimal.Decimal `db:"current_company_profit"`
}

type CurrentJackpot struct {
	CurrentPoolJackpot decimal.Decimal `db:"current_pool_jackpot"`
	ThresholdAmount    decimal.Decimal `db:"threshold_amount"`
}
