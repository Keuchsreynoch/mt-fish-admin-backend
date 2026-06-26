package auth

import (
	"errors"
	custom_log "fish_shooting_admin_backend/pkg/logs"
	"fish_shooting_admin_backend/pkg/postgres"
	"fish_shooting_admin_backend/pkg/responses"
	"fish_shooting_admin_backend/pkg/utils"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type LoginRequest struct {
	UserName string `json:"user_name"`
	LoginID  string `json:"login_id"`
	Password string `json:"password" validate:"required"`
}

func (au *LoginRequest) Bind(c *fiber.Ctx, v *utils.Validator) error {
	if err := c.BodyParser(au); err != nil {
		custom_log.NewCustomLog("login_failed", err.Error(), "error")
		return errors.New(utils.Translate("invalid_body", nil, c))
	}

	if err := v.Validate(au, c); err != nil {
		custom_log.NewCustomLog("login_failed", err.Error(), "error")
		return err
	}

	if strings.TrimSpace(au.UserName) == "" && strings.TrimSpace(au.LoginID) == "" {
		custom_log.NewCustomLog("login_failed", "missing login identifier", "error")
		return errors.New(utils.Translate("login_id_or_user_name_is_required", nil, c))
	}

	return nil
}

type LoginReponse struct {
	Auth Auth `json:"auth"`
}

type Auth struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
}

type User struct {
	UserUUID uuid.UUID `json:"user_uuid" db:"user_uuid"`
}

type UserInfo struct {
	ID           int    `json:"id" db:"id"`
	UserUUID     string `json:"user_uuid" db:"user_uuid"`
	UserName     string `json:"user_name" db:"user_name"`
	LoginID      string `json:"login_id" db:"login_id"`
	LoginSession string `json:"login_session" db:"login_session"`
	StatusID     int    `json:"status_id" db:"status_id"`
}

type UserMenu struct {
	ID        int       `json:"id" db:"id"`
	MenuUUID  uuid.UUID `json:"menu_uuid" db:"menu_uuid"`
	Name      string    `json:"name" db:"name"`
	Icon      string    `json:"icon" db:"icon"`
	Path      string    `json:"path" db:"path"`
	ParentID  int       `json:"parent_id" db:"parent_id"`
	StatusID  int       `json:"status_id" db:"status_id"`
	Order     int       `json:"order" db:"order"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UserInfoWithMenusResponse struct {
	UserInfo UserInfo   `json:"user_info"`
	Menus    []UserMenu `json:"menus"`
}

type RegisterRequest struct {
	UserName        string `json:"user_name" validate:"required,min=2,max=100"`
	LoginID         string `json:"login_id" validate:"required,min=3,max=100"`
	Password        string `json:"password" validate:"required,min=6,max=255"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=6,max=255"`
	PhoneNumber     string `json:"phone_number" validate:"required,max=100"`
	ProfilePhoto    string `json:"profile_photo" validate:"omitempty"`
}

func (au *RegisterRequest) Bind(c *fiber.Ctx, v *utils.Validator) error {
	if err := c.BodyParser(au); err != nil {
		custom_log.NewCustomLog("register_failed", err.Error(), "error")
		return errors.New(utils.Translate("invalid_body", nil, c))
	}

	if err := v.Validate(au, c); err != nil {
		custom_log.NewCustomLog("register_failed", err.Error(), "error")
		return err
	}

	if au.Password != au.ConfirmPassword {
		custom_log.NewCustomLog("register_failed", "passwords do not match", "error")
		return errors.New(utils.Translate("passwords_do_not_match", nil, c))
	}

	return nil
}

type RegisterModel struct {
	ID           uint64    `db:"id" json:"-"`
	UserUUID     string    `db:"user_uuid" json:"user_uuid"`
	UserName     string    `db:"user_name" json:"user_name"`
	LoginID      string    `db:"login_id" json:"login_id"`
	Password     string    `db:"password" json:"password"`
	PhoneNumber  string    `db:"phone_number" json:"phone_number"`
	ProfilePhoto string    `db:"profile_photo" json:"profile_photo"`
	StatusID     uint64    `db:"status_id" json:"status_id"`
	Order        uint64    `db:"order" json:"order"`
	CreatedBy    uint64    `db:"created_by" json:"created_by"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

func (au *RegisterModel) New(register_req RegisterRequest, conn *sqlx.DB) *responses.ErrorWithDetailResponse {
	// check login_id exists
	loginID := strings.TrimSpace(strings.ToLower(register_req.LoginID))
	is_login_id, err := postgres.IsExists("tbl_users", "login_id", loginID, conn)
	if err != nil {
		err_msg := &responses.ErrorWithDetailResponse{}
		return err_msg.NewErrorResponse("register_failed", fmt.Errorf("cannot_check_login_id"), err)
	} else {
		if is_login_id {
			err_msg := &responses.ErrorWithDetailResponse{}
			return err_msg.NewErrorResponse("register_failed", fmt.Errorf("login_id_already_exists"), fmt.Errorf("register login_id is already exists"))
		}
	}

	if register_req.ProfilePhoto == "" {
		register_req.ProfilePhoto = "user1.png"
	}

	// generate new UUID
	uuid, _ := uuid.NewV7()

	now, err := utils.GetCurrentAppTime()
	if err != nil {
		err_msg := &responses.ErrorWithDetailResponse{}
		return err_msg.NewErrorResponse("register_failed", fmt.Errorf("technical_error"), err)
	}

	// get next sequence ID
	id, err := postgres.GetSeqNextVal("tbl_users_id_seq", conn)
	if err != nil {
		err_msg := &responses.ErrorWithDetailResponse{}
		return err_msg.NewErrorResponse("register_failed", fmt.Errorf("technical_error"), err)
	}

	// assign values to the RegisterModel
	au.ID = uint64(*id)
	au.UserUUID = uuid.String()
	au.UserName = strings.TrimSpace(register_req.UserName)
	au.LoginID = loginID
	au.Password = register_req.Password
	au.PhoneNumber = strings.TrimSpace(register_req.PhoneNumber)
	au.ProfilePhoto = register_req.ProfilePhoto
	au.StatusID = 1
	au.Order = uint64(*id)
	au.CreatedBy = uint64(*id)
	au.CreatedAt = now

	return nil
}

type RegisterResponse struct {
	UserInfo RegisterModel `json:"user_info"`
}
