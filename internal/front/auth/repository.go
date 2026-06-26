package auth

import (
	custom_log "fish_shooting_admin_backend/pkg/logs"
	"fish_shooting_admin_backend/pkg/responses"
	"fish_shooting_admin_backend/pkg/utils"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

type AuthRepo interface {
	Login(login_id string, password string) (*LoginReponse, *responses.ErrorResponse)
	GetUserByUUID(user_uuid string) (*UserInfo, error)
	GetUserInfoWithMenus(user_uuid string) (*UserInfoWithMenusResponse, *responses.ErrorResponse)
	Register(register_req RegisterRequest) (*RegisterResponse, *responses.ErrorResponse)
}

type AuthRepoImpl struct {
	DBPool *sqlx.DB
}

func NewAuthRepoImpl(db_pool *sqlx.DB) *AuthRepoImpl {
	return &AuthRepoImpl{
		DBPool: db_pool,
	}
}

func (au *AuthRepoImpl) Login(login_id string, password string) (*LoginReponse, *responses.ErrorResponse) {
	var users []User

	// prepare sql
	sql := `
		SELECT
			user_uuid
		FROM tbl_users
		WHERE deleted_at IS NULL 
		AND (LOWER(login_id) = LOWER($1) OR LOWER(user_name) = LOWER($1))
		AND password = $2
		LIMIT 1
	`

	// execute request
	if err := au.DBPool.Select(&users, sql, login_id, password); err != nil {
		custom_log.NewCustomLog("login_failed", err.Error(), "error")
		err_msg := &responses.ErrorResponse{}
		return nil, err_msg.NewErrorResponse("login_failed", fmt.Errorf("login_id_or_password_invalid"))
	}

	if len(users) == 0 {
		custom_log.NewCustomLog("login_failed", "no_user_found", "error")
		err_msg := &responses.ErrorResponse{}
		return nil, err_msg.NewErrorResponse("login_failed", fmt.Errorf("login_id_or_password_invalid"))
	}

	user := users[0]

	hours := utils.GetenvInt("JWT_EXP_HOUR", 7)
	expirationTime := time.Now().Add(time.Duration(hours) * time.Hour)

	// create the JWT claims
	login_session, _ := uuid.NewV7()
	claims := jwt.MapClaims{
		"user_uuid":     user.UserUUID,
		"login_session": login_session.String(),
		"exp":           expirationTime.Unix(),
	}

	errs := godotenv.Load()
	if errs != nil {
		log.Fatalf("Error loading .env file")
	}
	secret_key := os.Getenv("JWT_SECRET_KEY")

	// update login session & last_login_at
	update_sql := `
		UPDATE tbl_users SET
			login_session = $1,
			last_login_at = $2
			WHERE deleted_at IS NULL
		AND user_uuid = $3
	`
	now, err := utils.GetCurrentAppTime()
	if err != nil {
		custom_log.NewCustomLog("login_failed", err.Error(), "error")
		err_msg := &responses.ErrorResponse{}
		return nil, err_msg.NewErrorResponse("login_failed", fmt.Errorf("error_database"))
	}

	// execute request
	_, err = au.DBPool.Exec(update_sql, login_session, now, user.UserUUID)
	if err != nil {
		custom_log.NewCustomLog("login_failed", err.Error(), "error")
		err_msg := &responses.ErrorResponse{}
		return nil, err_msg.NewErrorResponse("login_failed", fmt.Errorf("error_database"))
	}

	// create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret_key))
	if err != nil {
		custom_log.NewCustomLog("login_failed", err.Error(), "error")
		err_msg := &responses.ErrorResponse{}
		return nil, err_msg.NewErrorResponse("login_failed", fmt.Errorf("error_create_token"))
	}

	return &LoginReponse{
		Auth: Auth{
			Token:     tokenString,
			TokenType: "JWT",
		},
	}, nil
}

func (au *AuthRepoImpl) GetUserByUUID(user_uuid string) (*UserInfo, error) {
	var user_info UserInfo

	// prepare sql
	sql := `
		SELECT
			id, user_uuid, user_name,
			login_id, login_session, status_id
		FROM tbl_users
		WHERE deleted_at IS NULL
		AND user_uuid = $1
	`

	// execute request
	if err := au.DBPool.Get(&user_info, sql, user_uuid); err != nil {
		custom_log.NewCustomLog("get_user_info_failed", err.Error(), "error")
		return nil, err
	}

	return &user_info, nil
}

func (au *AuthRepoImpl) GetUserInfoWithMenus(user_uuid string) (*UserInfoWithMenusResponse, *responses.ErrorResponse) {
	userInfo, err := au.GetUserByUUID(user_uuid)
	if err != nil {
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("get_userinfo_failed", fmt.Errorf("error_database"))
	}

	menus := make([]UserMenu, 0)
	query := `
		SELECT
			m.id,
			m.menu_uuid,
			m.name,
			COALESCE(m.icon, '') AS icon,
			COALESCE(m.path, '') AS path,
			COALESCE(m.parent_id, 0) AS parent_id,
			COALESCE(m.status_id, 0) AS status_id,
			COALESCE(m."order", 0) AS "order",
			COALESCE(m.created_at, NOW()) AS created_at,
			COALESCE(m.updated_at, m.created_at, NOW()) AS updated_at
		FROM tbl_users_menus um
		JOIN tbl_menus m
			ON m.id = um.menu_id
		WHERE um.deleted_at IS NULL
			AND um.is_allow = TRUE
			AND m.deleted_at IS NULL
			AND um.user_id = $1
		ORDER BY m.parent_id ASC, m."order" ASC, m.id ASC
	`

	if err := au.DBPool.Select(&menus, query, userInfo.ID); err != nil {
		custom_log.NewCustomLog("get_user_info_failed", err.Error(), "error")
		e := &responses.ErrorResponse{}
		return nil, e.NewErrorResponse("get_userinfo_failed", fmt.Errorf("error_database"))
	}

	return &UserInfoWithMenusResponse{
		UserInfo: *userInfo,
		Menus:    menus,
	}, nil
}

func (au *AuthRepoImpl) Register(register_req RegisterRequest) (*RegisterResponse, *responses.ErrorResponse) {
	var register_model RegisterModel

	// create register model
	if err := register_model.New(register_req, au.DBPool); err != nil {
		custom_log.NewCustomLog(err.MessageID, err.Detail.Error(), "error")
		err_msg := &responses.ErrorResponse{}
		return nil, err_msg.NewErrorResponse(err.MessageID, err.Err)
	}

	// prepare query
	query := `
		INSERT INTO tbl_users (
			id, user_uuid, user_name, login_id,
			password, phone_number, profile_photo,
			status_id, "order", created_by, created_at
		) VALUES (
			:id, :user_uuid, :user_name, :login_id,
			:password, :phone_number, :profile_photo,
			:status_id, :order, :created_by, :created_at
		)
	`

	// execute request
	_, err := au.DBPool.NamedExec(query, register_model)
	if err != nil {
		custom_log.NewCustomLog("register_failed", err.Error(), "error")
		err_msg := &responses.ErrorResponse{}
		return nil, err_msg.NewErrorResponse("register_failed", fmt.Errorf("cannot_insert_db_error"))
	}

	return &RegisterResponse{
		UserInfo: register_model,
	}, nil
}
